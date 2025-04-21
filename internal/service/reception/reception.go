package reception

import (
	"context"
	"fmt"

	"github.com/kirillidk/pvz-service/internal/dto"
	"github.com/kirillidk/pvz-service/internal/model"
	"github.com/kirillidk/pvz-service/internal/repository"
)

type ReceptionServiceInterface interface {
	CreateReception(ctx context.Context, receptionCreateReq dto.ReceptionCreateRequest) (*model.Reception, error)
	CloseLastReception(ctx context.Context, pvzID string) (*model.Reception, error)
}

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

func (s *ReceptionService) CloseLastReception(ctx context.Context, pvzID string) (*model.Reception, error) {
	reception, err := s.receptionRepository.GetLastOpenReception(ctx, pvzID)
	if err != nil {
		return nil, fmt.Errorf("failed to find open reception: %w", err)
	}

	closedReception, err := s.receptionRepository.CloseReception(ctx, reception.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to close reception: %w", err)
	}

	return closedReception, nil
}
