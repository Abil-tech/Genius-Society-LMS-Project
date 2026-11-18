package teacherhandler

import (
	"net/http"
	"strings"

	"github.com/geniussociety/backend/internal/dto"
	teacherdto "github.com/geniussociety/backend/internal/dto/teacher"
	"github.com/geniussociety/backend/internal/handler"
	"github.com/geniussociety/backend/internal/middleware"
	teachersvc "github.com/geniussociety/backend/internal/service/teacher"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// GradeHandler handles HTTP requests for teacher grade endpoints.
type GradeHandler struct {
	svc teachersvc.GradeService
}

// NewGradeHandler constructs a new GradeHandler.
func NewGradeHandler(svc teachersvc.GradeService) *GradeHandler {
	return &GradeHandler{svc: svc}
}

// GetGrades handles GET /api/teacher/grades?classId=...&subjectId=...
func (h *GradeHandler) GetGrades(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)

	classIDStr := strings.TrimSpace(c.Query("classId"))
	subjectIDStr := strings.TrimSpace(c.Query("subjectId"))

	if classIDStr == "" || subjectIDStr == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error: dto.ErrorBody{
				Code:    "VALIDATION_ERROR",
				Message: "Parameter classId dan subjectId wajib disertakan",
			},
		})
		return
	}

	classID, err := bson.ObjectIDFromHex(classIDStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID kelas tidak valid"))
		return
	}

	subjectID, err := bson.ObjectIDFromHex(subjectIDStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID mata pelajaran tidak valid"))
		return
	}

	grades, err := h.svc.GetGrades(c.Request.Context(), scope, classID, subjectID)
	if err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    grades,
	})
}

// UpdateGrade handles PUT /api/teacher/grades
func (h *GradeHandler) UpdateGrade(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)

	var payload teacherdto.UpdateStudentGradePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		handler.BindError(c, err)
		return
	}

	res, err := h.svc.UpdateGrade(c.Request.Context(), scope, payload)
	if err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    res,
	})
}
