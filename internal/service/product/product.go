package product

import (
	"context"
	"fmt"

	"github.com/kirillidk/pvz-service/internal/dto"
	"github.com/kirillidk/pvz-service/internal/model"
	"github.com/kirillidk/pvz-service/internal/repository"
)

type ProductService struct {
	productRepository   repository.ProductRepositoryInterface
	receptionRepository repository.ReceptionRepositoryInterface
}

func NewProductService(
	productRepo repository.ProductRepositoryInterface,
	receptionRepo repository.ReceptionRepositoryInterface,
) *ProductService {
	return &ProductService{
		productRepository:   productRepo,
		receptionRepository: receptionRepo,
	}
}

func (s *ProductService) CreateProduct(ctx context.Context, req dto.ProductCreateRequest) (*model.Product, error) {
	reception, err := s.receptionRepository.GetLastOpenReception(ctx, req.PVZID)
	if err != nil {
		return nil, fmt.Errorf("failed to find open reception: %w", err)
	}

	product, err := s.productRepository.CreateProduct(ctx, req.Type, reception.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return product, nil
}

func (s *ProductService) DeleteLastProduct(ctx context.Context, pvzID string) error {
	reception, err := s.receptionRepository.GetLastOpenReception(ctx, pvzID)
	if err != nil {
		return fmt.Errorf("failed to find open reception: %w", err)
	}

	lastProduct, err := s.productRepository.GetLastProductInReception(ctx, reception.ID)
	if err != nil {
		return fmt.Errorf("failed to get last product: %w", err)
	}

	err = s.productRepository.DeleteProduct(ctx, lastProduct.ID)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}
