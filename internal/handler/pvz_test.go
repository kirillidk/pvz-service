package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kirillidk/pvz-service/internal/dto"
	"github.com/kirillidk/pvz-service/internal/handler"
	"github.com/kirillidk/pvz-service/internal/model"
)

type MockPVZService struct {
	CreatePVZFunc  func(ctx context.Context, req dto.PVZCreateRequest) (*model.PVZ, error)
	GetPVZListFunc func(ctx context.Context, filter dto.PVZFilterQuery) (*dto.PaginatedResponse, error)
}

func (m *MockPVZService) CreatePVZ(ctx context.Context, req dto.PVZCreateRequest) (*model.PVZ, error) {
	return m.CreatePVZFunc(ctx, req)
}

func (m *MockPVZService) GetPVZList(ctx context.Context, filter dto.PVZFilterQuery) (*dto.PaginatedResponse, error) {
	return m.GetPVZListFunc(ctx, filter)
}

func TestPVZHandler_CreatePVZ(t *testing.T) {
	testTime := time.Date(2025, 4, 15, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		mockService    MockPVZService
		requestBody    map[string]any
		expectedStatus int
		expectedBody   any
	}{
		{
			name: "Success",
			mockService: MockPVZService{
				CreatePVZFunc: func(ctx context.Context, req dto.PVZCreateRequest) (*model.PVZ, error) {
					return &model.PVZ{
						ID:               "123e4567-e89b-12d3-a456-426614174001",
						RegistrationDate: req.RegistrationDate,
						City:             req.City,
					}, nil
				},
			},
			requestBody: map[string]any{
				"registrationDate": testTime,
				"city":             "Москва",
			},
			expectedStatus: http.StatusCreated,
			expectedBody: model.PVZ{
				ID:               "123e4567-e89b-12d3-a456-426614174001",
				RegistrationDate: testTime,
				City:             "Москва",
			},
		},
		{
			name: "Invalid City",
			mockService: MockPVZService{
				CreatePVZFunc: func(ctx context.Context, req dto.PVZCreateRequest) (*model.PVZ, error) {
					return nil, nil
				},
			},
			requestBody: map[string]any{
				"registrationDate": testTime,
				"city":             "Новосибирск",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: model.Error{
				Message: "Invalid request data",
			},
		},
		{
			name: "Invalid Request Data",
			mockService: MockPVZService{
				CreatePVZFunc: func(ctx context.Context, req dto.PVZCreateRequest) (*model.PVZ, error) {
					return nil, nil
				},
			},
			requestBody: map[string]any{
				"registrationDate": testTime,
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: model.Error{
				Message: "Invalid request data",
			},
		},
		{
			name: "Service Error",
			mockService: MockPVZService{
				CreatePVZFunc: func(ctx context.Context, req dto.PVZCreateRequest) (*model.PVZ, error) {
					return nil, errors.New("failed to create PVZ")
				},
			},
			requestBody: map[string]any{
				"registrationDate": testTime,
				"city":             "Москва",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: model.Error{
				Message: "failed to create PVZ",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			pvzHandler := handler.NewPVZHandler(&tt.mockService)

			router.POST("/pvz", pvzHandler.CreatePVZ)

			requestBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/pvz", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var response any
			if tt.expectedStatus == http.StatusCreated {
				var pvz model.PVZ
				json.Unmarshal(w.Body.Bytes(), &pvz)
				response = pvz
			} else {
				var errResponse model.Error
				json.Unmarshal(w.Body.Bytes(), &errResponse)
				response = errResponse
			}

			if !reflect.DeepEqual(tt.expectedBody, response) {
				t.Errorf("Expected body %v, got %v", tt.expectedBody, response)
			}
		})
	}
}

func TestPVZHandler_GetPVZList(t *testing.T) {
	testTime := time.Date(2025, 4, 15, 12, 0, 0, 0, time.UTC)

	pvz := model.PVZ{
		ID:               "123e4567-e89b-12d3-a456-426614174001",
		RegistrationDate: testTime,
		City:             "Москва",
	}

	reception := model.Reception{
		ID:       "123e4567-e89b-12d3-a456-426614174002",
		DateTime: testTime,
		PVZID:    pvz.ID,
		Status:   "in_progress",
	}

	product := model.Product{
		ID:          "123e4567-e89b-12d3-a456-426614174003",
		DateTime:    testTime,
		Type:        "электроника",
		ReceptionID: reception.ID,
	}

	tests := []struct {
		name           string
		mockService    MockPVZService
		queryParams    string
		expectedStatus int
		expectedBody   any
	}{
		{
			name: "Success",
			mockService: MockPVZService{
				GetPVZListFunc: func(ctx context.Context, filter dto.PVZFilterQuery) (*dto.PaginatedResponse, error) {
					return &dto.PaginatedResponse{
						Data: []dto.PVZWithReceptionsResponse{
							{
								PVZ: pvz,
								Receptions: []dto.ReceptionWithProductsResponse{
									{
										Reception: reception,
										Products:  []model.Product{product},
									},
								},
							},
						},
					}, nil
				},
			},
			queryParams:    "?page=1&limit=10",
			expectedStatus: http.StatusOK,
			expectedBody: &dto.PaginatedResponse{
				Data: []dto.PVZWithReceptionsResponse{
					{
						PVZ: pvz,
						Receptions: []dto.ReceptionWithProductsResponse{
							{
								Reception: reception,
								Products:  []model.Product{product},
							},
						},
					},
				},
			},
		},
		{
			name: "Invalid Query Parameters",
			mockService: MockPVZService{
				GetPVZListFunc: func(ctx context.Context, filter dto.PVZFilterQuery) (*dto.PaginatedResponse, error) {
					return nil, nil
				},
			},
			queryParams:    "?page=0&limit=50",
			expectedStatus: http.StatusBadRequest,
			expectedBody: model.Error{
				Message: "Invalid query parameters",
			},
		},
		{
			name: "Service Error",
			mockService: MockPVZService{
				GetPVZListFunc: func(ctx context.Context, filter dto.PVZFilterQuery) (*dto.PaginatedResponse, error) {
					return nil, errors.New("failed to get PVZ list")
				},
			},
			queryParams:    "?page=1&limit=10",
			expectedStatus: http.StatusInternalServerError,
			expectedBody: model.Error{
				Message: "failed to get PVZ list",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			pvzHandler := handler.NewPVZHandler(&tt.mockService)

			router.GET("/pvz", pvzHandler.GetPVZList)

			req, _ := http.NewRequest(http.MethodGet, "/pvz"+tt.queryParams, nil)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var response any
			if tt.expectedStatus == http.StatusOK {
				var paginatedResponse dto.PaginatedResponse
				json.Unmarshal(w.Body.Bytes(), &paginatedResponse)
				response = &paginatedResponse
			} else {
				var errResponse model.Error
				json.Unmarshal(w.Body.Bytes(), &errResponse)
				response = errResponse
			}

			if !reflect.DeepEqual(tt.expectedBody, response) {
				t.Errorf("Expected body %v, got %v", tt.expectedBody, response)
			}
		})
	}
}
