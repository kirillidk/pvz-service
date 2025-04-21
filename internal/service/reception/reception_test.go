package reception_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/kirillidk/pvz-service/internal/dto"
	"github.com/kirillidk/pvz-service/internal/model"
	"github.com/kirillidk/pvz-service/internal/service/reception"
)

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

	for ttNum, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := reception.NewReceptionService(tt.mockRepo)
			got, err := s.CreateReception(context.Background(), tt.input)

			if (err != nil) != tt.expectedError {
				t.Errorf("Test %v: ReceptionService.CreateReception() error = %v, expectedError %v", ttNum, err, tt.expectedError)
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Test %v: ReceptionService.CreateReception() = %v, expected %v", ttNum, got, tt.expected)
			}
		})
	}
}

func TestReceptionService_CloseLastReception(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name          string
		mockRepo      *MockReceptionRepository
		pvzID         string
		expected      *model.Reception
		expectedError bool
	}{
		{
			name: "Success",
			mockRepo: &MockReceptionRepository{
				GetLastOpenReceptionFunc: func(ctx context.Context, pvzID string) (*model.Reception, error) {
					return &model.Reception{
						ID:       "123e4567-e89b-12d3-a456-426614174000",
						DateTime: now,
						PVZID:    pvzID,
						Status:   "in_progress",
					}, nil
				},
				CloseReceptionFunc: func(ctx context.Context, receptionID string) (*model.Reception, error) {
					return &model.Reception{
						ID:       receptionID,
						DateTime: now,
						PVZID:    "123e4567-e89b-12d3-a456-426614174000",
						Status:   "close",
					}, nil
				},
			},
			pvzID: "123e4567-e89b-12d3-a456-426614174000",
			expected: &model.Reception{
				ID:       "123e4567-e89b-12d3-a456-426614174000",
				DateTime: now,
				PVZID:    "123e4567-e89b-12d3-a456-426614174000",
				Status:   "close",
			},
			expectedError: false,
		},
		{
			name: "No Open Reception",
			mockRepo: &MockReceptionRepository{
				GetLastOpenReceptionFunc: func(ctx context.Context, pvzID string) (*model.Reception, error) {
					return nil, errors.New("no open reception found for this PVZ")
				},
			},
			pvzID:         "123e4567-e89b-12d3-a456-426614174000",
			expected:      nil,
			expectedError: true,
		},
		{
			name: "Close Reception Error",
			mockRepo: &MockReceptionRepository{
				GetLastOpenReceptionFunc: func(ctx context.Context, pvzID string) (*model.Reception, error) {
					return &model.Reception{
						ID:       "123e4567-e89b-12d3-a456-426614174000",
						DateTime: now,
						PVZID:    pvzID,
						Status:   "in_progress",
					}, nil
				},
				CloseReceptionFunc: func(ctx context.Context, receptionID string) (*model.Reception, error) {
					return nil, errors.New("failed to close reception")
				},
			},
			pvzID:         "123e4567-e89b-12d3-a456-426614174000",
			expected:      nil,
			expectedError: true,
		},
	}

	for ttNum, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := reception.NewReceptionService(tt.mockRepo)
			got, err := s.CloseLastReception(context.Background(), tt.pvzID)

			if (err != nil) != tt.expectedError {
				t.Errorf("Test %v: ReceptionService.CloseLastReception() error = %v, expectedError %v", ttNum, err, tt.expectedError)
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Test %v: ReceptionService.CloseLastReception() = %v, expected %v", ttNum, got, tt.expected)
			}
		})
	}
}
