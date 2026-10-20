package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/geniussociety/backend/internal/dto"
	"github.com/geniussociety/backend/internal/model"
	"github.com/geniussociety/backend/internal/service"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserHandler struct {
	userSvc *service.UserService
}

func NewUserHandler(userSvc *service.UserService) *UserHandler {
	return &UserHandler{userSvc: userSvc}
}

// GET /api/users
func (h *UserHandler) List(c *gin.Context) {
	role := c.Query("role")
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64)
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	users, total, err := h.userSvc.ListUsers(c.Request.Context(), role, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "INTERNAL_ERROR", Message: err.Error()},
		})
		return
	}

	resp := make([]dto.UserResponse, 0, len(users))
	for _, u := range users {
		resp = append(resp, dto.UserResponse{
			ID:       u.ID.Hex(),
			Name:     u.Name,
			Email:    u.Email,
			Username: u.Username,
			Role:     u.Role,
			IsActive: u.IsActive,
		})
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Data:    resp,
		Meta:    dto.PaginationMeta{Page: page, Limit: limit, Total: total},
	})
}

// POST /api/users
func (h *UserHandler) Create(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "VALIDATION_ERROR", Message: err.Error()},
		})
		return
	}

	user, err := h.userSvc.CreateUser(c.Request.Context(), req.Name, req.Email, req.Username, req.Password, req.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "BAD_REQUEST", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse{
		Success: true,
		Data: dto.UserResponse{
			ID:       user.ID.Hex(),
			Name:     user.Name,
			Email:    user.Email,
			Username: user.Username,
			Role:     user.Role,
			IsActive: user.IsActive,
		},
	})
}

// GET /api/users/:id
func (h *UserHandler) GetOne(c *gin.Context) {
	id, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "BAD_REQUEST", Message: "invalid id"},
		})
		return
	}

	user, err := h.userSvc.GetUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "NOT_FOUND", Message: "user not found"},
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Data: dto.UserResponse{
			ID:       user.ID.Hex(),
			Name:     user.Name,
			Email:    user.Email,
			Username: user.Username,
			Role:     user.Role,
			IsActive: user.IsActive,
		},
	})
}

// PUT /api/users/:id
func (h *UserHandler) Update(c *gin.Context) {
	id, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "BAD_REQUEST", Message: "invalid id"},
		})
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "VALIDATION_ERROR", Message: err.Error()},
		})
		return
	}

	if err := h.userSvc.UpdateUser(c.Request.Context(), id, req.Name, req.Email, req.Username, req.Password, req.IsActive); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "INTERNAL_ERROR", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true, Data: gin.H{"message": "updated"}})
}

// DELETE /api/users/:id
func (h *UserHandler) Delete(c *gin.Context) {
	id, err := bson.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "BAD_REQUEST", Message: "invalid id"},
		})
		return
	}

	if err := h.userSvc.DeleteUser(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "INTERNAL_ERROR", Message: err.Error()},
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true, Data: gin.H{"message": "deleted"}})
}

// Ensure model is imported (used for role constants elsewhere)
var _ model.Role = ""
