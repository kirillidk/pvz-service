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

type MockReceptionService struct {
	CreateReceptionFunc    func(ctx context.Context, req dto.ReceptionCreateRequest) (*model.Reception, error)
	CloseLastReceptionFunc func(ctx context.Context, pvzID string) (*model.Reception, error)
}

func (m *MockReceptionService) CreateReception(ctx context.Context, req dto.ReceptionCreateRequest) (*model.Reception, error) {
	return m.CreateReceptionFunc(ctx, req)
}

func (m *MockReceptionService) CloseLastReception(ctx context.Context, pvzID string) (*model.Reception, error) {
	return m.CloseLastReceptionFunc(ctx, pvzID)
}

func TestReceptionHandler_CreateReception(t *testing.T) {
	testTime := time.Now()

	tests := []struct {
		name           string
		mockService    MockReceptionService
		requestBody    map[string]any
		expectedStatus int
		expectedBody   any
	}{
		{
			name: "Success",
			mockService: MockReceptionService{
				CreateReceptionFunc: func(ctx context.Context, req dto.ReceptionCreateRequest) (*model.Reception, error) {
					return &model.Reception{
						ID:       "123e4567-e89b-12d3-a456-426614174001",
						DateTime: testTime,
						PVZID:    req.PVZID,
						Status:   "in_progress",
					}, nil
				},
			},
			requestBody: map[string]any{
				"pvzId": "123e4567-e89b-12d3-a456-426614174002",
			},
			expectedStatus: http.StatusCreated,
			expectedBody: model.Reception{
				ID:       "123e4567-e89b-12d3-a456-426614174001",
				DateTime: testTime,
				PVZID:    "123e4567-e89b-12d3-a456-426614174002",
				Status:   "in_progress",
			},
		},
		{
			name: "Invalid Request Data",
			mockService: MockReceptionService{
				CreateReceptionFunc: func(ctx context.Context, req dto.ReceptionCreateRequest) (*model.Reception, error) {
					return nil, nil
				},
			},
			requestBody: map[string]any{
				"pvzId": "invalid-uuid",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: model.Error{
				Message: "Invalid request data",
			},
		},
		{
			name: "Service Error",
			mockService: MockReceptionService{
				CreateReceptionFunc: func(ctx context.Context, req dto.ReceptionCreateRequest) (*model.Reception, error) {
					return nil, errors.New("failed to create reception")
				},
			},
			requestBody: map[string]any{
				"pvzId": "123e4567-e89b-12d3-a456-426614174002",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: model.Error{
				Message: "failed to create reception",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			receptionHandler := handler.NewReceptionHandler(&tt.mockService)

			router.POST("/receptions", receptionHandler.CreateReception)

			requestBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/receptions", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var response any
			if tt.expectedStatus == http.StatusCreated {
				var reception model.Reception
				json.Unmarshal(w.Body.Bytes(), &reception)

				if !reception.DateTime.IsZero() {
					reception.DateTime = testTime
				}

				response = reception
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

func TestReceptionHandler_CloseLastReception(t *testing.T) {
	testTime := time.Now()

	tests := []struct {
		name           string
		mockService    MockReceptionService
		pvzID          string
		expectedStatus int
		expectedBody   any
	}{
		{
			name: "Success",
			mockService: MockReceptionService{
				CloseLastReceptionFunc: func(ctx context.Context, pvzID string) (*model.Reception, error) {
					return &model.Reception{
						ID:       "123e4567-e89b-12d3-a456-426614174001",
						DateTime: testTime,
						PVZID:    pvzID,
						Status:   "close",
					}, nil
				},
			},
			pvzID:          "123e4567-e89b-12d3-a456-426614174002",
			expectedStatus: http.StatusOK,
			expectedBody: model.Reception{
				ID:       "123e4567-e89b-12d3-a456-426614174001",
				DateTime: testTime,
				PVZID:    "123e4567-e89b-12d3-a456-426614174002",
				Status:   "close",
			},
		},
		{
			name: "Missing PVZ ID",
			mockService: MockReceptionService{
				CloseLastReceptionFunc: func(ctx context.Context, pvzID string) (*model.Reception, error) {
					return nil, nil
				},
			},
			pvzID:          "",
			expectedStatus: http.StatusBadRequest,
			expectedBody: model.Error{
				Message: "PVZ ID is required",
			},
		},
		{
			name: "Service Error",
			mockService: MockReceptionService{
				CloseLastReceptionFunc: func(ctx context.Context, pvzID string) (*model.Reception, error) {
					return nil, errors.New("failed to close reception")
				},
			},
			pvzID:          "123e4567-e89b-12d3-a456-426614174002",
			expectedStatus: http.StatusBadRequest,
			expectedBody: model.Error{
				Message: "failed to close reception",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			receptionHandler := handler.NewReceptionHandler(&tt.mockService)

			router.POST("/pvz/:pvzId/receptions/close", receptionHandler.CloseLastReception)

			url := "/pvz/" + tt.pvzID + "/receptions/close"
			req, _ := http.NewRequest(http.MethodPost, url, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var response any
			if tt.expectedStatus == http.StatusOK {
				var reception model.Reception
				json.Unmarshal(w.Body.Bytes(), &reception)

				if !reception.DateTime.IsZero() {
					reception.DateTime = testTime
				}

				response = reception
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
