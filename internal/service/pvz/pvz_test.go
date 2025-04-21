package pvz_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/kirillidk/pvz-service/internal/dto"
	"github.com/kirillidk/pvz-service/internal/model"
	service "github.com/kirillidk/pvz-service/internal/service/pvz"
)

type MockRepositories struct {
	MockPVZRepository       *MockPVZRepository
	MockReceptionRepository *MockReceptionRepository
	MockProductRepository   *MockProductRepository
}

type MockPVZRepository struct {
	CreatePVZFunc  func(ctx context.Context, req dto.PVZCreateRequest) (*model.PVZ, error)
	GetPVZListFunc func(ctx context.Context, filter dto.PVZFilterQuery) ([]model.PVZ, error)
	GetPVZByIDFunc func(ctx context.Context, pvzID string) (*model.PVZ, error)
}

func (m *MockPVZRepository) CreatePVZ(ctx context.Context, req dto.PVZCreateRequest) (*model.PVZ, error) {
	return m.CreatePVZFunc(ctx, req)
}

func (m *MockPVZRepository) GetPVZList(ctx context.Context, filter dto.PVZFilterQuery) ([]model.PVZ, error) {
	return m.GetPVZListFunc(ctx, filter)
}

func (m *MockPVZRepository) GetPVZByID(ctx context.Context, pvzID string) (*model.PVZ, error) {
	return m.GetPVZByIDFunc(ctx, pvzID)
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

func (m *MockReceptionRepository) CloseReception(ctx context.Context, pvzID string) (*model.Reception, error) {
	return m.CloseReceptionFunc(ctx, pvzID)
}

func (m *MockReceptionRepository) GetReceptionsByPVZID(ctx context.Context, pvzID string, startDate, endDate *time.Time) ([]model.Reception, error) {
	return m.GetReceptionsByPVZIDFunc(ctx, pvzID, startDate, endDate)
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

func TestPVZService_CreatePVZ(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name          string
		mockRepos     *MockRepositories
		input         dto.PVZCreateRequest
		expected      *model.PVZ
		expectedError bool
	}{
		{
			name: "Success",
			mockRepos: &MockRepositories{
				MockPVZRepository: &MockPVZRepository{
					CreatePVZFunc: func(ctx context.Context, req dto.PVZCreateRequest) (*model.PVZ, error) {
						return &model.PVZ{
							ID:               "123e4567-e89b-12d3-a456-426614174000",
							RegistrationDate: now,
							City:             req.City,
						}, nil
					},
				},
				MockReceptionRepository: &MockReceptionRepository{},
				MockProductRepository:   &MockProductRepository{},
			},
			input: dto.PVZCreateRequest{
				RegistrationDate: now,
				City:             "Москва",
			},
			expected: &model.PVZ{
				ID:               "123e4567-e89b-12d3-a456-426614174000",
				RegistrationDate: now,
				City:             "Москва",
			},
			expectedError: false,
		},
		{
			name: "Repository Error",
			mockRepos: &MockRepositories{
				MockPVZRepository: &MockPVZRepository{
					CreatePVZFunc: func(ctx context.Context, req dto.PVZCreateRequest) (*model.PVZ, error) {
						return nil, errors.New("repository error")
					},
				},
				MockReceptionRepository: &MockReceptionRepository{},
				MockProductRepository:   &MockProductRepository{},
			},
			input: dto.PVZCreateRequest{
				RegistrationDate: now,
				City:             "Москва",
			},
			expected:      nil,
			expectedError: true,
		},
		{
			name: "Invalid City",
			mockRepos: &MockRepositories{
				MockPVZRepository: &MockPVZRepository{
					CreatePVZFunc: func(ctx context.Context, req dto.PVZCreateRequest) (*model.PVZ, error) {
						return nil, errors.New("invalid city")
					},
				},
				MockReceptionRepository: &MockReceptionRepository{},
				MockProductRepository:   &MockProductRepository{},
			},
			input: dto.PVZCreateRequest{
				RegistrationDate: now,
				City:             "Новосибирск",
			},
			expected:      nil,
			expectedError: true,
		},
	}

	for ttNum, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service.NewPVZService(
				tt.mockRepos.MockPVZRepository,
				tt.mockRepos.MockReceptionRepository,
				tt.mockRepos.MockProductRepository,
			)
			got, err := s.CreatePVZ(context.Background(), tt.input)

			if (err != nil) != tt.expectedError {
				t.Errorf("Test %v: PVZService.CreatePVZ() error = %v, expectedError %v", ttNum, err, tt.expectedError)
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Test %v: PVZService.CreatePVZ() = %v, expected %v", ttNum, got, tt.expected)
			}
		})
	}
}

func TestPVZService_GetPVZList(t *testing.T) {
	now := time.Now()
	startDate := now.Add(-24 * time.Hour)
	endDate := now

	pvz1 := model.PVZ{
		ID:               "pvz-id-1",
		RegistrationDate: now,
		City:             "Москва",
	}

	pvz2 := model.PVZ{
		ID:               "pvz-id-2",
		RegistrationDate: now.Add(-1 * time.Hour),
		City:             "Санкт-Петербург",
	}

	reception1 := model.Reception{
		ID:       "reception-id-1",
		DateTime: now.Add(-30 * time.Minute),
		PVZID:    "pvz-id-1",
		Status:   "in_progress",
	}

	reception2 := model.Reception{
		ID:       "reception-id-2",
		DateTime: now.Add(-2 * time.Hour),
		PVZID:    "pvz-id-2",
		Status:   "close",
	}

	product1 := model.Product{
		ID:          "product-id-1",
		DateTime:    now.Add(-15 * time.Minute),
		Type:        "type1",
		ReceptionID: "reception-id-1",
	}

	product2 := model.Product{
		ID:          "product-id-2",
		DateTime:    now.Add(-2 * time.Hour),
		Type:        "type2",
		ReceptionID: "reception-id-2",
	}

	tests := []struct {
		name          string
		mockRepos     *MockRepositories
		filter        dto.PVZFilterQuery
		expected      *dto.PaginatedResponse
		expectedError bool
	}{
		{
			name: "Success - With Filters",
			mockRepos: &MockRepositories{
				MockPVZRepository: &MockPVZRepository{
					GetPVZListFunc: func(ctx context.Context, filter dto.PVZFilterQuery) ([]model.PVZ, error) {
						return []model.PVZ{pvz1, pvz2}, nil
					},
				},
				MockReceptionRepository: &MockReceptionRepository{
					GetReceptionsByPVZIDFunc: func(ctx context.Context, pvzID string, startDate, endDate *time.Time) ([]model.Reception, error) {
						if pvzID == "pvz-id-1" {
							return []model.Reception{reception1}, nil
						}
						return []model.Reception{reception2}, nil
					},
				},
				MockProductRepository: &MockProductRepository{
					GetProductsByReceptionIDFunc: func(ctx context.Context, receptionID string) ([]model.Product, error) {
						if receptionID == "reception-id-1" {
							return []model.Product{product1}, nil
						}
						return []model.Product{product2}, nil
					},
				},
			},
			filter: dto.PVZFilterQuery{
				Page:      1,
				Limit:     10,
				StartDate: &startDate,
				EndDate:   &endDate,
			},
			expected: &dto.PaginatedResponse{
				Data: []dto.PVZWithReceptionsResponse{
					{
						PVZ: pvz1,
						Receptions: []dto.ReceptionWithProductsResponse{
							{
								Reception: reception1,
								Products:  []model.Product{product1},
							},
						},
					},
					{
						PVZ: pvz2,
						Receptions: []dto.ReceptionWithProductsResponse{
							{
								Reception: reception2,
								Products:  []model.Product{product2},
							},
						},
					},
				},
			},
			expectedError: false,
		},
		{
			name: "Success - No Filters",
			mockRepos: &MockRepositories{
				MockPVZRepository: &MockPVZRepository{
					GetPVZListFunc: func(ctx context.Context, filter dto.PVZFilterQuery) ([]model.PVZ, error) {
						return []model.PVZ{pvz1}, nil
					},
				},
				MockReceptionRepository: &MockReceptionRepository{
					GetReceptionsByPVZIDFunc: func(ctx context.Context, pvzID string, startDate, endDate *time.Time) ([]model.Reception, error) {
						return []model.Reception{reception1}, nil
					},
				},
				MockProductRepository: &MockProductRepository{
					GetProductsByReceptionIDFunc: func(ctx context.Context, receptionID string) ([]model.Product, error) {
						return []model.Product{product1}, nil
					},
				},
			},
			filter: dto.PVZFilterQuery{
				Page:  1,
				Limit: 10,
			},
			expected: &dto.PaginatedResponse{
				Data: []dto.PVZWithReceptionsResponse{
					{
						PVZ: pvz1,
						Receptions: []dto.ReceptionWithProductsResponse{
							{
								Reception: reception1,
								Products:  []model.Product{product1},
							},
						},
					},
				},
			},
			expectedError: false,
		},
		{
			name: "Empty PVZ List",
			mockRepos: &MockRepositories{
				MockPVZRepository: &MockPVZRepository{
					GetPVZListFunc: func(ctx context.Context, filter dto.PVZFilterQuery) ([]model.PVZ, error) {
						return []model.PVZ{}, nil
					},
				},
				MockReceptionRepository: &MockReceptionRepository{},
				MockProductRepository:   &MockProductRepository{},
			},
			filter: dto.PVZFilterQuery{
				Page:  1,
				Limit: 10,
			},
			expected: &dto.PaginatedResponse{
				Data: []dto.PVZWithReceptionsResponse{},
			},
			expectedError: false,
		},
		{
			name: "PVZ Repository Error",
			mockRepos: &MockRepositories{
				MockPVZRepository: &MockPVZRepository{
					GetPVZListFunc: func(ctx context.Context, filter dto.PVZFilterQuery) ([]model.PVZ, error) {
						return nil, errors.New("pvz repository error")
					},
				},
				MockReceptionRepository: &MockReceptionRepository{},
				MockProductRepository:   &MockProductRepository{},
			},
			filter: dto.PVZFilterQuery{
				Page:  1,
				Limit: 10,
			},
			expected:      nil,
			expectedError: true,
		},
		{
			name: "Reception Repository Error",
			mockRepos: &MockRepositories{
				MockPVZRepository: &MockPVZRepository{
					GetPVZListFunc: func(ctx context.Context, filter dto.PVZFilterQuery) ([]model.PVZ, error) {
						return []model.PVZ{pvz1}, nil
					},
				},
				MockReceptionRepository: &MockReceptionRepository{
					GetReceptionsByPVZIDFunc: func(ctx context.Context, pvzID string, startDate, endDate *time.Time) ([]model.Reception, error) {
						return nil, errors.New("reception repository error")
					},
				},
				MockProductRepository: &MockProductRepository{},
			},
			filter: dto.PVZFilterQuery{
				Page:  1,
				Limit: 10,
			},
			expected:      nil,
			expectedError: true,
		},
		{
			name: "Product Repository Error",
			mockRepos: &MockRepositories{
				MockPVZRepository: &MockPVZRepository{
					GetPVZListFunc: func(ctx context.Context, filter dto.PVZFilterQuery) ([]model.PVZ, error) {
						return []model.PVZ{pvz1}, nil
					},
				},
				MockReceptionRepository: &MockReceptionRepository{
					GetReceptionsByPVZIDFunc: func(ctx context.Context, pvzID string, startDate, endDate *time.Time) ([]model.Reception, error) {
						return []model.Reception{reception1}, nil
					},
				},
				MockProductRepository: &MockProductRepository{
					GetProductsByReceptionIDFunc: func(ctx context.Context, receptionID string) ([]model.Product, error) {
						return nil, errors.New("product repository error")
					},
				},
			},
			filter: dto.PVZFilterQuery{
				Page:  1,
				Limit: 10,
			},
			expected:      nil,
			expectedError: true,
		},
		{
			name: "PVZ With No Receptions",
			mockRepos: &MockRepositories{
				MockPVZRepository: &MockPVZRepository{
					GetPVZListFunc: func(ctx context.Context, filter dto.PVZFilterQuery) ([]model.PVZ, error) {
						return []model.PVZ{pvz1}, nil
					},
				},
				MockReceptionRepository: &MockReceptionRepository{
					GetReceptionsByPVZIDFunc: func(ctx context.Context, pvzID string, startDate, endDate *time.Time) ([]model.Reception, error) {
						return []model.Reception{}, nil
					},
				},
				MockProductRepository: &MockProductRepository{},
			},
			filter: dto.PVZFilterQuery{
				Page:  1,
				Limit: 10,
			},
			expected: &dto.PaginatedResponse{
				Data: []dto.PVZWithReceptionsResponse{
					{
						PVZ:        pvz1,
						Receptions: []dto.ReceptionWithProductsResponse{},
					},
				},
			},
			expectedError: false,
		},
		{
			name: "Reception With No Products",
			mockRepos: &MockRepositories{
				MockPVZRepository: &MockPVZRepository{
					GetPVZListFunc: func(ctx context.Context, filter dto.PVZFilterQuery) ([]model.PVZ, error) {
						return []model.PVZ{pvz1}, nil
					},
				},
				MockReceptionRepository: &MockReceptionRepository{
					GetReceptionsByPVZIDFunc: func(ctx context.Context, pvzID string, startDate, endDate *time.Time) ([]model.Reception, error) {
						return []model.Reception{reception1}, nil
					},
				},
				MockProductRepository: &MockProductRepository{
					GetProductsByReceptionIDFunc: func(ctx context.Context, receptionID string) ([]model.Product, error) {
						return []model.Product{}, nil
					},
				},
			},
			filter: dto.PVZFilterQuery{
				Page:  1,
				Limit: 10,
			},
			expected: &dto.PaginatedResponse{
				Data: []dto.PVZWithReceptionsResponse{
					{
						PVZ: pvz1,
						Receptions: []dto.ReceptionWithProductsResponse{
							{
								Reception: reception1,
								Products:  []model.Product{},
							},
						},
					},
				},
			},
			expectedError: false,
		},
		{
			name: "Pagination Test",
			mockRepos: &MockRepositories{
				MockPVZRepository: &MockPVZRepository{
					GetPVZListFunc: func(ctx context.Context, filter dto.PVZFilterQuery) ([]model.PVZ, error) {
						// Verify pagination parameters are passed correctly
						if filter.Page != 2 || filter.Limit != 5 {
							t.Errorf("Expected page 2, limit 5, got page %d, limit %d", filter.Page, filter.Limit)
						}
						return []model.PVZ{pvz2}, nil
					},
				},
				MockReceptionRepository: &MockReceptionRepository{
					GetReceptionsByPVZIDFunc: func(ctx context.Context, pvzID string, startDate, endDate *time.Time) ([]model.Reception, error) {
						return []model.Reception{reception2}, nil
					},
				},
				MockProductRepository: &MockProductRepository{
					GetProductsByReceptionIDFunc: func(ctx context.Context, receptionID string) ([]model.Product, error) {
						return []model.Product{product2}, nil
					},
				},
			},
			filter: dto.PVZFilterQuery{
				Page:  2,
				Limit: 5,
			},
			expected: &dto.PaginatedResponse{
				Data: []dto.PVZWithReceptionsResponse{
					{
						PVZ: pvz2,
						Receptions: []dto.ReceptionWithProductsResponse{
							{
								Reception: reception2,
								Products:  []model.Product{product2},
							},
						},
					},
				},
			},
			expectedError: false,
		},
	}

	for ttNum, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service.NewPVZService(
				tt.mockRepos.MockPVZRepository,
				tt.mockRepos.MockReceptionRepository,
				tt.mockRepos.MockProductRepository,
			)
			got, err := s.GetPVZList(context.Background(), tt.filter)

			if (err != nil) != tt.expectedError {
				t.Errorf("Test %v: PVZService.GetPVZList() error = %v, expectedError %v", ttNum, err, tt.expectedError)
				return
			}

			if !tt.expectedError {
				if !reflect.DeepEqual(got, tt.expected) {
					t.Errorf("Test %v: PVZService.GetPVZList() = %v, expected %v", ttNum, got, tt.expected)
				}
			}
		})
	}
}
