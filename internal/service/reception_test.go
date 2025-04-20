package service_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/kirillidk/pvz-service/internal/dto"
	"github.com/kirillidk/pvz-service/internal/model"
	"github.com/kirillidk/pvz-service/internal/service"
)

type MockReceptionRepository struct {
	CreateReceptionFunc  func(ctx context.Context, req dto.ReceptionCreateRequest) (*model.Reception, error)
	HasOpenReceptionFunc func(ctx context.Context, req dto.ReceptionCreateRequest) (bool, error)
}

func (m *MockReceptionRepository) CreateReception(ctx context.Context, req dto.ReceptionCreateRequest) (*model.Reception, error) {
	return m.CreateReceptionFunc(ctx, req)
}

func (m *MockReceptionRepository) HasOpenReception(ctx context.Context, req dto.ReceptionCreateRequest) (bool, error) {
	return m.HasOpenReceptionFunc(ctx, req)
}

func TestReceptionService_CreateReception(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name          string
		mockRepo      *MockReceptionRepository
		input         dto.ReceptionCreateRequest
		expected      *model.Reception
		expectedError bool
	}{
		{
			name: "Success",
			mockRepo: &MockReceptionRepository{
				CreateReceptionFunc: func(ctx context.Context, req dto.ReceptionCreateRequest) (*model.Reception, error) {
					return &model.Reception{
						ID:       "123e4567-e89b-12d3-a456-426614174000",
						DateTime: now,
						PVZID:    req.PVZID,
						Status:   "in_progress",
					}, nil
				},
			},
			input: dto.ReceptionCreateRequest{
				PVZID: "123e4567-e89b-12d3-a456-426614174000",
			},
			expected: &model.Reception{
				ID:       "123e4567-e89b-12d3-a456-426614174000",
				DateTime: now,
				PVZID:    "123e4567-e89b-12d3-a456-426614174000",
				Status:   "in_progress",
			},
			expectedError: false,
		},
		{
			name: "Repository Error",
			mockRepo: &MockReceptionRepository{
				CreateReceptionFunc: func(ctx context.Context, req dto.ReceptionCreateRequest) (*model.Reception, error) {
					return nil, errors.New("repository error")
				},
			},
			input: dto.ReceptionCreateRequest{
				PVZID: "123e4567-e89b-12d3-a456-426614174000",
			},
			expected:      nil,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service.NewReceptionService(tt.mockRepo)
			got, err := s.CreateReception(context.Background(), tt.input)

			if (err != nil) != tt.expectedError {
				t.Errorf("ReceptionService.CreateReception() error = %v, expectedError %v", err, tt.expectedError)
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ReceptionService.CreateReception() = %v, expected %v", got, tt.expected)
			}
		})
	}
}
