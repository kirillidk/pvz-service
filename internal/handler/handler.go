package handler

import "github.com/kirillidk/pvz-service/internal/config"

type Handler struct {
	AuthHandler *AuthHandler
}

func NewHandler(cfg *config.Config) *Handler {
	return &Handler{
		AuthHandler: NewAuthHandler(cfg.JWT.JWTSecret),
	}
}
