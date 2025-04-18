package dto

import "github.com/kirillidk/pvz-service/internal/model"

type DummyLoginRequest struct {
	Role model.UserRole `json:"role" binding:"required,oneof=employee moderator"`
}

type RegisterRequest struct {
	Email    string         `json:"email" binding:"required,email"`
	Password string         `json:"password" binding:"required,min=6"`
	Role     model.UserRole `json:"role" binding:"required,oneof=employee moderator"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}
