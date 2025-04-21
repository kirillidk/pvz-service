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
