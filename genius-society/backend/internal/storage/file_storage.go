package storage

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/geniussociety/backend/internal/model"
)

// Sentinel errors for storage operations
var (
	ErrStorageNotConfigured = errors.New("layanan penyimpanan belum terkonfigurasi")
	ErrInvalidFile          = errors.New("file tidak valid")
	ErrFileTooLarge         = errors.New("ukuran file melebihi batas maksimal")
	ErrMimeNotAllowed       = errors.New("tipe file tidak diizinkan")
	ErrFileCountExceeded    = errors.New("maksimal 5 file per resource")
)

const (
	MaxFilesPerResource = 5
	MaxPresignDuration  = 5 * time.Minute
)

// FileStorage defines the standard interface for file operations.
type FileStorage interface {
	Upload(ctx context.Context, kind, filename, mimeType string, size int64, body io.Reader) (*model.FileRef, error)
	PresignGet(ctx context.Context, key string, duration time.Duration) (string, error)
	Delete(ctx context.Context, key string) error
}

// B2StorageAdapter adapts B2Storage to FileStorage interface.
type B2StorageAdapter struct {
	b2 *B2Storage
}

// NewFileStorage creates a new FileStorage instance wrapping B2Storage.
func NewFileStorage(b2 *B2Storage) FileStorage {
	return &B2StorageAdapter{b2: b2}
}

// Upload uploads a file to B2 storage.
func (a *B2StorageAdapter) Upload(ctx context.Context, kind, filename, mimeType string, size int64, body io.Reader) (*model.FileRef, error) {
	if a == nil || a.b2 == nil {
		return nil, ErrStorageNotConfigured
	}
	return a.b2.UploadFile(ctx, kind, filename, mimeType, size, body)
}

// PresignGet creates a presigned download URL capped at 5 minutes (300 seconds).
func (a *B2StorageAdapter) PresignGet(ctx context.Context, key string, duration time.Duration) (string, error) {
	if a == nil || a.b2 == nil {
		return "", ErrStorageNotConfigured
	}
	if duration <= 0 || duration > MaxPresignDuration {
		duration = MaxPresignDuration
	}
	return a.b2.GetPresignedURL(ctx, key, int64(duration.Seconds()))
}

// Delete deletes a file from B2 storage.
func (a *B2StorageAdapter) Delete(ctx context.Context, key string) error {
	if a == nil || a.b2 == nil {
		return ErrStorageNotConfigured
	}
	return a.b2.DeleteFile(ctx, key)
}

// ValidateFileCount checks if file count is within allowed limit.
func ValidateFileCount(count int) error {
	if count > MaxFilesPerResource {
		return ErrFileCountExceeded
	}
	return nil
}

var invalidFilenameChars = regexp.MustCompile(`[^a-zA-Z0-9._\-\s]`)

// SanitizeFilename cleans the filename, removes path traversal and invalid characters.
func SanitizeFilename(name string) string {
	base := filepath.Base(name)
	base = strings.ReplaceAll(base, "\\", "")
	base = strings.ReplaceAll(base, "\x00", "")
	base = strings.TrimSpace(base)
	base = invalidFilenameChars.ReplaceAllString(base, "_")
	if len(base) > 200 {
		ext := filepath.Ext(base)
		base = base[:200-len(ext)] + ext
	}
	if base == "" || base == "." {
		base = "file"
	}
	return base
}

// ValidateAndSniffMultipart validates size, sanitizes filename, and sniffs MIME type.
func ValidateAndSniffMultipart(fh *multipart.FileHeader) (multipart.File, string, string, int64, error) {
	if fh == nil {
		return nil, "", "", 0, ErrInvalidFile
	}

	size := fh.Size
	if size <= 0 {
		return nil, "", "", 0, errors.New("file kosong")
	}
	if size > MaxFileSize {
		return nil, "", "", 0, ErrFileTooLarge
	}

	cleanName := SanitizeFilename(fh.Filename)
	ext := strings.ToLower(filepath.Ext(cleanName))

	file, err := fh.Open()
	if err != nil {
		return nil, "", "", 0, fmt.Errorf("gagal membuka file: %w", err)
	}

	// Sniff first 512 bytes
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		file.Close()
		return nil, "", "", 0, fmt.Errorf("gagal membaca file untuk sniffing: %w", err)
	}

	// Reset read pointer
	if seeker, ok := file.(io.Seeker); ok {
		if _, err := seeker.Seek(0, io.SeekStart); err != nil {
			file.Close()
			return nil, "", "", 0, fmt.Errorf("gagal mereset pointer file: %w", err)
		}
	}

	detectedMime := http.DetectContentType(buf[:n])
	validatedMime, err := resolveAndValidateMime(file, size, ext, detectedMime)
	if err != nil {
		file.Close()
		return nil, "", "", 0, err
	}

	return file, cleanName, validatedMime, size, nil
}

// resolveAndValidateMime verifies sniffed MIME against whitelist and handles OOXML (docx, xlsx, pptx).
func resolveAndValidateMime(r io.ReaderAt, size int64, ext, detectedMime string) (string, error) {
	// Base detected MIME normalization
	baseDetected := strings.Split(detectedMime, ";")[0]
	baseDetected = strings.TrimSpace(baseDetected)

	switch ext {
	case ".pdf":
		if baseDetected == "application/pdf" {
			return "application/pdf", nil
		}
	case ".jpg", ".jpeg":
		if baseDetected == "image/jpeg" {
			return "image/jpeg", nil
		}
	case ".png":
		if baseDetected == "image/png" {
			return "image/png", nil
		}
	case ".gif":
		if baseDetected == "image/gif" {
			return "image/gif", nil
		}
	case ".webp":
		if baseDetected == "image/webp" {
			return "image/webp", nil
		}
	case ".mp4":
		if baseDetected == "video/mp4" {
			return "video/mp4", nil
		}
	case ".mp3":
		if baseDetected == "audio/mpeg" || baseDetected == "audio/mp3" {
			return "audio/mpeg", nil
		}
	case ".txt":
		if strings.HasPrefix(baseDetected, "text/plain") {
			return "text/plain", nil
		}
	case ".zip":
		if baseDetected == "application/zip" || baseDetected == "application/x-zip-compressed" {
			return "application/zip", nil
		}
	case ".docx":
		if isZipCompatible(baseDetected) {
			if err := verifyOfficeZip(r, size, "word/"); err == nil {
				return "application/vnd.openxmlformats-officedocument.wordprocessingml.document", nil
			}
		}
	case ".xlsx":
		if isZipCompatible(baseDetected) {
			if err := verifyOfficeZip(r, size, "xl/"); err == nil {
				return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", nil
			}
		}
	case ".pptx":
		if isZipCompatible(baseDetected) {
			if err := verifyOfficeZip(r, size, "ppt/"); err == nil {
				return "application/vnd.openxmlformats-officedocument.presentationml.presentation", nil
			}
		}
	case ".doc":
		if baseDetected == "application/msword" || baseDetected == "application/octet-stream" {
			return "application/msword", nil
		}
	case ".xls":
		if baseDetected == "application/vnd.ms-excel" || baseDetected == "application/octet-stream" {
			return "application/vnd.ms-excel", nil
		}
	case ".ppt":
		if baseDetected == "application/vnd.ms-powerpoint" || baseDetected == "application/octet-stream" {
			return "application/vnd.ms-powerpoint", nil
		}
	}

	return "", fmt.Errorf("%w: ekstensi %s dengan tipe konten terdeteksi %s", ErrMimeNotAllowed, ext, detectedMime)
}

func isZipCompatible(mime string) bool {
	return mime == "application/zip" ||
		mime == "application/x-zip-compressed" ||
		mime == "application/octet-stream"
}

// verifyOfficeZip opens the zip archive and checks for [Content_Types].xml and the required folder.
func verifyOfficeZip(r io.ReaderAt, size int64, requiredFolder string) error {
	zr, err := zip.NewReader(r, size)
	if err != nil {
		return fmt.Errorf("bukan arsip zip yang valid: %w", err)
	}

	hasContentTypes := false
	hasRequiredFolder := false

	for _, f := range zr.File {
		if f.Name == "[Content_Types].xml" {
			hasContentTypes = true
		}
		if strings.HasPrefix(f.Name, requiredFolder) {
			hasRequiredFolder = true
		}
		if hasContentTypes && hasRequiredFolder {
			return nil
		}
	}

	if !hasContentTypes {
		return errors.New("tidak memiliki [Content_Types].xml")
	}
	if !hasRequiredFolder {
		return fmt.Errorf("tidak memiliki direktori %s", requiredFolder)
	}
	return nil
}
