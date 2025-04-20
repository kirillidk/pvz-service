package handler

import (
	"github.com/kirillidk/pvz-service/internal/config"
	"github.com/kirillidk/pvz-service/internal/service"
)

type Handler struct {
	AuthHandler      *AuthHandler
	PVZHandler       *PVZHandler
	ReceptionHandler *ReceptionHandler
	ProductHandler   *ProductHandler
}

func NewHandler(serv *service.Service, cfg *config.Config) *Handler {
	return &Handler{
		AuthHandler:      NewAuthHandler(serv.AuthService, cfg.JWT.JWTSecret),
		PVZHandler:       NewPVZHandler(serv.PVZService),
		ReceptionHandler: NewReceptionHandler(serv.ReceptionService),
		ProductHandler:   NewProductHandler(serv.ProductService),
	}
}
