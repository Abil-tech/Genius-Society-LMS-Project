package teacherhandler

import (
	"net/http"

	"github.com/geniussociety/backend/internal/handler"
	"github.com/geniussociety/backend/internal/middleware"
	teachersvc "github.com/geniussociety/backend/internal/service/teacher"
	"github.com/gin-gonic/gin"
)

// DashboardHandler handles HTTP requests for teacher dashboard.
type DashboardHandler struct {
	svc teachersvc.DashboardService
}

// NewDashboardHandler constructs a new DashboardHandler instance.
func NewDashboardHandler(svc teachersvc.DashboardService) *DashboardHandler {
	return &DashboardHandler{svc: svc}
}

// GetDashboard handles GET /api/teacher/dashboard.
func (h *DashboardHandler) GetDashboard(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)
	data, err := h.svc.GetDashboardData(c.Request.Context(), scope)
	if err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}
