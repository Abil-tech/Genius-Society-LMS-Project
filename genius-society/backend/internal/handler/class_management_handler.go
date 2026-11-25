package handler

import (
	"net/http"
	"strconv"

	"github.com/geniussociety/backend/internal/dto"
	"github.com/geniussociety/backend/internal/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// ClassManagementHandler handles class‑management API endpoints.
type ClassManagementHandler struct {
	svc *service.ClassService
}

func NewClassManagementHandler(svc *service.ClassService) *ClassManagementHandler {
	return &ClassManagementHandler{svc: svc}
}

// GET /api/classes
func (h *ClassManagementHandler) List(c *gin.Context) {
	// Query parameters for filters and pagination
	search := c.Query("search")
	educationLevel := c.Query("educationLevel")
	department := c.Query("department")
	academicYearId := c.Query("academicYearId")
	status := c.Query("status")
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64)

	// Build filter DTO
	filter := &dto.ClassListFilter{
		Code:   nil,
		Name:   nil,
		Status: nil,
		Page:   page,
		Limit:  limit,
	}
	if search != "" {
		filter.Name = &search
	}
	if educationLevel != "" {
		filter.EducationLevel = &educationLevel
	}
	if department != "" {
		filter.Department = &department
	}
	if academicYearId != "" {
		filter.AcademicYearID = &academicYearId
	}
	if status != "" {
		filter.Status = &status
	}

	result, err := h.svc.ListClasses(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "INTERNAL_ERROR", Message: err.Error()}})
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true, Data: result.Items, Meta: dto.PaginationMeta{Page: result.Page, Limit: result.Limit, Total: result.Total}})
}

// GET /api/classes/:id
func (h *ClassManagementHandler) GetOne(c *gin.Context) {
	classIDHex := c.Param("id")
	if classIDHex == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: "missing ID"}})
		return
	}
	classID, err := bson.ObjectIDFromHex(classIDHex)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: "invalid class ID"}})
		return
	}
	class, err := h.svc.GetClassDetail(c.Request.Context(), classID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "NOT_FOUND", Message: err.Error()}})
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true, Data: class})
}

// POST /api/classes
func (h *ClassManagementHandler) Create(c *gin.Context) {
	var req dto.CreateClassReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "VALIDATION_ERROR", Message: err.Error()}})
		return
	}
	academicYearID, err := bson.ObjectIDFromHex(req.AcademicYearID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: "invalid academicYearId"}})
		return
	}
	created, err := h.svc.CreateClass(c.Request.Context(), academicYearID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: err.Error()}})
		return
	}
	c.JSON(http.StatusCreated, dto.SuccessResponse{Success: true, Data: created})
}

// PUT /api/classes/:id
func (h *ClassManagementHandler) Update(c *gin.Context) {
	var req dto.UpdateClassReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "VALIDATION_ERROR", Message: err.Error()}})
		return
	}
	classIDHex := c.Param("id")
	if classIDHex == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: "missing ID"}})
		return
	}
	classID, err := bson.ObjectIDFromHex(classIDHex)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: "invalid class ID"}})
		return
	}
	updated, err := h.svc.UpdateClass(c.Request.Context(), classID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: err.Error()}})
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true, Data: updated})
}

// DELETE /api/classes/:id (soft‑archive)
func (h *ClassManagementHandler) Archive(c *gin.Context) {
	classIDHex := c.Param("id")
	if classIDHex == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: "missing ID"}})
		return
	}
	classID, err := bson.ObjectIDFromHex(classIDHex)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: "invalid class ID"}})
		return
	}
	if err := h.svc.ArchiveClass(c.Request.Context(), classID); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: err.Error()}})
		return
	}
	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true, Data: gin.H{"message": "kelas berhasil diarsipkan"}})
}
