package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/geniussociety/backend/internal/config"
	"github.com/geniussociety/backend/internal/dto"
	"github.com/geniussociety/backend/internal/middleware"
	"github.com/geniussociety/backend/internal/service"
)

type AuthHandler struct {
	authSvc *service.AuthService
	cfg     *config.Config
}

func NewAuthHandler(authSvc *service.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{authSvc: authSvc, cfg: cfg}
}

// POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "VALIDATION_ERROR", Message: err.Error()},
		})
		return
	}

	user, token, err := h.authSvc.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "UNAUTHORIZED", Message: err.Error()},
		})
		return
	}

	// Set HTTP-only cookie
	secure := h.cfg.AppEnv == "production"
	sameSite := http.SameSiteLaxMode
	if h.cfg.CookieSameSite == "None" {
		sameSite = http.SameSiteNoneMode
	} else if h.cfg.CookieSameSite == "Strict" {
		sameSite = http.SameSiteStrictMode
	}

	c.SetCookie("token", token, int(24*time.Hour/time.Second), "/", h.cfg.CookieDomain, secure, true)
	c.SetSameSite(sameSite)

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

// POST /api/auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", h.cfg.CookieDomain, h.cfg.AppEnv == "production", true)
	c.JSON(http.StatusOK, dto.SuccessResponse{Success: true, Data: gin.H{"message": "logged out"}})
}

// GET /api/auth/me
func (h *AuthHandler) Me(c *gin.Context) {
	userID := middleware.GetUserID(c)
	user, err := h.authSvc.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Error:   dto.ErrorBody{Code: "UNAUTHORIZED", Message: "user not found"},
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
