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

type MockPVZRepository struct {
	CreatePVZFunc func(ctx context.Context, req dto.PVZRequest) (*model.PVZ, error)
}

func (m *MockPVZRepository) CreatePVZ(ctx context.Context, req dto.PVZRequest) (*model.PVZ, error) {
	return m.CreatePVZFunc(ctx, req)
}

func TestPVZService_CreatePVZ(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name          string
		mockRepo      *MockPVZRepository
		input         dto.PVZRequest
		expected      *model.PVZ
		expectedError bool
	}{
		{
			name: "Success",
			mockRepo: &MockPVZRepository{
				CreatePVZFunc: func(ctx context.Context, req dto.PVZRequest) (*model.PVZ, error) {
					return &model.PVZ{
						ID:               "123e4567-e89b-12d3-a456-426614174000",
						RegistrationDate: now,
						City:             req.City,
					}, nil
				},
			},
			input: dto.PVZRequest{
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
			mockRepo: &MockPVZRepository{
				CreatePVZFunc: func(ctx context.Context, req dto.PVZRequest) (*model.PVZ, error) {
					return nil, errors.New("repository error")
				},
			},
			input: dto.PVZRequest{
				RegistrationDate: now,
				City:             "Москва",
			},
			expected:      nil,
			expectedError: true,
		},
		{
			name: "Invalid City",
			mockRepo: &MockPVZRepository{
				CreatePVZFunc: func(ctx context.Context, req dto.PVZRequest) (*model.PVZ, error) {
					return nil, errors.New("invalid city")
				},
			},
			input: dto.PVZRequest{
				RegistrationDate: now,
				City:             "Новосибирск",
			},
			expected:      nil,
			expectedError: true,
		},
	}

	for ttNum, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service.NewPVZService(tt.mockRepo)
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
