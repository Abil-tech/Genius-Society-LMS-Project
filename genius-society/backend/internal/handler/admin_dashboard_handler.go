package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/geniussociety/backend/internal/dto"
	"github.com/geniussociety/backend/internal/service"
)

type AdminDashboardHandler struct {
	svc *service.AdminDashboardService
}

func NewAdminDashboardHandler(svc *service.AdminDashboardService) *AdminDashboardHandler {
	return &AdminDashboardHandler{svc: svc}
}

// GET /api/admin/dashboard
func (h *AdminDashboardHandler) GetMetrics(c *gin.Context) {
	metrics, err := h.svc.GetMetrics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "INTERNAL_ERROR", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Data:    metrics,
	})
}
