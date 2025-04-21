package pvz

import (
	"context"
	"fmt"

	"github.com/kirillidk/pvz-service/internal/dto"
	"github.com/kirillidk/pvz-service/internal/model"
	"github.com/kirillidk/pvz-service/internal/repository"
)

type PVZServiceInterface interface {
	CreatePVZ(ctx context.Context, pvzReq dto.PVZCreateRequest) (*model.PVZ, error)
	GetPVZList(ctx context.Context, filter dto.PVZFilterQuery) (*dto.PaginatedResponse, error)
}

type PVZService struct {
	pvzRepository       repository.PVZRepositoryInterface
	receptionRepository repository.ReceptionRepositoryInterface
	productRepository   repository.ProductRepositoryInterface
}

func NewPVZService(
	pvzRepo repository.PVZRepositoryInterface,
	receptionRepo repository.ReceptionRepositoryInterface,
	productRepo repository.ProductRepositoryInterface,
) *PVZService {
	return &PVZService{
		pvzRepository:       pvzRepo,
		receptionRepository: receptionRepo,
		productRepository:   productRepo,
	}
}

func (s *PVZService) CreatePVZ(ctx context.Context, pvzReq dto.PVZCreateRequest) (*model.PVZ, error) {
	createdPVZ, err := s.pvzRepository.CreatePVZ(ctx, pvzReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create PVZ: %w", err)
	}

	return createdPVZ, nil
}

func (s *PVZService) GetPVZList(ctx context.Context, filter dto.PVZFilterQuery) (*dto.PaginatedResponse, error) {
	pvzList, err := s.pvzRepository.GetPVZList(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get PVZ list: %w", err)
	}

	result := &dto.PaginatedResponse{
		Data: make([]dto.PVZWithReceptionsResponse, 0, len(pvzList)),
	}

	for _, pvz := range pvzList {
		receptions, err := s.receptionRepository.GetReceptionsByPVZID(ctx, pvz.ID, filter.StartDate, filter.EndDate)
		if err != nil {
			return nil, fmt.Errorf("failed to get receptions for PVZ %s: %w", pvz.ID, err)
		}

		pvzResponse := dto.PVZWithReceptionsResponse{
			PVZ:        pvz,
			Receptions: make([]dto.ReceptionWithProductsResponse, 0, len(receptions)),
		}

		for _, reception := range receptions {
			products, err := s.productRepository.GetProductsByReceptionID(ctx, reception.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to get products for reception %s: %w", reception.ID, err)
			}

			receptionResponse := dto.ReceptionWithProductsResponse{
				Reception: reception,
				Products:  products,
			}

			pvzResponse.Receptions = append(pvzResponse.Receptions, receptionResponse)
		}

		result.Data = append(result.Data, pvzResponse)
	}

	return result, nil
}
