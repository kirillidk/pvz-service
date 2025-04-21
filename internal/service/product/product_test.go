package product_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/kirillidk/pvz-service/internal/dto"
	"github.com/kirillidk/pvz-service/internal/model"
	"github.com/kirillidk/pvz-service/internal/service/product"
)

type MockRepositories struct {
	MockProductRepository   *MockProductRepository
	MockReceptionRepository *MockReceptionRepository
}

type MockProductRepository struct {
	CreateProductFunc             func(ctx context.Context, productType string, receptionID string) (*model.Product, error)
	GetLastProductInReceptionFunc func(ctx context.Context, receptionID string) (*model.Product, error)
	DeleteProductFunc             func(ctx context.Context, productID string) error
	GetProductsByReceptionIDFunc  func(ctx context.Context, receptionID string) ([]model.Product, error)
}

func (m *MockProductRepository) CreateProduct(ctx context.Context, productType string, receptionID string) (*model.Product, error) {
	return m.CreateProductFunc(ctx, productType, receptionID)
}

func (m *MockProductRepository) GetLastProductInReception(ctx context.Context, receptionID string) (*model.Product, error) {
	return m.GetLastProductInReceptionFunc(ctx, receptionID)
}

func (m *MockProductRepository) DeleteProduct(ctx context.Context, productID string) error {
	return m.DeleteProductFunc(ctx, productID)
}

func (m *MockProductRepository) GetProductsByReceptionID(ctx context.Context, receptionID string) ([]model.Product, error) {
	return m.GetProductsByReceptionIDFunc(ctx, receptionID)
}

type MockReceptionRepository struct {
	CreateReceptionFunc      func(ctx context.Context, req dto.ReceptionCreateRequest) (*model.Reception, error)
	HasOpenReceptionFunc     func(ctx context.Context, pvzID string) (bool, error)
	GetLastOpenReceptionFunc func(ctx context.Context, pvzID string) (*model.Reception, error)
	CloseReceptionFunc       func(ctx context.Context, receptionID string) (*model.Reception, error)
	GetReceptionsByPVZIDFunc func(ctx context.Context, pvzID string, startDate, endDate *time.Time) ([]model.Reception, error)
}

func (m *MockReceptionRepository) CreateReception(ctx context.Context, req dto.ReceptionCreateRequest) (*model.Reception, error) {
	return m.CreateReceptionFunc(ctx, req)
}

func (m *MockReceptionRepository) HasOpenReception(ctx context.Context, pvzID string) (bool, error) {
	return m.HasOpenReceptionFunc(ctx, pvzID)
}

func (m *MockReceptionRepository) GetLastOpenReception(ctx context.Context, pvzID string) (*model.Reception, error) {
	return m.GetLastOpenReceptionFunc(ctx, pvzID)
}

func (m *MockReceptionRepository) CloseReception(ctx context.Context, receptionID string) (*model.Reception, error) {
	return m.CloseReceptionFunc(ctx, receptionID)
}

func (m *MockReceptionRepository) GetReceptionsByPVZID(ctx context.Context, pvzID string, startDate, endDate *time.Time) ([]model.Reception, error) {
	return m.GetReceptionsByPVZIDFunc(ctx, pvzID, startDate, endDate)
}

func TestProductService_CreateProduct(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name          string
		mocks         MockRepositories
		input         dto.ProductCreateRequest
		expected      *model.Product
		expectedError bool
	}{
		{
			name: "Success",
			mocks: MockRepositories{
				MockProductRepository: &MockProductRepository{
					CreateProductFunc: func(ctx context.Context, productType string, receptionID string) (*model.Product, error) {
						return &model.Product{
							ID:          "123e4567-e89b-12d3-a456-426614174001",
							DateTime:    now,
							Type:        productType,
							ReceptionID: receptionID,
						}, nil
					},
				},
				MockReceptionRepository: &MockReceptionRepository{
					GetLastOpenReceptionFunc: func(ctx context.Context, pvzID string) (*model.Reception, error) {
						return &model.Reception{
							ID:       "123e4567-e89b-12d3-a456-426614174002",
							DateTime: now,
							PVZID:    pvzID,
							Status:   "in_progress",
						}, nil
					},
				},
			},
			input: dto.ProductCreateRequest{
				Type:  "electronics",
				PVZID: "123e4567-e89b-12d3-a456-426614174003",
			},
			expected: &model.Product{
				ID:          "123e4567-e89b-12d3-a456-426614174001",
				DateTime:    now,
				Type:        "electronics",
				ReceptionID: "123e4567-e89b-12d3-a456-426614174002",
			},
			expectedError: false,
		},
		{
			name: "No Open Reception",
			mocks: MockRepositories{
				MockProductRepository: &MockProductRepository{},
				MockReceptionRepository: &MockReceptionRepository{
					GetLastOpenReceptionFunc: func(ctx context.Context, pvzID string) (*model.Reception, error) {
						return nil, errors.New("no open reception found for this PVZ")
					},
				},
			},
			input: dto.ProductCreateRequest{
				Type:  "electronics",
				PVZID: "123e4567-e89b-12d3-a456-426614174003",
			},
			expected:      nil,
			expectedError: true,
		},
		{
			name: "Product Creation Error",
			mocks: MockRepositories{
				MockProductRepository: &MockProductRepository{
					CreateProductFunc: func(ctx context.Context, productType string, receptionID string) (*model.Product, error) {
						return nil, errors.New("failed to create product")
					},
				},
				MockReceptionRepository: &MockReceptionRepository{
					GetLastOpenReceptionFunc: func(ctx context.Context, pvzID string) (*model.Reception, error) {
						return &model.Reception{
							ID:       "123e4567-e89b-12d3-a456-426614174002",
							DateTime: now,
							PVZID:    pvzID,
							Status:   "in_progress",
						}, nil
					},
				},
			},
			input: dto.ProductCreateRequest{
				Type:  "electronics",
				PVZID: "123e4567-e89b-12d3-a456-426614174003",
			},
			expected:      nil,
			expectedError: true,
		},
	}

	for ttNum, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := product.NewProductService(tt.mocks.MockProductRepository, tt.mocks.MockReceptionRepository)
			got, err := s.CreateProduct(context.Background(), tt.input)

			if (err != nil) != tt.expectedError {
				t.Errorf("Test %v: ProductService.CreateProduct() error = %v, expectedError %v", ttNum, err, tt.expectedError)
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Test %v: ProductService.CreateProduct() = %v, expected %v", ttNum, got, tt.expected)
			}
		})
	}
}

func TestProductService_DeleteLastProduct(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name          string
		mocks         MockRepositories
		pvzID         string
		expectedError bool
	}{
		{
			name: "Success",
			mocks: MockRepositories{
				MockProductRepository: &MockProductRepository{
					GetLastProductInReceptionFunc: func(ctx context.Context, receptionID string) (*model.Product, error) {
						return &model.Product{
							ID:          "123e4567-e89b-12d3-a456-426614174001",
							DateTime:    now,
							Type:        "electronics",
							ReceptionID: receptionID,
						}, nil
					},
					DeleteProductFunc: func(ctx context.Context, productID string) error {
						return nil
					},
				},
				MockReceptionRepository: &MockReceptionRepository{
					GetLastOpenReceptionFunc: func(ctx context.Context, pvzID string) (*model.Reception, error) {
						return &model.Reception{
							ID:       "123e4567-e89b-12d3-a456-426614174002",
							DateTime: now,
							PVZID:    pvzID,
							Status:   "in_progress",
						}, nil
					},
				},
			},
			pvzID:         "123e4567-e89b-12d3-a456-426614174003",
			expectedError: false,
		},
		{
			name: "No Open Reception",
			mocks: MockRepositories{
				MockReceptionRepository: &MockReceptionRepository{
					GetLastOpenReceptionFunc: func(ctx context.Context, pvzID string) (*model.Reception, error) {
						return nil, errors.New("no open reception found for this PVZ")
					},
				},
			},
			pvzID:         "123e4567-e89b-12d3-a456-426614174003",
			expectedError: true,
		},
		{
			name: "No Product Found",
			mocks: MockRepositories{
				MockProductRepository: &MockProductRepository{
					GetLastProductInReceptionFunc: func(ctx context.Context, receptionID string) (*model.Product, error) {
						return nil, errors.New("no products found for this reception")
					},
					DeleteProductFunc: func(ctx context.Context, productID string) error {
						return nil
					},
				},
				MockReceptionRepository: &MockReceptionRepository{
					GetLastOpenReceptionFunc: func(ctx context.Context, pvzID string) (*model.Reception, error) {
						return &model.Reception{
							ID:       "123e4567-e89b-12d3-a456-426614174002",
							DateTime: now,
							PVZID:    pvzID,
							Status:   "in_progress",
						}, nil
					},
				},
			},
			pvzID:         "123e4567-e89b-12d3-a456-426614174003",
			expectedError: true,
		},
		{
			name: "Delete Product Error",
			mocks: MockRepositories{
				MockProductRepository: &MockProductRepository{
					GetLastProductInReceptionFunc: func(ctx context.Context, receptionID string) (*model.Product, error) {
						return &model.Product{
							ID:          "123e4567-e89b-12d3-a456-426614174001",
							DateTime:    now,
							Type:        "electronics",
							ReceptionID: receptionID,
						}, nil
					},
					DeleteProductFunc: func(ctx context.Context, productID string) error {
						return errors.New("failed to delete product")
					},
				},
				MockReceptionRepository: &MockReceptionRepository{
					GetLastOpenReceptionFunc: func(ctx context.Context, pvzID string) (*model.Reception, error) {
						return &model.Reception{
							ID:       "123e4567-e89b-12d3-a456-426614174002",
							DateTime: now,
							PVZID:    pvzID,
							Status:   "in_progress",
						}, nil
					},
				},
			},
			pvzID:         "123e4567-e89b-12d3-a456-426614174003",
			expectedError: true,
		},
	}

	for ttNum, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := product.NewProductService(tt.mocks.MockProductRepository, tt.mocks.MockReceptionRepository)
			err := s.DeleteLastProduct(context.Background(), tt.pvzID)

			if (err != nil) != tt.expectedError {
				t.Errorf("Test %v: ProductService.DeleteLastProduct() error = %v, expectedError %v", ttNum, err, tt.expectedError)
			}
		})
	}
}
