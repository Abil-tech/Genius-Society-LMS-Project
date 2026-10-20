package storage

import (
	"context"
	"fmt"
	"io"
	"mime"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/geniussociety/backend/internal/config"
	"github.com/geniussociety/backend/internal/model"
)


const (
	MaxFileSize = 50 * 1024 * 1024 // 50 MB

	// Allowed MIME types
	allowedMimes = "application/pdf,image/jpeg,image/png,image/gif,image/webp," +
		"video/mp4,video/mpeg,audio/mpeg,audio/mp3," +
		"application/msword,application/vnd.openxmlformats-officedocument.wordprocessingml.document," +
		"application/vnd.ms-excel,application/vnd.openxmlformats-officedocument.spreadsheetml.sheet," +
		"application/vnd.ms-powerpoint,application/vnd.openxmlformats-officedocument.presentationml.presentation," +
		"text/plain,application/zip,application/x-zip-compressed"
)

type B2Storage struct {
	client     *s3.Client
	bucketName string
	publicURL  string
}

func NewB2Storage(cfg *config.Config) (*B2Storage, error) {
	if cfg.B2ApplicationKeyID == "" || cfg.B2ApplicationKey == "" || cfg.B2BucketName == "" || cfg.B2S3Endpoint == "" {
		return nil, fmt.Errorf("B2 storage config incomplete")
	}

	// Normalize endpoint: ensure https:// prefix
	endpoint := cfg.B2S3Endpoint
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		endpoint = "https://" + endpoint
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.B2ApplicationKeyID,
			cfg.B2ApplicationKey,
			"",
		)),
		awsconfig.WithRegion("auto"),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.BaseEndpoint = &endpoint
	})

	return &B2Storage{
		client:     client,
		bucketName: cfg.B2BucketName,
		publicURL:  cfg.B2PublicURL,
	}, nil
}

// UploadFile uploads a file to B2 and returns a FileRef.
// kind = materials | assignments | projects | submissions
func (s *B2Storage) UploadFile(ctx context.Context, kind, originalFilename, mimeType string, size int64, body io.Reader) (*model.FileRef, error) {
	if size > MaxFileSize {
		return nil, fmt.Errorf("file size exceeds limit of %d MB", MaxFileSize/1024/1024)
	}

	if !isAllowedMIME(mimeType) {
		return nil, fmt.Errorf("mime type %s is not allowed", mimeType)
	}

	ext := filepath.Ext(originalFilename)
	if ext == "" {
		exts, _ := mime.ExtensionsByType(mimeType)
		if len(exts) > 0 {
			ext = exts[0]
		}
	}

	key := fmt.Sprintf("%s/%s/%s%s", kind, time.Now().Format("2006/01"), uuid.NewString(), ext)

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(s.bucketName),
		Key:           aws.String(key),
		Body:          body,
		ContentType:   aws.String(mimeType),
		ContentLength: aws.Int64(size),
	})
	if err != nil {
		return nil, fmt.Errorf("upload failed: %w", err)
	}

	return &model.FileRef{
		Key:  key,
		Name: originalFilename,
		Size: size,
		Mime: mimeType,
	}, nil
}

// GetPresignedURL generates a short-lived presigned download URL.
func (s *B2Storage) GetPresignedURL(ctx context.Context, key string, expirySeconds int64) (string, error) {
	presigner := s3.NewPresignClient(s.client)
	req, err := presigner.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(time.Duration(expirySeconds)*time.Second))
	if err != nil {
		return "", err
	}
	return req.URL, nil
}

// DeleteFile removes a file from B2.
func (s *B2Storage) DeleteFile(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	return err
}

func isAllowedMIME(mimeType string) bool {
	allowed := strings.Split(allowedMimes, ",")
	for _, a := range allowed {
		if strings.EqualFold(strings.TrimSpace(a), strings.TrimSpace(mimeType)) {
			return true
		}
	}
	return false
}
