package auth_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kirillidk/pvz-service/internal/model"
	"github.com/kirillidk/pvz-service/internal/service/auth"
)

func TestGenerateAndValidateToken(t *testing.T) {
	tests := []struct {
		name          string
		role          model.UserRole
		jwtSecret     string
		expectedError bool
	}{
		{
			name:          "Employee Token",
			role:          model.EmployeeRole,
			jwtSecret:     "test-secret",
			expectedError: false,
		},
		{
			name:          "Moderator Token",
			role:          model.ModeratorRole,
			jwtSecret:     "test-secret",
			expectedError: false,
		},
		{
			name:          "Invalid Role",
			role:          "invalid-role",
			jwtSecret:     "test-secret",
			expectedError: false,
		},
		{
			name:          "Empty Secret",
			role:          model.EmployeeRole,
			jwtSecret:     "",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := auth.GenerateToken(tt.role, tt.jwtSecret)
			if (err != nil) != tt.expectedError {
				t.Errorf("GenerateToken() error = %v, expectedError %v", err, tt.expectedError)
				return
			}

			if !tt.expectedError && token == "" {
				t.Errorf("GenerateToken() returned empty token")
				return
			}

			if !tt.expectedError {
				claims, err := auth.ValidateToken(token, tt.jwtSecret)
				if err != nil {
					t.Errorf("ValidateToken() error = %v", err)
					return
				}

				if claims.Role != tt.role {
					t.Errorf("ValidateToken() role = %v, expected %v", claims.Role, tt.role)
				}
			}
		})
	}
}

func TestValidateToken_InvalidToken(t *testing.T) {
	tests := []struct {
		name      string
		token     string
		jwtSecret string
	}{
		{
			name:      "Invalid Format",
			token:     "invalid-token",
			jwtSecret: "test-secret",
		},
		{
			name:      "Wrong Secret",
			token:     "",
			jwtSecret: "wrong-secret",
		},
		{
			name:      "Expired Token",
			token:     "",
			jwtSecret: "test-secret",
		},
	}

	validToken, _ := auth.GenerateToken(model.EmployeeRole, "test-secret")
	tests[1].token = validToken

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tokenToTest string

			if tt.name == "Expired Token" {
				claims := auth.Claims{
					Role: model.EmployeeRole,
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
						IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
						NotBefore: jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
					},
				}

				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				expiredToken, _ := token.SignedString([]byte("test-secret"))
				tokenToTest = expiredToken
			} else {
				tokenToTest = tt.token
			}

			_, err := auth.ValidateToken(tokenToTest, tt.jwtSecret)
			if err == nil {
				t.Errorf("ValidateToken() expected error but got nil")
			}
		})
	}
}
