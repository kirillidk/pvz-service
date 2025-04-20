package pvz

import (
	"context"
	"fmt"

	"github.com/kirillidk/pvz-service/internal/dto"
	"github.com/kirillidk/pvz-service/internal/model"
	"github.com/kirillidk/pvz-service/internal/repository"
)

type PVZService struct {
	pvzRepository repository.PVZRepositoryInterface
}

func NewPVZService(pvzRepo repository.PVZRepositoryInterface) *PVZService {
	return &PVZService{
		pvzRepository: pvzRepo,
	}
}

func (s *PVZService) CreatePVZ(ctx context.Context, pvzReq dto.PVZCreateRequest) (*model.PVZ, error) {
	createdPVZ, err := s.pvzRepository.CreatePVZ(ctx, pvzReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create PVZ: %w", err)
	}

	return createdPVZ, nil
}
