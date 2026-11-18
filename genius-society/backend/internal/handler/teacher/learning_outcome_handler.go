package teacherhandler

import (
	"net/http"
	"strings"

	teacherdto "github.com/geniussociety/backend/internal/dto/teacher"
	"github.com/geniussociety/backend/internal/handler"
	"github.com/geniussociety/backend/internal/middleware"
	teachersvc "github.com/geniussociety/backend/internal/service/teacher"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// LearningOutcomeHandler handles HTTP requests for Pencapaian Pembelajaran (PP).
type LearningOutcomeHandler struct {
	svc teachersvc.LearningOutcomeService
}

// NewLearningOutcomeHandler constructs a new LearningOutcomeHandler instance.
func NewLearningOutcomeHandler(svc teachersvc.LearningOutcomeService) *LearningOutcomeHandler {
	return &LearningOutcomeHandler{svc: svc}
}

// ListOutcomes handles GET /api/teacher/learning-outcomes.
func (h *LearningOutcomeHandler) ListOutcomes(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)

	subjectID := strings.TrimSpace(c.Query("subjectId"))
	classID := strings.TrimSpace(c.Query("classId"))

	list, err := h.svc.ListOutcomes(c.Request.Context(), scope, subjectID, classID)
	if err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    list,
	})
}

// GetOutcomeDetail handles GET /api/teacher/learning-outcomes/:id.
func (h *LearningOutcomeHandler) GetOutcomeDetail(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)

	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID pencapaian pembelajaran tidak valid"))
		return
	}

	res, err := h.svc.GetOutcomeDetail(c.Request.Context(), scope, id)
	if err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    res,
	})
}

// CreateOutcome handles POST /api/teacher/learning-outcomes.
func (h *LearningOutcomeHandler) CreateOutcome(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)

	var payload teacherdto.LearningOutcomeCreatePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		handler.BindError(c, err)
		return
	}

	res, err := h.svc.CreateOutcome(c.Request.Context(), scope, &payload)
	if err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    res,
	})
}

// UpdateOutcome handles PUT /api/teacher/learning-outcomes/:id.
func (h *LearningOutcomeHandler) UpdateOutcome(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)

	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID pencapaian pembelajaran tidak valid"))
		return
	}

	var payload teacherdto.LearningOutcomeUpdatePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		handler.BindError(c, err)
		return
	}

	res, err := h.svc.UpdateOutcome(c.Request.Context(), scope, id, &payload)
	if err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    res,
	})
}

// DeleteOutcome handles DELETE /api/teacher/learning-outcomes/:id.
func (h *LearningOutcomeHandler) DeleteOutcome(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)

	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID pencapaian pembelajaran tidak valid"))
		return
	}

	if err := h.svc.DeleteOutcome(c.Request.Context(), scope, id); err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Pencapaian pembelajaran berhasil dihapus",
	})
}

// MapAssignmentOutcome handles PUT /api/teacher/assignments/:id/learning-outcome.
func (h *LearningOutcomeHandler) MapAssignmentOutcome(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)

	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID tugas tidak valid"))
		return
	}

	var payload teacherdto.LearningOutcomeItemMappingPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		handler.BindError(c, err)
		return
	}

	if err := h.svc.MapAssignmentOutcome(c.Request.Context(), scope, id, payload.LearningOutcomeID); err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Pemetaan pencapaian pembelajaran berhasil diperbarui",
	})
}

// MapProjectOutcome handles PUT /api/teacher/projects/:id/learning-outcome.
func (h *LearningOutcomeHandler) MapProjectOutcome(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)

	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID projek tidak valid"))
		return
	}

	var payload teacherdto.LearningOutcomeItemMappingPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		handler.BindError(c, err)
		return
	}

	if err := h.svc.MapProjectOutcome(c.Request.Context(), scope, id, payload.LearningOutcomeID); err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Pemetaan pencapaian pembelajaran berhasil diperbarui",
	})
}

// MapAssessmentOutcome handles PUT /api/teacher/assessments/:id/learning-outcome.
func (h *LearningOutcomeHandler) MapAssessmentOutcome(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)

	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID asesmen tidak valid"))
		return
	}

	var payload teacherdto.LearningOutcomeItemMappingPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		handler.BindError(c, err)
		return
	}

	if err := h.svc.MapAssessmentOutcome(c.Request.Context(), scope, id, payload.LearningOutcomeID); err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Pemetaan pencapaian pembelajaran berhasil diperbarui",
	})
}
