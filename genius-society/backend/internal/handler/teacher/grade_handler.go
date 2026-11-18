package teacherhandler

import (
	"net/http"
	"strconv"
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

// InputSTSExams handles PUT /api/teacher/grades/exams
func (h *GradeHandler) InputSTSExams(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)

	var payload teacherdto.BulkExamsPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		handler.BindError(c, err)
		return
	}

	if err := h.svc.InputSTSExams(c.Request.Context(), scope, &payload); err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Nilai ujian berhasil disimpan",
	})
}

// SetManualAssignmentGrade handles PUT /api/teacher/assignments/:id/students/:studentId/grade
func (h *GradeHandler) SetManualAssignmentGrade(c *gin.Context) {
	h.setManualGrade(c, "assignment")
}

// DeleteManualAssignmentGrade handles DELETE /api/teacher/assignments/:id/students/:studentId/grade
func (h *GradeHandler) DeleteManualAssignmentGrade(c *gin.Context) {
	h.deleteManualGrade(c, "assignment")
}

// SetManualProjectGrade handles PUT /api/teacher/projects/:id/students/:studentId/grade
func (h *GradeHandler) SetManualProjectGrade(c *gin.Context) {
	h.setManualGrade(c, "project")
}

// DeleteManualProjectGrade handles DELETE /api/teacher/projects/:id/students/:studentId/grade
func (h *GradeHandler) DeleteManualProjectGrade(c *gin.Context) {
	h.deleteManualGrade(c, "project")
}

// SetManualAssessmentGrade handles PUT /api/teacher/assessments/:id/students/:studentId/grade
func (h *GradeHandler) SetManualAssessmentGrade(c *gin.Context) {
	h.setManualGrade(c, "assessment")
}

// DeleteManualAssessmentGrade handles DELETE /api/teacher/assessments/:id/students/:studentId/grade
func (h *GradeHandler) DeleteManualAssessmentGrade(c *gin.Context) {
	h.deleteManualGrade(c, "assessment")
}

func (h *GradeHandler) setManualGrade(c *gin.Context, itemType string) {
	scope := middleware.MustGetTeacherScope(c)

	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID item tidak valid"))
		return
	}

	studentIDStr := strings.TrimSpace(c.Param("studentId"))
	studentID, err := bson.ObjectIDFromHex(studentIDStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID siswa tidak valid"))
		return
	}

	var payload teacherdto.ManualGradePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		handler.BindError(c, err)
		return
	}

	if err := h.svc.SetManualGrade(c.Request.Context(), scope, itemType, id, studentID, payload.Score, payload.Feedback, payload.Reason); err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Nilai manual berhasil disimpan",
	})
}

func (h *GradeHandler) deleteManualGrade(c *gin.Context, itemType string) {
	scope := middleware.MustGetTeacherScope(c)

	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID item tidak valid"))
		return
	}

	studentIDStr := strings.TrimSpace(c.Param("studentId"))
	studentID, err := bson.ObjectIDFromHex(studentIDStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID siswa tidak valid"))
		return
	}

	if err := h.svc.DeleteManualGrade(c.Request.Context(), scope, itemType, id, studentID); err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Nilai manual berhasil dihapus",
	})
}

// GetGradeRecap handles GET /api/teacher/grades/recap?classId=...&subjectId=...&semester=...
func (h *GradeHandler) GetGradeRecap(c *gin.Context) {
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

	semester := 0
	if sStr := strings.TrimSpace(c.Query("semester")); sStr != "" {
		if s, err := strconv.Atoi(sStr); err == nil {
			semester = s
		}
	}

	recap, err := h.svc.GetGradeRecap(c.Request.Context(), scope, classID, subjectID, semester)
	if err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    recap,
	})
}
