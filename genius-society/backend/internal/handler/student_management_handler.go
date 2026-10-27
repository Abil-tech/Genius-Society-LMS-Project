package handler

import (
	"net/http"
	"strconv"

	"github.com/geniussociety/backend/internal/dto"
	"github.com/geniussociety/backend/internal/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type StudentManagementHandler struct {
	svc *service.StudentManagementService
}

func NewStudentManagementHandler(svc *service.StudentManagementService) *StudentManagementHandler {
	return &StudentManagementHandler{svc: svc}
}

// GET /api/students
func (h *StudentManagementHandler) List(c *gin.Context) {
	search := c.Query("search")
	classID := c.Query("classId")
	status := c.Query("status")
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64)

	items, total, err := h.svc.ListStudents(c.Request.Context(), search, classID, status, page, limit)
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

// GET /api/students/stats
func (h *StudentManagementHandler) Stats(c *gin.Context) {
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

// GET /api/students/:id
func (h *StudentManagementHandler) GetOne(c *gin.Context) {
	id, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "BAD_REQUEST", Message: "invalid student id"},
		})
		return
	}

	st, err := h.svc.GetStudentDetail(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "NOT_FOUND", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true, Data: st})
}

// POST /api/students
func (h *StudentManagementHandler) Create(c *gin.Context) {
	var req dto.CreateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "VALIDATION_ERROR", Message: err.Error()},
		})
		return
	}

	st, err := h.svc.CreateStudent(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "BAD_REQUEST", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse{Success: true, Data: st})
}

// PUT /api/students/:id
func (h *StudentManagementHandler) Update(c *gin.Context) {
	id, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "BAD_REQUEST", Message: "invalid student id"},
		})
		return
	}

	var req dto.UpdateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "VALIDATION_ERROR", Message: err.Error()},
		})
		return
	}

	if err := h.svc.UpdateStudent(c.Request.Context(), id, req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "BAD_REQUEST", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true, Data: gin.H{"message": "data siswa berhasil diperbarui"}})
}

// POST /api/students/bulk/class
func (h *StudentManagementHandler) BulkClass(c *gin.Context) {
	var req dto.AssignClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "VALIDATION_ERROR", Message: err.Error()},
		})
		return
	}

	if err := h.svc.BulkAssignClass(c.Request.Context(), req.StudentIDs, req.ClassID); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "INTERNAL_ERROR", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true, Data: gin.H{"message": "kelas siswa berhasil dialokasikan"}})
}

// POST /api/students/bulk/status
func (h *StudentManagementHandler) BulkStatus(c *gin.Context) {
	var req dto.BulkStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "VALIDATION_ERROR", Message: err.Error()},
		})
		return
	}

	if err := h.svc.BulkUpdateStatus(c.Request.Context(), req.StudentIDs, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "INTERNAL_ERROR", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true, Data: gin.H{"message": "status massal siswa berhasil diubah"}})
}

// POST /api/students/:id/account
func (h *StudentManagementHandler) CreateAccount(c *gin.Context) {
	id, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "BAD_REQUEST", Message: "invalid student id"},
		})
		return
	}

	var req dto.CreateStudentAccountRequest
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

	c.JSON(http.StatusCreated, dto.SuccessResponse{Success: true, Data: gin.H{"message": "akun murid berhasil dibuat dan ditautkan"}})
}

// PATCH /api/students/:id/account/toggle
func (h *StudentManagementHandler) ToggleAccount(c *gin.Context) {
	id, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "BAD_REQUEST", Message: "invalid student id"},
		})
		return
	}

	newStatus, err := h.svc.ToggleAccountStatus(c.Request.Context(), id)
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

// DELETE /api/students/:id
func (h *StudentManagementHandler) Delete(c *gin.Context) {
	id, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "BAD_REQUEST", Message: "invalid student id"},
		})
		return
	}

	if err := h.svc.DeleteStudent(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "BAD_REQUEST", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Data:    gin.H{"message": "data siswa dan akun terkait berhasil dihapus"},
	})
}
