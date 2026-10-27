package handler

import (
	"net/http"
	"strconv"

	"github.com/geniussociety/backend/internal/dto"
	"github.com/geniussociety/backend/internal/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type SubjectManagementHandler struct {
	svc *service.SubjectManagementService
}

func NewSubjectManagementHandler(svc *service.SubjectManagementService) *SubjectManagementHandler {
	return &SubjectManagementHandler{svc: svc}
}

// GET /api/subjects
func (h *SubjectManagementHandler) List(c *gin.Context) {
	search := c.Query("search")
	category := c.Query("category")
	gradeLevel := c.Query("gradeLevel")
	status := c.Query("status")
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64)

	items, total, err := h.svc.ListSubjects(c.Request.Context(), search, category, gradeLevel, status, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "INTERNAL_ERROR", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Data:    items,
		Meta:    dto.PaginationMeta{Page: page, Limit: limit, Total: total},
	})
}

// GET /api/subjects/stats
func (h *SubjectManagementHandler) Stats(c *gin.Context) {
	stats, err := h.svc.GetStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "INTERNAL_ERROR", Message: err.Error()},
		})
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true, Data: stats})
}

// GET /api/subjects/:id
func (h *SubjectManagementHandler) GetOne(c *gin.Context) {
	id, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "BAD_REQUEST", Message: "invalid subject id"},
		})
		return
	}

	sub, err := h.svc.GetSubjectDetail(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "NOT_FOUND", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true, Data: sub})
}

// POST /api/subjects
func (h *SubjectManagementHandler) Create(c *gin.Context) {
	var req dto.CreateSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "VALIDATION_ERROR", Message: err.Error()},
		})
		return
	}

	sub, err := h.svc.CreateSubject(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "BAD_REQUEST", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse{Success: true, Data: sub})
}

// PUT /api/subjects/:id
func (h *SubjectManagementHandler) Update(c *gin.Context) {
	id, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "BAD_REQUEST", Message: "invalid subject id"},
		})
		return
	}

	var req dto.UpdateSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "VALIDATION_ERROR", Message: err.Error()},
		})
		return
	}

	if err := h.svc.UpdateSubject(c.Request.Context(), id, req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "BAD_REQUEST", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Data:    gin.H{"message": "mata pelajaran berhasil diperbarui"},
	})
}

// DELETE /api/subjects/:id
func (h *SubjectManagementHandler) Delete(c *gin.Context) {
	id, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "BAD_REQUEST", Message: "invalid subject id"},
		})
		return
	}

	if err := h.svc.DeleteSubject(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "BAD_REQUEST", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Data:    gin.H{"message": "mata pelajaran berhasil dihapus"},
	})
}
