package grpc_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	pvz_v1 "github.com/kirillidk/pvz-service/api/proto/pvz/pvz_v1"
	"github.com/kirillidk/pvz-service/internal/dto"
	"github.com/kirillidk/pvz-service/internal/model"
	grpcservice "github.com/kirillidk/pvz-service/internal/service/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

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

func TestPVZService_GetPVZList(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name          string
		mockRepo      *MockPVZRepository
		request       *pvz_v1.GetPVZListRequest
		expected      *pvz_v1.GetPVZListResponse
		expectedError bool
	}{
		{
			name: "Success - Multiple PVZs",
			mockRepo: &MockPVZRepository{
				GetPVZListFunc: func(ctx context.Context, filter dto.PVZFilterQuery) ([]model.PVZ, error) {
					return []model.PVZ{
						{
							ID:               "pvz-id-1",
							RegistrationDate: now,
							City:             "Москва",
						},
						{
							ID:               "pvz-id-2",
							RegistrationDate: now.Add(-1 * time.Hour),
							City:             "Санкт-Петербург",
						},
					}, nil
				},
			},
			request: &pvz_v1.GetPVZListRequest{},
			expected: &pvz_v1.GetPVZListResponse{
				Pvzs: []*pvz_v1.PVZ{
					{
						Id:               "pvz-id-1",
						RegistrationDate: timestamppb.New(now),
						City:             "Москва",
					},
					{
						Id:               "pvz-id-2",
						RegistrationDate: timestamppb.New(now.Add(-1 * time.Hour)),
						City:             "Санкт-Петербург",
					},
				},
			},
			expectedError: false,
		},
		{
			name: "Success - Single PVZ",
			mockRepo: &MockPVZRepository{
				GetPVZListFunc: func(ctx context.Context, filter dto.PVZFilterQuery) ([]model.PVZ, error) {
					return []model.PVZ{
						{
							ID:               "pvz-id-1",
							RegistrationDate: now,
							City:             "Москва",
						},
					}, nil
				},
			},
			request: &pvz_v1.GetPVZListRequest{},
			expected: &pvz_v1.GetPVZListResponse{
				Pvzs: []*pvz_v1.PVZ{
					{
						Id:               "pvz-id-1",
						RegistrationDate: timestamppb.New(now),
						City:             "Москва",
					},
				},
			},
			expectedError: false,
		},
		{
			name: "Success - Empty PVZ List",
			mockRepo: &MockPVZRepository{
				GetPVZListFunc: func(ctx context.Context, filter dto.PVZFilterQuery) ([]model.PVZ, error) {
					return []model.PVZ{}, nil
				},
			},
			request: &pvz_v1.GetPVZListRequest{},
			expected: &pvz_v1.GetPVZListResponse{
				Pvzs: []*pvz_v1.PVZ{},
			},
			expectedError: false,
		},
		{
			name: "Repository Error",
			mockRepo: &MockPVZRepository{
				GetPVZListFunc: func(ctx context.Context, filter dto.PVZFilterQuery) ([]model.PVZ, error) {
					return nil, errors.New("repository error")
				},
			},
			request:       &pvz_v1.GetPVZListRequest{},
			expected:      nil,
			expectedError: true,
		},
		{
			name: "Filter Verification",
			mockRepo: &MockPVZRepository{
				GetPVZListFunc: func(ctx context.Context, filter dto.PVZFilterQuery) ([]model.PVZ, error) {
					if filter.Page != 1 || filter.Limit != 1000 {
						t.Errorf("Expected page 1, limit 1000, got page %d, limit %d", filter.Page, filter.Limit)
					}
					return []model.PVZ{
						{
							ID:               "pvz-id-1",
							RegistrationDate: now,
							City:             "Москва",
						},
					}, nil
				},
			},
			request: &pvz_v1.GetPVZListRequest{},
			expected: &pvz_v1.GetPVZListResponse{
				Pvzs: []*pvz_v1.PVZ{
					{
						Id:               "pvz-id-1",
						RegistrationDate: timestamppb.New(now),
						City:             "Москва",
					},
				},
			},
			expectedError: false,
		},
	}

	for ttNum, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := grpcservice.NewPVZService(tt.mockRepo)

			got, err := s.GetPVZList(context.Background(), tt.request)

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
