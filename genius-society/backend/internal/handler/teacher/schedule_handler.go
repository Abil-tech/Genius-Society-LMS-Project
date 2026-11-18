package teacherhandler

import (
	"net/http"

	"github.com/geniussociety/backend/internal/handler"
	"github.com/geniussociety/backend/internal/middleware"
	teachersvc "github.com/geniussociety/backend/internal/service/teacher"
	"github.com/gin-gonic/gin"
)

// ScheduleHandler handles requests for teacher schedules.
type ScheduleHandler struct {
	svc teachersvc.ScheduleService
}

// NewScheduleHandler constructs a new ScheduleHandler.
func NewScheduleHandler(svc teachersvc.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{svc: svc}
}

// GetTeacherSchedule handles GET /api/teacher/schedule.
func (h *ScheduleHandler) GetTeacherSchedule(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)
	schedules, err := h.svc.GetTeacherSchedule(c.Request.Context(), scope)
	if err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    schedules,
	})
}
