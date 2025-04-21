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

	"github.com/gin-gonic/gin"
	"github.com/kirillidk/pvz-service/internal/dto"
	"github.com/kirillidk/pvz-service/internal/handler"
	"github.com/kirillidk/pvz-service/internal/model"
)

type MockProductService struct {
	CreateProductFunc     func(ctx context.Context, req dto.ProductCreateRequest) (*model.Product, error)
	DeleteLastProductFunc func(ctx context.Context, pvzID string) error
}

func (m *MockProductService) CreateProduct(ctx context.Context, req dto.ProductCreateRequest) (*model.Product, error) {
	return m.CreateProductFunc(ctx, req)
}

func (m *MockProductService) DeleteLastProduct(ctx context.Context, pvzID string) error {
	return m.DeleteLastProductFunc(ctx, pvzID)
}

func TestProductHandler_CreateProduct(t *testing.T) {
	tests := []struct {
		name           string
		mockService    MockProductService
		requestBody    map[string]any
		expectedStatus int
		expectedBody   any
	}{
		{
			name: "Success",
			mockService: MockProductService{
				CreateProductFunc: func(ctx context.Context, req dto.ProductCreateRequest) (*model.Product, error) {
					return &model.Product{
						ID:          "123e4567-e89b-12d3-a456-426614174001",
						Type:        req.Type,
						ReceptionID: "123e4567-e89b-12d3-a456-426614174002",
					}, nil
				},
			},
			requestBody: map[string]any{
				"type":  "электроника",
				"pvzId": "123e4567-e89b-12d3-a456-426614174003",
			},
			expectedStatus: http.StatusCreated,
			expectedBody: model.Product{
				ID:          "123e4567-e89b-12d3-a456-426614174001",
				Type:        "электроника",
				ReceptionID: "123e4567-e89b-12d3-a456-426614174002",
			},
		},
		{
			name: "Invalid Request Data",
			mockService: MockProductService{
				CreateProductFunc: func(ctx context.Context, req dto.ProductCreateRequest) (*model.Product, error) {
					return nil, nil
				},
			},
			requestBody: map[string]any{
				"type": "invalid_type",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: model.Error{
				Message: "Invalid request data",
			},
		},
		{
			name: "Service Error",
			mockService: MockProductService{
				CreateProductFunc: func(ctx context.Context, req dto.ProductCreateRequest) (*model.Product, error) {
					return nil, errors.New("failed to create product")
				},
			},
			requestBody: map[string]any{
				"type":  "электроника",
				"pvzId": "123e4567-e89b-12d3-a456-426614174003",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: model.Error{
				Message: "failed to create product",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			productHandler := handler.NewProductHandler(&tt.mockService)

			router.POST("/products", productHandler.CreateProduct)

			requestBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var response any
			if tt.expectedStatus == http.StatusCreated {
				var product model.Product
				json.Unmarshal(w.Body.Bytes(), &product)
				response = product
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

func TestProductHandler_DeleteLastProduct(t *testing.T) {
	tests := []struct {
		name            string
		mockService     MockProductService
		pvzID           string
		expectedStatus  int
		expectedMessage string
	}{
		{
			name: "Success",
			mockService: MockProductService{
				DeleteLastProductFunc: func(ctx context.Context, pvzID string) error {
					return nil
				},
			},
			pvzID:           "123e4567-e89b-12d3-a456-426614174003",
			expectedStatus:  http.StatusOK,
			expectedMessage: "Last product deleted successfully",
		},
		{
			name: "Missing PVZ ID",
			mockService: MockProductService{
				DeleteLastProductFunc: func(ctx context.Context, pvzID string) error {
					return nil
				},
			},
			pvzID:           "",
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "PVZ ID is required",
		},
		{
			name: "Service Error",
			mockService: MockProductService{
				DeleteLastProductFunc: func(ctx context.Context, pvzID string) error {
					return errors.New("failed to delete last product")
				},
			},
			pvzID:           "123e4567-e89b-12d3-a456-426614174003",
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "failed to delete last product",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			productHandler := handler.NewProductHandler(&tt.mockService)

			router.DELETE("/pvz/:pvzId/products/last", productHandler.DeleteLastProduct)

			url := "/pvz/" + tt.pvzID + "/products/last"
			if tt.pvzID == "" {
				url = "/pvz//products/last"
			}

			req, _ := http.NewRequest(http.MethodDelete, url, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var responseMap map[string]string
			if err := json.Unmarshal(w.Body.Bytes(), &responseMap); err != nil {
				t.Errorf("Failed to parse response: %v", err)
				return
			}

			message, exists := responseMap["message"]
			if !exists && tt.expectedStatus == http.StatusOK {
				t.Errorf("Expected message field in response")
				return
			}

			if tt.expectedStatus == http.StatusBadRequest {
				errorMessage, exists := responseMap["message"]
				if !exists {
					t.Errorf("Expected error message in response")
					return
				}
				if errorMessage != tt.expectedMessage {
					t.Errorf("Expected error message %s, got %s", tt.expectedMessage, errorMessage)
				}
			} else if message != tt.expectedMessage {
				t.Errorf("Expected message %s, got %s", tt.expectedMessage, message)
			}
		})
	}
}
