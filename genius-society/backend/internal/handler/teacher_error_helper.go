package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/geniussociety/backend/internal/dto"
	"github.com/geniussociety/backend/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Standard domain error sentinel definitions
var (
	ErrNotFound               = errors.New("data tidak ditemukan")
	ErrForbidden              = errors.New("akses ditolak")
	ErrConflict               = errors.New("terjadi konflik pada data")
	ErrValidation             = errors.New("validasi gagal")
	ErrNotAssigned            = errors.New("guru tidak ditugaskan pada kelas atau mata pelajaran ini")
	ErrAcademicYearNotActive  = errors.New("tahun ajaran tidak aktif")
	ErrStorageNotConfigured   = errors.New("layanan penyimpanan belum terkonfigurasi")
	ErrUnauthorized           = errors.New("autentikasi gagal atau sesi berakhir")
)

// DomainError represents a structured domain error with code, message, and HTTP status.
type DomainError struct {
	Status  int
	Code    string
	Message string
	Err     error
}

func (e *DomainError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Code
}

func (e *DomainError) Unwrap() error {
	return e.Err
}

func NewDomainError(status int, code, message string) *DomainError {
	return &DomainError{
		Status:  status,
		Code:    code,
		Message: message,
	}
}

func NewNotFoundError(message string) *DomainError {
	return &DomainError{
		Status:  http.StatusNotFound,
		Code:    "NOT_FOUND",
		Message: message,
	}
}

func NewForbiddenError(message string) *DomainError {
	return &DomainError{
		Status:  http.StatusForbidden,
		Code:    "FORBIDDEN",
		Message: message,
	}
}

func NewConflictError(message string) *DomainError {
	return &DomainError{
		Status:  http.StatusConflict,
		Code:    "CONFLICT",
		Message: message,
	}
}

func NewValidationError(message string) *DomainError {
	return &DomainError{
		Status:  http.StatusBadRequest,
		Code:    "VALIDATION_ERROR",
		Message: message,
	}
}

// respondError maps domain errors to proper HTTP status and dto.ErrorBody.
// Unknown errors are logged and answered with safe "Terjadi kesalahan internal pada server".
func respondError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	var dErr *DomainError
	if errors.As(err, &dErr) {
		c.JSON(dErr.Status, dto.ErrorResponse{
			Success: false,
			Error: dto.ErrorBody{
				Code:    dErr.Code,
				Message: dErr.Message,
			},
		})
		return
	}

	// Match sentinel errors
	switch {
	case errors.Is(err, ErrNotFound):
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Success: false,
			Error: dto.ErrorBody{
				Code:    "NOT_FOUND",
				Message: "Data tidak ditemukan",
			},
		})
	case errors.Is(err, ErrForbidden):
		c.JSON(http.StatusForbidden, dto.ErrorResponse{
			Success: false,
			Error: dto.ErrorBody{
				Code:    "FORBIDDEN",
				Message: "Akses ditolak",
			},
		})
	case errors.Is(err, ErrNotAssigned):
		c.JSON(http.StatusForbidden, dto.ErrorResponse{
			Success: false,
			Error: dto.ErrorBody{
				Code:    "NOT_ASSIGNED",
				Message: "Guru tidak ditugaskan pada kelas atau mata pelajaran ini",
			},
		})
	case errors.Is(err, ErrConflict):
		c.JSON(http.StatusConflict, dto.ErrorResponse{
			Success: false,
			Error: dto.ErrorBody{
				Code:    "CONFLICT",
				Message: "Terjadi konflik pada data",
			},
		})
	case errors.Is(err, ErrAcademicYearNotActive):
		c.JSON(http.StatusConflict, dto.ErrorResponse{
			Success: false,
			Error: dto.ErrorBody{
				Code:    "ACADEMIC_YEAR_NOT_ACTIVE",
				Message: "Tahun ajaran tidak aktif",
			},
		})
	case errors.Is(err, ErrValidation):
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error: dto.ErrorBody{
				Code:    "VALIDATION_ERROR",
				Message: "Data yang dikirim tidak valid",
			},
		})
	case errors.Is(err, ErrStorageNotConfigured), errors.Is(err, storage.ErrStorageNotConfigured):
		c.JSON(http.StatusServiceUnavailable, dto.ErrorResponse{
			Success: false,
			Error: dto.ErrorBody{
				Code:    "STORAGE_NOT_CONFIGURED",
				Message: "Layanan penyimpanan berkas belum terkonfigurasi",
			},
		})
	case errors.Is(err, ErrUnauthorized):
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Error: dto.ErrorBody{
				Code:    "UNAUTHORIZED",
				Message: "Autentikasi gagal atau sesi berakhir",
			},
		})
	default:
		// Log internal error safely
		log.Printf("[TeacherInternalError] %v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Error: dto.ErrorBody{
				Code:    "INTERNAL_ERROR",
				Message: "Terjadi kesalahan internal pada server",
			},
		})
	}
}

// bindError transforms binding and JSON validation errors into safe client-facing messages.
func bindError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		var fieldNames []string
		for _, fe := range ve {
			fieldNames = append(fieldNames, strings.ToLower(fe.Field()))
		}
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error: dto.ErrorBody{
				Code:    "VALIDATION_ERROR",
				Message: fmt.Sprintf("Data tidak valid: %s", strings.Join(fieldNames, ", ")),
			},
		})
		return
	}

	c.JSON(http.StatusBadRequest, dto.ErrorResponse{
		Success: false,
		Error: dto.ErrorBody{
			Code:    "VALIDATION_ERROR",
			Message: "Data tidak valid: format request salah",
		},
	})
}
