package teacherhandler

import (
	"net/http"
	"strings"
	"time"

	"github.com/geniussociety/backend/internal/handler"
	"github.com/geniussociety/backend/internal/storage"
	"github.com/gin-gonic/gin"
)

// FileHandler handles generic file operations for teacher routes.
type FileHandler struct {
	storage storage.FileStorage
}

// NewFileHandler constructs a new FileHandler.
func NewFileHandler(st storage.FileStorage) *FileHandler {
	return &FileHandler{storage: st}
}

// GetPresignedDownloadURL handles GET /api/teacher/files/download?key=...
// Returns a short-lived presigned URL for the given R2 object key.
// The teacher must supply a non-empty key; access control is enforced by R2 scope.
func (h *FileHandler) GetPresignedDownloadURL(c *gin.Context) {
	key := strings.TrimSpace(c.Query("key"))
	if key == "" {
		handler.RespondError(c, handler.NewValidationError("Parameter key tidak boleh kosong"))
		return
	}

	// Basic sanity check: reject path traversal attempts.
	if strings.Contains(key, "..") {
		handler.RespondError(c, handler.NewValidationError("Parameter key tidak valid"))
		return
	}

	url, err := h.storage.PresignGet(c.Request.Context(), key, 5*time.Minute)
	if err != nil {
		if strings.Contains(err.Error(), "not configured") {
			handler.RespondError(c, handler.ErrStorageNotConfigured)
			return
		}
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"url":       url,
			"expiresIn": 300,
		},
	})
}
