package teacher

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/geniussociety/backend/internal/handler"
	"github.com/geniussociety/backend/internal/middleware"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// AssessmentHandler handles HTTP requests for teacher assessment endpoints.
type AssessmentHandler struct {
	svc AssessmentService
}

// NewAssessmentHandler constructs a new AssessmentHandler instance.
func NewAssessmentHandler(svc AssessmentService) *AssessmentHandler {
	return &AssessmentHandler{svc: svc}
}

// ListAssessments handles GET /api/teacher/assessments.
func (h *AssessmentHandler) ListAssessments(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)

	classID := strings.TrimSpace(c.Query("classId"))
	subjectID := strings.TrimSpace(c.Query("subjectId"))
	search := strings.TrimSpace(c.Query("q"))
	if search == "" {
		search = strings.TrimSpace(c.Query("search"))
	}

	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64)

	list, meta, err := h.svc.ListAssessments(c.Request.Context(), scope, classID, subjectID, search, page, limit)
	if err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    list,
		"meta":    meta,
	})
}

// GetAssessmentDetail handles GET /api/teacher/assessments/:id.
func (h *AssessmentHandler) GetAssessmentDetail(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)

	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID assessment tidak valid"))
		return
	}

	as, err := h.svc.GetAssessmentDetail(c.Request.Context(), scope, id)
	if err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    as,
	})
}

// CreateAssessment handles POST /api/teacher/assessments.
func (h *AssessmentHandler) CreateAssessment(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)

	var payload AssessmentCreatePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		handler.BindError(c, err)
		return
	}

	res, err := h.svc.CreateAssessment(c.Request.Context(), scope, &payload)
	if err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    res,
	})
}

// UpdateAssessment handles PUT /api/teacher/assessments/:id.
func (h *AssessmentHandler) UpdateAssessment(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)

	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID assessment tidak valid"))
		return
	}

	var payload AssessmentUpdatePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		handler.BindError(c, err)
		return
	}

	res, err := h.svc.UpdateAssessment(c.Request.Context(), scope, id, &payload)
	if err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    res,
	})
}

// DeleteAssessment handles DELETE /api/teacher/assessments/:id.
func (h *AssessmentHandler) DeleteAssessment(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)

	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID assessment tidak valid"))
		return
	}

	if err := h.svc.DeleteAssessment(c.Request.Context(), scope, id); err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Assessment berhasil dihapus",
	})
}

// AddQuestion handles POST /api/teacher/assessments/:id/questions.
func (h *AssessmentHandler) AddQuestion(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)

	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID assessment tidak valid"))
		return
	}

	var payload QuestionCreateUpdatePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		handler.BindError(c, err)
		return
	}

	q, err := h.svc.AddQuestion(c.Request.Context(), scope, id, &payload)
	if err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    q,
	})
}

// UpdateQuestion handles PUT /api/teacher/assessments/:id/questions/:questionId.
func (h *AssessmentHandler) UpdateQuestion(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)

	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID assessment tidak valid"))
		return
	}

	qIDStr := strings.TrimSpace(c.Param("questionId"))
	qID, err := bson.ObjectIDFromHex(qIDStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID soal tidak valid"))
		return
	}

	var payload QuestionCreateUpdatePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		handler.BindError(c, err)
		return
	}

	q, err := h.svc.UpdateQuestion(c.Request.Context(), scope, id, qID, &payload)
	if err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    q,
	})
}

// DeleteQuestion handles DELETE /api/teacher/assessments/:id/questions/:questionId.
func (h *AssessmentHandler) DeleteQuestion(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)

	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID assessment tidak valid"))
		return
	}

	qIDStr := strings.TrimSpace(c.Param("questionId"))
	qID, err := bson.ObjectIDFromHex(qIDStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID soal tidak valid"))
		return
	}

	if err := h.svc.DeleteQuestion(c.Request.Context(), scope, id, qID); err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Soal berhasil dihapus",
	})
}

// ReorderQuestions handles PUT /api/teacher/assessments/:id/questions/reorder.
func (h *AssessmentHandler) ReorderQuestions(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)

	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID assessment tidak valid"))
		return
	}

	var payload QuestionReorderPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		handler.BindError(c, err)
		return
	}

	if err := h.svc.ReorderQuestions(c.Request.Context(), scope, id, &payload); err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Urutan soal berhasil diperbarui",
	})
}

// RegradeMCQ handles POST /api/teacher/assessments/:id/regrade-mcq.
func (h *AssessmentHandler) RegradeMCQ(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)

	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID assessment tidak valid"))
		return
	}

	var payload RegradeMCQPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		handler.BindError(c, err)
		return
	}

	if err := h.svc.RegradeMCQ(c.Request.Context(), scope, id, &payload); err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Kunci jawaban berhasil diperbarui dan seluruh pengerjaan telah dihitung ulang",
	})
}

// PublishAssessment handles POST /api/teacher/assessments/:id/publish.
func (h *AssessmentHandler) PublishAssessment(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)

	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID assessment tidak valid"))
		return
	}

	res, err := h.svc.PublishAssessment(c.Request.Context(), scope, id)
	if err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    res,
	})
}

// CloseAssessment handles POST /api/teacher/assessments/:id/close.
func (h *AssessmentHandler) CloseAssessment(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)

	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID assessment tidak valid"))
		return
	}

	res, err := h.svc.CloseAssessment(c.Request.Context(), scope, id)
	if err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    res,
	})
}

// ListAttempts handles GET /api/teacher/assessments/:id/attempts.
func (h *AssessmentHandler) ListAttempts(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)

	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID assessment tidak valid"))
		return
	}

	status := strings.TrimSpace(c.Query("status"))
	classID := strings.TrimSpace(c.Query("classId"))
	search := strings.TrimSpace(c.Query("q"))
	if search == "" {
		search = strings.TrimSpace(c.Query("search"))
	}

	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64)

	attempts, meta, err := h.svc.ListAttempts(c.Request.Context(), scope, id, status, classID, search, page, limit)
	if err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    attempts,
		"meta":    meta,
	})
}

// GetAttemptDetail handles GET /api/teacher/assessments/:id/attempts/:attemptId.
func (h *AssessmentHandler) GetAttemptDetail(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)

	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID assessment tidak valid"))
		return
	}

	attemptIDStr := strings.TrimSpace(c.Param("attemptId"))
	attemptID, err := bson.ObjectIDFromHex(attemptIDStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID attempt tidak valid"))
		return
	}

	res, err := h.svc.GetAttemptDetail(c.Request.Context(), scope, id, attemptID)
	if err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    res,
	})
}

// GradeEssayAnswer handles PUT /api/teacher/assessments/:id/answers/:answerId/grade.
func (h *AssessmentHandler) GradeEssayAnswer(c *gin.Context) {
	scope := middleware.MustGetTeacherScope(c)

	idStr := strings.TrimSpace(c.Param("id"))
	id, err := bson.ObjectIDFromHex(idStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID assessment tidak valid"))
		return
	}

	answerIDStr := strings.TrimSpace(c.Param("answerId"))
	answerID, err := bson.ObjectIDFromHex(answerIDStr)
	if err != nil {
		handler.RespondError(c, handler.NewValidationError("ID jawaban tidak valid"))
		return
	}

	var payload GradeEssayAnswerPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		handler.BindError(c, err)
		return
	}

	if err := h.svc.GradeEssayAnswer(c.Request.Context(), scope, id, answerID, payload.Score, payload.Feedback); err != nil {
		handler.RespondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Nilai essay berhasil disimpan",
	})
}

// StartAssessment handles POST /api/teacher/assessments/:id/start.
// Rule: POST /start adalah milik Murid; Guru 403.
func (h *AssessmentHandler) StartAssessment(c *gin.Context) {
	handler.RespondError(c, handler.NewForbiddenError("Pengerjaan assessment hanya dapat dilakukan oleh murid"))
}

// SubmitAssessment handles POST /api/teacher/assessments/:id/submit.
// Rule: POST /submit adalah milik Murid; Guru 403.
func (h *AssessmentHandler) SubmitAssessment(c *gin.Context) {
	handler.RespondError(c, handler.NewForbiddenError("Pengumpulan assessment hanya dapat dilakukan oleh murid"))
}
