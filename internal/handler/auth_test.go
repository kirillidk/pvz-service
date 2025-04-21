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

type MockAuthService struct {
	RegisterFunc func(ctx context.Context, req dto.RegisterRequest) (*model.User, error)
	LoginFunc    func(ctx context.Context, req dto.LoginRequest) (string, error)
}

func (m *MockAuthService) Register(ctx context.Context, req dto.RegisterRequest) (*model.User, error) {
	return m.RegisterFunc(ctx, req)
}

func (m *MockAuthService) Login(ctx context.Context, req dto.LoginRequest) (string, error) {
	return m.LoginFunc(ctx, req)
}

func TestAuthHandler_DummyLogin(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]any
		expectedStatus int
		expectedBody   any
	}{
		{
			name: "Success - Employee Role",
			requestBody: map[string]any{
				"role": "employee",
			},
			expectedStatus: http.StatusOK,
			expectedBody: model.Token{
				Value: "mock-token",
			},
		},
		{
			name: "Success - Moderator Role",
			requestBody: map[string]any{
				"role": "moderator",
			},
			expectedStatus: http.StatusOK,
			expectedBody: model.Token{
				Value: "mock-token",
			},
		},
		{
			name: "Invalid Role",
			requestBody: map[string]any{
				"role": "invalid_role",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: model.Error{
				Message: "Invalid role. Must be 'employee' or 'moderator'",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			authHandler := handler.NewAuthHandler(&MockAuthService{}, "test-secret")

			router.POST("/dummy-login", authHandler.DummyLogin)

			requestBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/dummy-login", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var tokenResponse model.Token
				json.Unmarshal(w.Body.Bytes(), &tokenResponse)

				if tokenResponse.Value == "" {
					t.Errorf("Expected token value to be non-empty")
				}

				authHeader := w.Header().Get("Authorization")
				if authHeader == "" || authHeader != "Bearer "+tokenResponse.Value {
					t.Errorf("Expected Authorization header to be 'Bearer %s', got %s", tokenResponse.Value, authHeader)
				}
			} else {
				var errResponse model.Error
				json.Unmarshal(w.Body.Bytes(), &errResponse)

				if !reflect.DeepEqual(tt.expectedBody, errResponse) {
					t.Errorf("Expected body %v, got %v", tt.expectedBody, errResponse)
				}
			}
		})
	}
}

func TestAuthHandler_Register(t *testing.T) {
	tests := []struct {
		name           string
		mockService    MockAuthService
		requestBody    map[string]any
		expectedStatus int
		expectedBody   any
	}{
		{
			name: "Success",
			mockService: MockAuthService{
				RegisterFunc: func(ctx context.Context, req dto.RegisterRequest) (*model.User, error) {
					return &model.User{
						ID:    "123e4567-e89b-12d3-a456-426614174001",
						Email: req.Email,
						Role:  req.Role,
					}, nil
				},
			},
			requestBody: map[string]any{
				"email":    "test@example.com",
				"password": "password123",
				"role":     "employee",
			},
			expectedStatus: http.StatusCreated,
			expectedBody: model.User{
				ID:    "123e4567-e89b-12d3-a456-426614174001",
				Email: "test@example.com",
				Role:  model.EmployeeRole,
			},
		},
		{
			name: "Invalid Request Data",
			mockService: MockAuthService{
				RegisterFunc: func(ctx context.Context, req dto.RegisterRequest) (*model.User, error) {
					return nil, nil
				},
			},
			requestBody: map[string]any{
				"email":    "invalid-email",
				"password": "pwd",
				"role":     "employee",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: model.Error{
				Message: "Invalid request data",
			},
		},
		{
			name: "Service Error",
			mockService: MockAuthService{
				RegisterFunc: func(ctx context.Context, req dto.RegisterRequest) (*model.User, error) {
					return nil, errors.New("user with this email already exists")
				},
			},
			requestBody: map[string]any{
				"email":    "existing@example.com",
				"password": "password123",
				"role":     "employee",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: model.Error{
				Message: "user with this email already exists",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			authHandler := handler.NewAuthHandler(&tt.mockService, "test-secret")

			router.POST("/register", authHandler.Register)

			requestBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var response any
			if tt.expectedStatus == http.StatusCreated {
				var user model.User
				json.Unmarshal(w.Body.Bytes(), &user)
				response = user
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

func TestAuthHandler_Login(t *testing.T) {
	tests := []struct {
		name           string
		mockService    MockAuthService
		requestBody    map[string]any
		expectedStatus int
		expectedBody   any
	}{
		{
			name: "Success",
			mockService: MockAuthService{
				LoginFunc: func(ctx context.Context, req dto.LoginRequest) (string, error) {
					return "valid-token", nil
				},
			},
			requestBody: map[string]any{
				"email":    "test@example.com",
				"password": "password123",
			},
			expectedStatus: http.StatusOK,
			expectedBody: model.Token{
				Value: "valid-token",
			},
		},
		{
			name: "Invalid Request Data",
			mockService: MockAuthService{
				LoginFunc: func(ctx context.Context, req dto.LoginRequest) (string, error) {
					return "", nil
				},
			},
			requestBody: map[string]any{
				"email":    "invalid-email",
				"password": "pwd",
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: model.Error{
				Message: "Invalid request data",
			},
		},
		{
			name: "Invalid Credentials",
			mockService: MockAuthService{
				LoginFunc: func(ctx context.Context, req dto.LoginRequest) (string, error) {
					return "", errors.New("invalid email or password")
				},
			},
			requestBody: map[string]any{
				"email":    "test@example.com",
				"password": "wrong-password",
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: model.Error{
				Message: "invalid email or password",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			authHandler := handler.NewAuthHandler(&tt.mockService, "test-secret")

			router.POST("/login", authHandler.Login)

			requestBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var response any
			if tt.expectedStatus == http.StatusOK {
				var token model.Token
				json.Unmarshal(w.Body.Bytes(), &token)
				response = token

				authHeader := w.Header().Get("Authorization")
				if authHeader == "" || authHeader != "Bearer "+token.Value {
					t.Errorf("Expected Authorization header to be 'Bearer %s', got %s", token.Value, authHeader)
				}
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
