package service

import (
	"context"
	"fmt"

	"github.com/kirillidk/pvz-service/internal/dto"
	"github.com/kirillidk/pvz-service/internal/model"
	"github.com/kirillidk/pvz-service/internal/repository"
)

type PVZService struct {
	pvzRepository *repository.PVZRepository
}

func NewPVZService(pvzRepo *repository.PVZRepository) *PVZService {
	return &PVZService{
		pvzRepository: pvzRepo,
	}
}

func (s *PVZService) CreatePVZ(ctx context.Context, pvzReq dto.PVZReq) (*model.PVZ, error) {
	createdPVZ, err := s.pvzRepository.CreatePVZ(ctx, pvzReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create PVZ: %w", err)
	}

	return createdPVZ, nil
}
