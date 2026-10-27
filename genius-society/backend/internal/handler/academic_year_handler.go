// File: d:/projectME/Genius-Society-Baru/genius-society/backend/internal/handler/academic_year_handler.go
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

// NewAcademicYearHandler creates the handler (injection of the service)
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

	year, err := h.svc.Create(c.Request.Context(), schoolID, req)
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

	ayID, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: "invalid id"}})
		return
	}
	year, err := h.svc.Update(c.Request.Context(), schoolID, ayID, req)
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

	ayID, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: "invalid id"}})
		return
	}
	if err := h.svc.Activate(c.Request.Context(), schoolID, ayID); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "INTERNAL_ERROR", Message: err.Error()}})
		return
	}
	c.Status(http.StatusNoContent)
}

// DELETE /api/academic-years/:id
func (h *AcademicYearHandler) Delete(c *gin.Context) {
	schoolIDHex := middleware.GetSchoolID(c)
	schoolID, _ := bson.ObjectIDFromHex(schoolIDHex)

	ayID, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: "invalid id"}})
		return
	}
	if err := h.svc.Delete(c.Request.Context(), schoolID, ayID); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Success: false, Error: dto.ErrorBody{Code: "INTERNAL_ERROR", Message: err.Error()}})
		return
	}
	c.Status(http.StatusNoContent)
}
