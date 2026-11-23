package studenthandler

import (
	"net/http"
	"strconv"
	"strings"

	studentdto "github.com/geniussociety/backend/internal/dto/student"
	"github.com/geniussociety/backend/internal/handler"
	"github.com/geniussociety/backend/internal/middleware"
	studentsvc "github.com/geniussociety/backend/internal/service/student"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// StudentHandler handles incoming HTTP requests for the student portal.
type StudentHandler struct {
	svc studentsvc.StudentService
}

// NewStudentHandler constructs a StudentHandler instance.
func NewStudentHandler(svc studentsvc.StudentService) *StudentHandler {
	return &StudentHandler{svc: svc}
}

// GetDashboard handles GET /api/student/dashboard
func (h *StudentHandler) GetDashboard(c *gin.Context) {
	scope := middleware.MustGetStudentScope(c)
	dto, err := h.svc.GetDashboard(c.Request.Context(), scope)
	if err != nil {
		handler.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": dto,
	})
}

// -------------------------------------------------------------
// Materials
// -------------------------------------------------------------

// ListMaterials handles GET /api/student/materials
func (h *StudentHandler) ListMaterials(c *gin.Context) {
	scope := middleware.MustGetStudentScope(c)
	subjectID := strings.TrimSpace(c.Query("subjectId"))
	list, err := h.svc.ListMaterials(c.Request.Context(), scope, subjectID)
	if err != nil {
		handler.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": list,
	})
}

// GetMaterialDetail handles GET /api/student/materials/:id
func (h *StudentHandler) GetMaterialDetail(c *gin.Context) {
	scope := middleware.MustGetStudentScope(c)
	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID materi tidak valid"))
		return
	}

	detail, err := h.svc.GetMaterialDetail(c.Request.Context(), scope, id)
	if err != nil {
		handler.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": detail,
	})
}

// DownloadMaterialFile handles GET /api/student/materials/:id/files/:index/download
func (h *StudentHandler) DownloadMaterialFile(c *gin.Context) {
	scope := middleware.MustGetStudentScope(c)
	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID materi tidak valid"))
		return
	}

	indexStr := strings.TrimSpace(c.Param("index"))
	index, err := strconv.Atoi(indexStr)
	if err != nil || index < 0 {
		handler.RespondError(c, handler.NewValidationError("Index berkas tidak valid"))
		return
	}

	url, err := h.svc.GetMaterialDownloadURL(c.Request.Context(), scope, id, index)
	if err != nil {
		handler.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"url":       url,
			"expiresIn": 300,
		},
	})
}

// -------------------------------------------------------------
// Assignments
// -------------------------------------------------------------

// ListAssignments handles GET /api/student/assignments
func (h *StudentHandler) ListAssignments(c *gin.Context) {
	scope := middleware.MustGetStudentScope(c)
	subjectID := strings.TrimSpace(c.Query("subjectId"))
	list, err := h.svc.ListAssignments(c.Request.Context(), scope, subjectID)
	if err != nil {
		handler.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": list,
	})
}

// GetAssignmentDetail handles GET /api/student/assignments/:id
func (h *StudentHandler) GetAssignmentDetail(c *gin.Context) {
	scope := middleware.MustGetStudentScope(c)
	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID tugas tidak valid"))
		return
	}

	detail, err := h.svc.GetAssignmentDetail(c.Request.Context(), scope, id)
	if err != nil {
		handler.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": detail,
	})
}

// SubmitAssignment handles POST /api/student/assignments/:id/submit
func (h *StudentHandler) SubmitAssignment(c *gin.Context) {
	h.processAssignmentSubmission(c, false)
}

// ResubmitAssignment handles PUT /api/student/assignments/:id/resubmit
func (h *StudentHandler) ResubmitAssignment(c *gin.Context) {
	h.processAssignmentSubmission(c, true)
}

func (h *StudentHandler) processAssignmentSubmission(c *gin.Context, isResubmit bool) {
	scope := middleware.MustGetStudentScope(c)
	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID tugas tidak valid"))
		return
	}

	file, _ := c.FormFile("file")
	note := strings.TrimSpace(c.PostForm("note"))

	sub, err := h.svc.SubmitAssignment(c.Request.Context(), scope, id, file, note, isResubmit)
	if err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tugas berhasil dikumpulkan",
		"data":    sub,
	})
}

// DownloadAssignmentFile handles GET /api/student/assignments/:id/files/:index/download
func (h *StudentHandler) DownloadAssignmentFile(c *gin.Context) {
	scope := middleware.MustGetStudentScope(c)
	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID tugas tidak valid"))
		return
	}

	indexStr := strings.TrimSpace(c.Param("index"))
	index, err := strconv.Atoi(indexStr)
	if err != nil || index < 0 {
		handler.RespondError(c, handler.NewValidationError("Index berkas tidak valid"))
		return
	}

	url, err := h.svc.GetAssignmentFileDownloadURL(c.Request.Context(), scope, id, index)
	if err != nil {
		handler.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"url":       url,
			"expiresIn": 300,
		},
	})
}

// -------------------------------------------------------------
// Projects
// -------------------------------------------------------------

// ListProjects handles GET /api/student/projects
func (h *StudentHandler) ListProjects(c *gin.Context) {
	scope := middleware.MustGetStudentScope(c)
	subjectID := strings.TrimSpace(c.Query("subjectId"))
	list, err := h.svc.ListProjects(c.Request.Context(), scope, subjectID)
	if err != nil {
		handler.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": list,
	})
}

// GetProjectDetail handles GET /api/student/projects/:id
func (h *StudentHandler) GetProjectDetail(c *gin.Context) {
	scope := middleware.MustGetStudentScope(c)
	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID projek tidak valid"))
		return
	}

	detail, err := h.svc.GetProjectDetail(c.Request.Context(), scope, id)
	if err != nil {
		handler.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": detail,
	})
}

// SubmitProject handles POST /api/student/projects/:id/submit
func (h *StudentHandler) SubmitProject(c *gin.Context) {
	h.processProjectSubmission(c, false)
}

// ResubmitProject handles PUT /api/student/projects/:id/resubmit
func (h *StudentHandler) ResubmitProject(c *gin.Context) {
	h.processProjectSubmission(c, true)
}

func (h *StudentHandler) processProjectSubmission(c *gin.Context, isResubmit bool) {
	scope := middleware.MustGetStudentScope(c)
	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID projek tidak valid"))
		return
	}

	file, _ := c.FormFile("file")
	sub, err := h.svc.SubmitProject(c.Request.Context(), scope, id, file, isResubmit)
	if err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Hasil projek berhasil dikumpulkan",
		"data":    sub,
	})
}

// DownloadProjectFile handles GET /api/student/projects/:id/files/:index/download
func (h *StudentHandler) DownloadProjectFile(c *gin.Context) {
	scope := middleware.MustGetStudentScope(c)
	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID projek tidak valid"))
		return
	}

	indexStr := strings.TrimSpace(c.Param("index"))
	index, err := strconv.Atoi(indexStr)
	if err != nil || index < 0 {
		handler.RespondError(c, handler.NewValidationError("Index berkas tidak valid"))
		return
	}

	url, err := h.svc.GetProjectFileDownloadURL(c.Request.Context(), scope, id, index)
	if err != nil {
		handler.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"url":       url,
			"expiresIn": 300,
		},
	})
}

// -------------------------------------------------------------
// Assessments
// -------------------------------------------------------------

// ListAssessments handles GET /api/student/assessments
func (h *StudentHandler) ListAssessments(c *gin.Context) {
	scope := middleware.MustGetStudentScope(c)
	subjectID := strings.TrimSpace(c.Query("subjectId"))
	list, err := h.svc.ListAssessments(c.Request.Context(), scope, subjectID)
	if err != nil {
		handler.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": list,
	})
}

// GetAssessmentDetail handles GET /api/student/assessments/:id
func (h *StudentHandler) GetAssessmentDetail(c *gin.Context) {
	scope := middleware.MustGetStudentScope(c)
	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID assessment tidak valid"))
		return
	}

	detail, err := h.svc.GetAssessmentDetail(c.Request.Context(), scope, id)
	if err != nil {
		handler.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": detail,
	})
}

// StartAssessment handles POST /api/student/assessments/:id/start
func (h *StudentHandler) StartAssessment(c *gin.Context) {
	scope := middleware.MustGetStudentScope(c)
	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID assessment tidak valid"))
		return
	}

	res, err := h.svc.StartAssessment(c.Request.Context(), scope, id)
	if err != nil {
		handler.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Assessment berhasil dimulai",
		"data":    res,
	})
}

// SubmitAssessment handles POST /api/student/assessments/:id/submit
func (h *StudentHandler) SubmitAssessment(c *gin.Context) {
	scope := middleware.MustGetStudentScope(c)
	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID assessment tidak valid"))
		return
	}

	var payload studentdto.AssessmentSubmitPayload
	_ = c.ShouldBindJSON(&payload)

	res, err := h.svc.SubmitAssessment(c.Request.Context(), scope, id, &payload)
	if err != nil {
		handler.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Assessment berhasil dikumpulkan",
		"data":    res,
	})
}

// GetAssessmentResult handles GET /api/student/assessments/:id/result
func (h *StudentHandler) GetAssessmentResult(c *gin.Context) {
	scope := middleware.MustGetStudentScope(c)
	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID assessment tidak valid"))
		return
	}

	res, err := h.svc.GetAssessmentResult(c.Request.Context(), scope, id)
	if err != nil {
		handler.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": res,
	})
}

// -------------------------------------------------------------
// Grades
// -------------------------------------------------------------

// GetGrades handles GET /api/student/grades
func (h *StudentHandler) GetGrades(c *gin.Context) {
	scope := middleware.MustGetStudentScope(c)
	grades, err := h.svc.GetGrades(c.Request.Context(), scope)
	if err != nil {
		handler.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": grades,
	})
}

// GetSubjectGradeDetail handles GET /api/student/grades/:subjectId
func (h *StudentHandler) GetSubjectGradeDetail(c *gin.Context) {
	scope := middleware.MustGetStudentScope(c)
	subjectIDStr := strings.TrimSpace(c.Param("subjectId"))
	subjectID, err := bson.ObjectIDFromHex(subjectIDStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID mata pelajaran tidak valid"))
		return
	}

	detail, err := h.svc.GetSubjectGradeDetail(c.Request.Context(), scope, subjectID)
	if err != nil {
		handler.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": detail,
	})
}

// -------------------------------------------------------------
// Schedule & Ranking
// -------------------------------------------------------------

// GetSchedule handles GET /api/student/schedule
func (h *StudentHandler) GetSchedule(c *gin.Context) {
	scope := middleware.MustGetStudentScope(c)
	schedule, err := h.svc.GetSchedule(c.Request.Context(), scope)
	if err != nil {
		handler.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": schedule,
	})
}

// GetRanking handles GET /api/student/ranking
func (h *StudentHandler) GetRanking(c *gin.Context) {
	scope := middleware.MustGetStudentScope(c)
	ranking, err := h.svc.GetRanking(c.Request.Context(), scope)
	if err != nil {
		handler.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": ranking,
	})
}
