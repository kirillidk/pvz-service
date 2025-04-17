package dto

import "github.com/kirillidk/pvz-service/internal/model"

type DummyLoginRequest struct {
	Role model.UserRole `json:"role" binding:"required,oneof=employee moderator"`
}
