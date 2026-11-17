package storage

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"testing"
	"time"
)

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"normal.pdf", "normal.pdf"},
		{"../../../etc/passwd", "passwd"},
		{"nested/path/to/file.docx", "file.docx"},
		{"nested\\windows\\path.xlsx", "path.xlsx"},
		{"file with spaces.png", "file with spaces.png"},
		{"file$with@bad!chars.jpg", "file_with_bad_chars.jpg"},
		{"", "file"},
		{".", "file"},
	}

	for _, tt := range tests {
		got := SanitizeFilename(tt.input)
		if got != tt.expected {
			t.Errorf("SanitizeFilename(%q) = %q; want %q", tt.input, got, tt.expected)
		}
	}
}

func TestValidateFileCount(t *testing.T) {
	if err := ValidateFileCount(5); err != nil {
		t.Errorf("ValidateFileCount(5) returned error: %v", err)
	}
	if err := ValidateFileCount(6); !errors.Is(err, ErrFileCountExceeded) {
		t.Errorf("ValidateFileCount(6) expected ErrFileCountExceeded, got: %v", err)
	}
}

func TestB2StorageAdapter_Unconfigured(t *testing.T) {
	adapter := NewFileStorage(nil)
	ctx := context.Background()

	_, err := adapter.Upload(ctx, "materials", "test.pdf", "application/pdf", 100, bytes.NewReader([]byte("test")))
	if !errors.Is(err, ErrStorageNotConfigured) {
		t.Errorf("expected ErrStorageNotConfigured, got %v", err)
	}

	_, err = adapter.PresignGet(ctx, "key", 1*time.Minute)
	if !errors.Is(err, ErrStorageNotConfigured) {
		t.Errorf("expected ErrStorageNotConfigured, got %v", err)
	}

	err = adapter.Delete(ctx, "key")
	if !errors.Is(err, ErrStorageNotConfigured) {
		t.Errorf("expected ErrStorageNotConfigured, got %v", err)
	}
}

func createZipBuffer(files map[string]string) ([]byte, error) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	for name, content := range files {
		w, err := zw.Create(name)
		if err != nil {
			return nil, err
		}
		if _, err := w.Write([]byte(content)); err != nil {
			return nil, err
		}
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func TestResolveAndValidateMime(t *testing.T) {
	// PDF test
	pdfBytes := []byte("%PDF-1.4 test")
	mime, err := resolveAndValidateMime(bytes.NewReader(pdfBytes), int64(len(pdfBytes)), ".pdf", "application/pdf")
	if err != nil || mime != "application/pdf" {
		t.Fatalf("expected application/pdf, got %s, err: %v", mime, err)
	}

	// Invalid PDF MIME
	_, err = resolveAndValidateMime(bytes.NewReader(pdfBytes), int64(len(pdfBytes)), ".pdf", "text/plain")
	if !errors.Is(err, ErrMimeNotAllowed) {
		t.Fatalf("expected ErrMimeNotAllowed, got %v", err)
	}

	// Valid DOCX zip archive
	docxBytes, err := createZipBuffer(map[string]string{
		"[Content_Types].xml": "<Types></Types>",
		"word/document.xml":   "<document></document>",
	})
	if err != nil {
		t.Fatalf("failed to create zip buffer: %v", err)
	}
	mime, err = resolveAndValidateMime(bytes.NewReader(docxBytes), int64(len(docxBytes)), ".docx", "application/zip")
	if err != nil || mime != "application/vnd.openxmlformats-officedocument.wordprocessingml.document" {
		t.Fatalf("expected docx mime, got %s, err: %v", mime, err)
	}

	// Corrupt / non-office zip as DOCX
	badZipBytes, err := createZipBuffer(map[string]string{
		"unrelated.txt": "hello",
	})
	if err != nil {
		t.Fatalf("failed to create bad zip buffer: %v", err)
	}
	_, err = resolveAndValidateMime(bytes.NewReader(badZipBytes), int64(len(badZipBytes)), ".docx", "application/zip")
	if !errors.Is(err, ErrMimeNotAllowed) {
		t.Fatalf("expected ErrMimeNotAllowed for invalid zip content, got %v", err)
	}
}
