package handler

import (
	"github.com/kirillidk/pvz-service/internal/config"
	"github.com/kirillidk/pvz-service/internal/service"
)

type Handler struct {
	AuthHandler *AuthHandler
}

func NewHandler(serv *service.Service, cfg *config.Config) *Handler {
	return &Handler{
		AuthHandler: NewAuthHandler(serv.AuthService, cfg.JWT.JWTSecret),
	}
}
