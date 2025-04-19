package service

import (
	"github.com/kirillidk/pvz-service/internal/repository"
	"github.com/kirillidk/pvz-service/internal/service/auth"
)

type Service struct {
	AuthService *auth.AuthService
	PVZService  *PVZService
}

func NewService(repository *repository.Repository, jwtSecret string) *Service {
	return &Service{
		AuthService: auth.NewAuthService(repository.UserRepository, jwtSecret),
		PVZService:  NewPVZService(repository.PVZRepository),
	}
}
