package middleware_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kirillidk/pvz-service/internal/middleware"
	"github.com/kirillidk/pvz-service/internal/model"
	"github.com/kirillidk/pvz-service/internal/service/auth"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func TestAuthMiddleware(t *testing.T) {
	jwtSecret := "test-secret"

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectedBody   any
	}{
		{
			name:           "No Authorization Header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody: model.Error{
				Message: "Authorization header is required",
			},
		},
		{
			name:           "Invalid Authorization Format",
			authHeader:     "InvalidFormat token",
			expectedStatus: http.StatusUnauthorized,
			expectedBody: model.Error{
				Message: "Authorization header must be in format: Bearer {token}",
			},
		},
		{
			name:           "Invalid Token",
			authHeader:     "Bearer invalid-token",
			expectedStatus: http.StatusUnauthorized,
			expectedBody: model.Error{
				Message: "Invalid or expired token",
			},
		},
		{
			name:           "Valid Employee Token",
			authHeader:     "",
			expectedStatus: http.StatusOK,
			expectedBody:   nil,
		},
		{
			name:           "Valid Moderator Token",
			authHeader:     "",
			expectedStatus: http.StatusOK,
			expectedBody:   nil,
		},
	}

	employeeToken, _ := auth.GenerateToken(model.EmployeeRole, jwtSecret)
	moderatorToken, _ := auth.GenerateToken(model.ModeratorRole, jwtSecret)

	tests[3].authHeader = "Bearer " + employeeToken
	tests[4].authHeader = "Bearer " + moderatorToken

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(middleware.AuthMiddleware(jwtSecret))

			router.GET("/protected", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedStatus != http.StatusOK && tt.expectedBody != nil {
				var responseError model.Error
				if err := json.Unmarshal(w.Body.Bytes(), &responseError); err != nil {
					t.Errorf("Failed to unmarshal response body: %v", err)
				}

				expectedError, ok := tt.expectedBody.(model.Error)
				if !ok {
					t.Errorf("Expected body is not of type model.Error")
				}

				if responseError.Message != expectedError.Message {
					t.Errorf("Expected error message '%s', got '%s'", expectedError.Message, responseError.Message)
				}
			}
		})
	}
}

func TestRoleMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		role           model.UserRole
		requiredRoles  []model.UserRole
		expectedStatus int
		expectedBody   any
	}{
		{
			name:           "Employee can access Employee-only route",
			role:           model.EmployeeRole,
			requiredRoles:  []model.UserRole{model.EmployeeRole},
			expectedStatus: http.StatusOK,
			expectedBody:   nil,
		},
		{
			name:           "Employee cannot access Moderator-only route",
			role:           model.EmployeeRole,
			requiredRoles:  []model.UserRole{model.ModeratorRole},
			expectedStatus: http.StatusForbidden,
			expectedBody: model.Error{
				Message: "Operation not permitted for this user role",
			},
		},
		{
			name:           "Moderator can access Moderator-only route",
			role:           model.ModeratorRole,
			requiredRoles:  []model.UserRole{model.ModeratorRole},
			expectedStatus: http.StatusOK,
			expectedBody:   nil,
		},
		{
			name:           "Any role can access multi-role route",
			role:           model.EmployeeRole,
			requiredRoles:  []model.UserRole{model.EmployeeRole, model.ModeratorRole},
			expectedStatus: http.StatusOK,
			expectedBody:   nil,
		},
		{
			name:           "Missing user role",
			role:           "",
			requiredRoles:  []model.UserRole{model.EmployeeRole},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: model.Error{
				Message: "User not authenticated",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()

			router.Use(func(c *gin.Context) {
				if tt.role != "" {
					c.Set("userRole", tt.role)
				}
				c.Next()
			})

			router.Use(middleware.RoleMiddleware(tt.requiredRoles...))

			router.GET("/protected", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedStatus != http.StatusOK && tt.expectedBody != nil {
				var responseError model.Error
				if err := json.Unmarshal(w.Body.Bytes(), &responseError); err != nil {
					t.Errorf("Failed to unmarshal response body: %v", err)
				}

				expectedError, ok := tt.expectedBody.(model.Error)
				if !ok {
					t.Errorf("Expected body is not of type model.Error")
				}

				if responseError.Message != expectedError.Message {
					t.Errorf("Expected error message '%s', got '%s'", expectedError.Message, responseError.Message)
				}
			}
		})
	}
}
