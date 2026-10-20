package dto

import "github.com/geniussociety/backend/internal/model"

// Auth
type LoginRequest struct {
	Login    string `json:"login" binding:"required"`    // email or username
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token,omitempty"` // only if not cookie-mode
}

// User
type UserResponse struct {
	ID       string     `json:"id"`
	Name     string     `json:"name"`
	Email    string     `json:"email"`
	Username string     `json:"username"`
	Role     model.Role `json:"role"`
	IsActive bool       `json:"isActive"`
}

type CreateUserRequest struct {
	Name     string     `json:"name" binding:"required"`
	Email    string     `json:"email" binding:"required,email"`
	Username string     `json:"username" binding:"required"`
	Password string     `json:"password" binding:"required,min=8"`
	Role     model.Role `json:"role" binding:"required"`
}

type UpdateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email" binding:"omitempty,email"`
	Username string `json:"username"`
	Password string `json:"password" binding:"omitempty,min=8"`
	IsActive *bool  `json:"isActive"`
}

// Pagination
type PaginationMeta struct {
	Page  int64 `json:"page"`
	Limit int64 `json:"limit"`
	Total int64 `json:"total"`
}

// API response wrappers
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

type ErrorResponse struct {
	Success bool      `json:"success"`
	Error   ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
