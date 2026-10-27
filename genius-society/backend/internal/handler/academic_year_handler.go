package handler

import (
	"net/http"

	"github.com/geniussociety/backend/internal/dto"
	"github.com/geniussociety/backend/internal/middleware"
	"github.com/geniussociety/backend/internal/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type AcademicYearHandler struct {
	svc *service.AcademicYearService
}

func NewAcademicYearHandler(svc *service.AcademicYearService) *AcademicYearHandler {
	return &AcademicYearHandler{svc: svc}
}

// GET /api/academic-years
func (h *AcademicYearHandler) List(c *gin.Context) {
	schoolIDHex := middleware.GetSchoolID(c)
	schoolID, _ := bson.ObjectIDFromHex(schoolIDHex)

	years, err := h.svc.List(c.Request.Context(), schoolID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "INTERNAL_ERROR", Message: err.Error()}})
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true, Data: years})
}

// GET /api/academic-years/:id
func (h *AcademicYearHandler) GetOne(c *gin.Context) {
	schoolIDHex := middleware.GetSchoolID(c)
	schoolID, _ := bson.ObjectIDFromHex(schoolIDHex)

	ayID, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: "invalid academic year id"}})
		return
	}
	year, err := h.svc.GetOne(c.Request.Context(), schoolID, ayID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "NOT_FOUND", Message: "academic year not found"}})
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true, Data: year})
}

// POST /api/academic-years
func (h *AcademicYearHandler) Create(c *gin.Context) {
	var req dto.CreateAcademicYearDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "VALIDATION_ERROR", Message: err.Error()}})
		return
	}
	schoolIDHex := middleware.GetSchoolID(c)
	schoolID, _ := bson.ObjectIDFromHex(schoolIDHex)
	userIDHex := middleware.GetUserID(c)
	userID, _ := bson.ObjectIDFromHex(userIDHex)

	year, err := h.svc.Create(c.Request.Context(), schoolID, userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: err.Error()}})
		return
	}
	c.JSON(http.StatusCreated, dto.SuccessResponse{Success: true, Data: year})
}

// PUT /api/academic-years/:id
func (h *AcademicYearHandler) Update(c *gin.Context) {
	var req dto.UpdateAcademicYearDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "VALIDATION_ERROR", Message: err.Error()}})
		return
	}
	schoolIDHex := middleware.GetSchoolID(c)
	schoolID, _ := bson.ObjectIDFromHex(schoolIDHex)
	userIDHex := middleware.GetUserID(c)
	userID, _ := bson.ObjectIDFromHex(userIDHex)

	ayID, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: "invalid id"}})
		return
	}
	year, err := h.svc.Update(c.Request.Context(), schoolID, userID, ayID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: err.Error()}})
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true, Data: year})
}

// PATCH /api/academic-years/:id/activate
func (h *AcademicYearHandler) Activate(c *gin.Context) {
	schoolIDHex := middleware.GetSchoolID(c)
	schoolID, _ := bson.ObjectIDFromHex(schoolIDHex)
	userIDHex := middleware.GetUserID(c)
	userID, _ := bson.ObjectIDFromHex(userIDHex)

	ayID, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: "invalid id"}})
		return
	}
	if err := h.svc.Activate(c.Request.Context(), schoolID, userID, ayID); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: err.Error()}})
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true, Data: gin.H{"message": "tahun ajaran berhasil diaktifkan"}})
}

// PATCH /api/academic-years/:id/archive
func (h *AcademicYearHandler) Archive(c *gin.Context) {
	schoolIDHex := middleware.GetSchoolID(c)
	schoolID, _ := bson.ObjectIDFromHex(schoolIDHex)
	userIDHex := middleware.GetUserID(c)
	userID, _ := bson.ObjectIDFromHex(userIDHex)

	ayID, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: "invalid id"}})
		return
	}
	if err := h.svc.Archive(c.Request.Context(), schoolID, userID, ayID); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: err.Error()}})
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true, Data: gin.H{"message": "tahun ajaran berhasil diarsipkan"}})
}

// DELETE /api/academic-years/:id
func (h *AcademicYearHandler) Delete(c *gin.Context) {
	schoolIDHex := middleware.GetSchoolID(c)
	schoolID, _ := bson.ObjectIDFromHex(schoolIDHex)
	userIDHex := middleware.GetUserID(c)
	userID, _ := bson.ObjectIDFromHex(userIDHex)

	ayID, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: "invalid id"}})
		return
	}
	if err := h.svc.Delete(c.Request.Context(), schoolID, userID, ayID); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: err.Error()}})
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true, Data: gin.H{"message": "tahun ajaran berhasil dihapus"}})
}

// ============================================================================
// PROMOTION HANDLERS
// ============================================================================

// GET /api/academic-years/:id/promotion/candidates
func (h *AcademicYearHandler) GetPromotionCandidates(c *gin.Context) {
	schoolIDHex := middleware.GetSchoolID(c)
	schoolID, _ := bson.ObjectIDFromHex(schoolIDHex)

	sourceYearID, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: "invalid academic year id"}})
		return
	}

	candidates, err := h.svc.GetPromotionCandidates(c.Request.Context(), schoolID, sourceYearID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: err.Error()}})
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true, Data: candidates})
}

// POST /api/academic-years/:id/promotion/preview
func (h *AcademicYearHandler) PreviewPromotion(c *gin.Context) {
	schoolIDHex := middleware.GetSchoolID(c)
	schoolID, _ := bson.ObjectIDFromHex(schoolIDHex)

	var req dto.PromotionPreviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "VALIDATION_ERROR", Message: err.Error()}})
		return
	}

	// Use route param if not specified in body
	if req.SourceAcademicYearID == "" {
		req.SourceAcademicYearID = c.Param("id")
	}

	preview, err := h.svc.PreviewPromotion(c.Request.Context(), schoolID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: err.Error()}})
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true, Data: preview})
}

// POST /api/academic-years/:id/promotion/execute
func (h *AcademicYearHandler) ExecutePromotion(c *gin.Context) {
	schoolIDHex := middleware.GetSchoolID(c)
	schoolID, _ := bson.ObjectIDFromHex(schoolIDHex)
	userIDHex := middleware.GetUserID(c)
	userID, _ := bson.ObjectIDFromHex(userIDHex)

	var req dto.PromotionExecutionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "VALIDATION_ERROR", Message: err.Error()}})
		return
	}

	if req.SourceAcademicYearID == "" {
		req.SourceAcademicYearID = c.Param("id")
	}

	res, err := h.svc.ExecutePromotion(c.Request.Context(), schoolID, userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: err.Error()}})
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true, Data: res})
}

// GET /api/academic-years/:id/promotion/history
func (h *AcademicYearHandler) GetPromotionHistory(c *gin.Context) {
	schoolIDHex := middleware.GetSchoolID(c)
	schoolID, _ := bson.ObjectIDFromHex(schoolIDHex)

	sourceYearID, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: "invalid academic year id"}})
		return
	}

	history, err := h.svc.GetPromotionHistory(c.Request.Context(), schoolID, sourceYearID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: err.Error()}})
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true, Data: history})
}

// ============================================================================
// GRADUATION HANDLERS
// ============================================================================

// GET /api/academic-years/:id/graduation/candidates
func (h *AcademicYearHandler) GetGraduationCandidates(c *gin.Context) {
	schoolIDHex := middleware.GetSchoolID(c)
	schoolID, _ := bson.ObjectIDFromHex(schoolIDHex)

	ayID, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: "invalid academic year id"}})
		return
	}

	candidates, err := h.svc.GetGraduationCandidates(c.Request.Context(), schoolID, ayID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: err.Error()}})
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true, Data: candidates})
}

// POST /api/academic-years/:id/graduation/preview
func (h *AcademicYearHandler) PreviewGraduation(c *gin.Context) {
	schoolIDHex := middleware.GetSchoolID(c)
	schoolID, _ := bson.ObjectIDFromHex(schoolIDHex)

	var req dto.GraduationPreviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "VALIDATION_ERROR", Message: err.Error()}})
		return
	}

	if req.AcademicYearID == "" {
		req.AcademicYearID = c.Param("id")
	}

	preview, err := h.svc.PreviewGraduation(c.Request.Context(), schoolID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: err.Error()}})
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true, Data: preview})
}

// POST /api/academic-years/:id/graduation/execute
func (h *AcademicYearHandler) ExecuteGraduation(c *gin.Context) {
	schoolIDHex := middleware.GetSchoolID(c)
	schoolID, _ := bson.ObjectIDFromHex(schoolIDHex)
	userIDHex := middleware.GetUserID(c)
	userID, _ := bson.ObjectIDFromHex(userIDHex)

	var req dto.GraduationExecutionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "VALIDATION_ERROR", Message: err.Error()}})
		return
	}

	if req.AcademicYearID == "" {
		req.AcademicYearID = c.Param("id")
	}

	res, err := h.svc.ExecuteGraduation(c.Request.Context(), schoolID, userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: err.Error()}})
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true, Data: res})
}

// GET /api/academic-years/:id/graduation/history
func (h *AcademicYearHandler) GetGraduationHistory(c *gin.Context) {
	schoolIDHex := middleware.GetSchoolID(c)
	schoolID, _ := bson.ObjectIDFromHex(schoolIDHex)

	ayID, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: "invalid academic year id"}})
		return
	}

	history, err := h.svc.GetGraduationHistory(c.Request.Context(), schoolID, ayID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: err.Error()}})
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true, Data: history})
}

// ============================================================================
// SUSPENDED DECISION HANDLERS
// ============================================================================

// GET /api/academic-years/:id/decisions/suspended
func (h *AcademicYearHandler) GetSuspendedDecisions(c *gin.Context) {
	schoolIDHex := middleware.GetSchoolID(c)
	schoolID, _ := bson.ObjectIDFromHex(schoolIDHex)

	ayID, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: "invalid academic year id"}})
		return
	}

	decisions, err := h.svc.GetSuspendedDecisions(c.Request.Context(), schoolID, ayID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: err.Error()}})
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true, Data: decisions})
}

// GET /api/academic-decisions/:decisionId
func (h *AcademicYearHandler) GetDecisionDetail(c *gin.Context) {
	schoolIDHex := middleware.GetSchoolID(c)
	schoolID, _ := bson.ObjectIDFromHex(schoolIDHex)

	decID, err := bson.ObjectIDFromHex(c.Param("decisionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: "invalid decision id"}})
		return
	}

	decision, err := h.svc.GetDecisionDetail(c.Request.Context(), schoolID, decID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "NOT_FOUND", Message: err.Error()}})
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true, Data: decision})
}

// POST /api/academic-decisions/:decisionId/review
func (h *AcademicYearHandler) ReviewDecision(c *gin.Context) {
	schoolIDHex := middleware.GetSchoolID(c)
	schoolID, _ := bson.ObjectIDFromHex(schoolIDHex)
	userIDHex := middleware.GetUserID(c)
	userID, _ := bson.ObjectIDFromHex(userIDHex)

	decID, err := bson.ObjectIDFromHex(c.Param("decisionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: "invalid decision id"}})
		return
	}

	var req dto.ReviewDecisionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "VALIDATION_ERROR", Message: err.Error()}})
		return
	}

	if err := h.svc.ReviewDecision(c.Request.Context(), schoolID, userID, decID, req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: err.Error()}})
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true, Data: gin.H{"message": "peninjauan berhasil dicatat"}})
}

// POST /api/academic-decisions/:decisionId/resolve
func (h *AcademicYearHandler) ResolveDecision(c *gin.Context) {
	schoolIDHex := middleware.GetSchoolID(c)
	schoolID, _ := bson.ObjectIDFromHex(schoolIDHex)
	userIDHex := middleware.GetUserID(c)
	userID, _ := bson.ObjectIDFromHex(userIDHex)

	decID, err := bson.ObjectIDFromHex(c.Param("decisionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: "invalid decision id"}})
		return
	}

	var req dto.ResolveDecisionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "VALIDATION_ERROR", Message: err.Error()}})
		return
	}

	if err := h.svc.ResolveDecision(c.Request.Context(), schoolID, userID, decID, req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: err.Error()}})
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true, Data: gin.H{"message": "keputusan berhasil diselesaikan"}})
}
