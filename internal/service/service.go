package service

import (
	"github.com/kirillidk/pvz-service/internal/repository"
	"github.com/kirillidk/pvz-service/internal/service/auth"
)

type Service struct {
	AuthService *auth.AuthService
}

func NewService(repository *repository.Repository, jwtSecret string) *Service {
	return &Service{
		AuthService: auth.NewAuthService(repository.UserRepository, jwtSecret),
	}
}
