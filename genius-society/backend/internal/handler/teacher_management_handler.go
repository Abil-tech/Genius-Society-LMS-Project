package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/geniussociety/backend/internal/dto"
	"github.com/geniussociety/backend/internal/service"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type TeacherManagementHandler struct {
	svc *service.TeacherManagementService
}

func NewTeacherManagementHandler(svc *service.TeacherManagementService) *TeacherManagementHandler {
	return &TeacherManagementHandler{svc: svc}
}

// GET /api/teachers
func (h *TeacherManagementHandler) List(c *gin.Context) {
	search := c.Query("search")
	subject := c.Query("subject")
	status := c.Query("status")
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64)

	var isWalas *bool
	if walasStr := c.Query("isWalas"); walasStr != "" {
		val := walasStr == "true"
		isWalas = &val
	}

	items, total, err := h.svc.ListTeachers(c.Request.Context(), search, subject, status, isWalas, page, limit)
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

// GET /api/teachers/stats
func (h *TeacherManagementHandler) Stats(c *gin.Context) {
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

// GET /api/teachers/:id
func (h *TeacherManagementHandler) GetOne(c *gin.Context) {
	id, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "BAD_REQUEST", Message: "invalid teacher id"},
		})
		return
	}

	t, err := h.svc.GetTeacherDetail(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "NOT_FOUND", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true, Data: t})
}

// POST /api/teachers
func (h *TeacherManagementHandler) Create(c *gin.Context) {
	var req dto.CreateTeacherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "VALIDATION_ERROR", Message: err.Error()},
		})
		return
	}

	t, err := h.svc.CreateTeacher(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "BAD_REQUEST", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse{Success: true, Data: t})
}

// PUT /api/teachers/:id
func (h *TeacherManagementHandler) Update(c *gin.Context) {
	id, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "BAD_REQUEST", Message: "invalid teacher id"},
		})
		return
	}

	var req dto.UpdateTeacherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "VALIDATION_ERROR", Message: err.Error()},
		})
		return
	}

	if err := h.svc.UpdateTeacher(c.Request.Context(), id, req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "BAD_REQUEST", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Data:    gin.H{"message": "data guru berhasil diperbarui"},
	})
}

// DELETE /api/teachers/:id
func (h *TeacherManagementHandler) Delete(c *gin.Context) {
	id, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "BAD_REQUEST", Message: "invalid teacher id"},
		})
		return
	}

	if err := h.svc.DeleteTeacher(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "BAD_REQUEST", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Data:    gin.H{"message": "data guru dan akun terkait berhasil dihapus"},
	})
}

// POST /api/teachers/:id/account
func (h *TeacherManagementHandler) CreateAccount(c *gin.Context) {
	id, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "BAD_REQUEST", Message: "invalid teacher id"},
		})
		return
	}

	var req dto.CreateTeacherAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "VALIDATION_ERROR", Message: err.Error()},
		})
		return
	}

	if err := h.svc.CreateAccount(c.Request.Context(), id, req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "BAD_REQUEST", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse{
		Success: true,
		Data:    gin.H{"message": "akun login guru berhasil dibuat dan ditautkan"},
	})
}

// PATCH /api/teachers/:id/account/toggle
func (h *TeacherManagementHandler) ToggleAccount(c *gin.Context) {
	id, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "BAD_REQUEST", Message: "invalid teacher id"},
		})
		return
	}

	newStatus, err := h.svc.ToggleAccount(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "BAD_REQUEST", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Data:    gin.H{"isActive": newStatus, "message": "status akun berhasil diperbarui"},
	})
}
