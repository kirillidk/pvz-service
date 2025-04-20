package service

import (
	"context"
	"fmt"

	"github.com/kirillidk/pvz-service/internal/dto"
	"github.com/kirillidk/pvz-service/internal/model"
	"github.com/kirillidk/pvz-service/internal/repository"
)

type ReceptionService struct {
	receptionRepository repository.ReceptionRepositoryInterface
}

func NewReceptionService(receptionRepo repository.ReceptionRepositoryInterface) *ReceptionService {
	return &ReceptionService{
		receptionRepository: receptionRepo,
	}
}

func (s *ReceptionService) CreateReception(ctx context.Context, receptionCreateReq dto.ReceptionCreateRequest) (*model.Reception, error) {
	reception, err := s.receptionRepository.CreateReception(ctx, receptionCreateReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create reception: %w", err)
	}

	return reception, nil
}
