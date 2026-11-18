package teacherhandler

import (
	"net/http"

	"github.com/geniussociety/backend/internal/handler"
	"github.com/geniussociety/backend/internal/middleware"
	teachersvc "github.com/geniussociety/backend/internal/service/teacher"
	"github.com/gin-gonic/gin"
)

// WalasHandler handles HTTP requests for homeroom teacher views.
type WalasHandler struct {
	svc teachersvc.WalasService
}

// NewWalasHandler constructs a new WalasHandler.
func NewWalasHandler(svc teachersvc.WalasService) *WalasHandler {
	return &WalasHandler{svc: svc}
}

// GetWalasStudents handles GET /api/teacher/walas/students
func (h *WalasHandler) GetWalasStudents(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)
	data, err := h.svc.GetWalasStudents(c.Request.Context(), scope)
	if err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}

// GetWalasGrades handles GET /api/teacher/walas/grades
func (h *WalasHandler) GetWalasGrades(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)
	data, err := h.svc.GetWalasGrades(c.Request.Context(), scope)
	if err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}
