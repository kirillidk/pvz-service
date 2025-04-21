package service

import (
	"github.com/kirillidk/pvz-service/internal/repository"
	"github.com/kirillidk/pvz-service/internal/service/auth"
	"github.com/kirillidk/pvz-service/internal/service/product"
	"github.com/kirillidk/pvz-service/internal/service/pvz"
	"github.com/kirillidk/pvz-service/internal/service/reception"
)

type Service struct {
	AuthService      *auth.AuthService
	PVZService       *pvz.PVZService
	ReceptionService *reception.ReceptionService
	ProductService   *product.ProductService
}

func NewService(repository *repository.Repository, jwtSecret string) *Service {
	return &Service{
		AuthService:      auth.NewAuthService(repository.UserRepository, jwtSecret),
		PVZService:       pvz.NewPVZService(repository.PVZRepository, repository.ReceptionRepository, repository.ProductRepository),
		ReceptionService: reception.NewReceptionService(repository.ReceptionRepository),
		ProductService:   product.NewProductService(repository.ProductRepository, repository.ReceptionRepository),
	}
}
