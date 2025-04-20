package auth_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/kirillidk/pvz-service/internal/dto"
	"github.com/kirillidk/pvz-service/internal/model"
	"github.com/kirillidk/pvz-service/internal/service/auth"
	"golang.org/x/crypto/bcrypt"
)

type MockUserRepository struct {
	CreateUserFunc      func(ctx context.Context, req dto.RegisterRequest) (*model.User, error)
	FindUserByEmailFunc func(ctx context.Context, email string) (*model.User, string, error)
	UserExistsFunc      func(ctx context.Context, email string) (bool, error)
}

func (m *MockUserRepository) CreateUser(ctx context.Context, req dto.RegisterRequest) (*model.User, error) {
	return m.CreateUserFunc(ctx, req)
}

func (m *MockUserRepository) FindUserByEmail(ctx context.Context, email string) (*model.User, string, error) {
	return m.FindUserByEmailFunc(ctx, email)
}

func (m *MockUserRepository) UserExists(ctx context.Context, email string) (bool, error) {
	return m.UserExistsFunc(ctx, email)
}

func TestAuthService_Register(t *testing.T) {
	tests := []struct {
		name          string
		mockRepo      *MockUserRepository
		input         dto.RegisterRequest
		expected      *model.User
		expectedError bool
	}{
		{
			name: "Success",
			mockRepo: &MockUserRepository{
				UserExistsFunc: func(ctx context.Context, email string) (bool, error) {
					return false, nil
				},
				CreateUserFunc: func(ctx context.Context, req dto.RegisterRequest) (*model.User, error) {
					return &model.User{
						ID:    "123e4567-e89b-12d3-a456-426614174000",
						Email: req.Email,
						Role:  req.Role,
					}, nil
				},
			},
			input: dto.RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
				Role:     model.EmployeeRole,
			},
			expected: &model.User{
				ID:    "123e4567-e89b-12d3-a456-426614174000",
				Email: "test@example.com",
				Role:  model.EmployeeRole,
			},
			expectedError: false,
		},
		{
			name: "User Exists",
			mockRepo: &MockUserRepository{
				UserExistsFunc: func(ctx context.Context, email string) (bool, error) {
					return true, nil
				},
			},
			input: dto.RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
				Role:     model.EmployeeRole,
			},
			expected:      nil,
			expectedError: true,
		},
		{
			name: "Repository Error",
			mockRepo: &MockUserRepository{
				UserExistsFunc: func(ctx context.Context, email string) (bool, error) {
					return false, nil
				},
				CreateUserFunc: func(ctx context.Context, req dto.RegisterRequest) (*model.User, error) {
					return nil, errors.New("repository error")
				},
			},
			input: dto.RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
				Role:     model.EmployeeRole,
			},
			expected:      nil,
			expectedError: true,
		},
	}

	for tNum, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := auth.NewAuthService(tt.mockRepo, "secret")
			got, err := s.Register(context.Background(), tt.input)

			if (err != nil) != tt.expectedError {
				t.Errorf("Test %v: AuthService.Register() error = %v, expectedError %v", tNum, err, tt.expectedError)
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Test %v: AuthService.Register() = %v, expected %v", tNum, got, tt.expected)
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	validPasswordHash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		t.Errorf("failed to hash password: %v", err)
	}

	tests := []struct {
		name          string
		mockRepo      *MockUserRepository
		input         dto.LoginRequest
		jwtSecret     string
		expectedError bool
	}{
		{
			name: "Success",
			mockRepo: &MockUserRepository{
				FindUserByEmailFunc: func(ctx context.Context, email string) (*model.User, string, error) {
					return &model.User{
						ID:    "123e4567-e89b-12d3-a456-426614174000",
						Email: email,
						Role:  model.EmployeeRole,
					}, string(validPasswordHash), nil
				},
			},
			input: dto.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			jwtSecret:     "test-secret",
			expectedError: false,
		},
		{
			name: "User Not Found",
			mockRepo: &MockUserRepository{
				FindUserByEmailFunc: func(ctx context.Context, email string) (*model.User, string, error) {
					return nil, "", errors.New("user not found")
				},
			},
			input: dto.LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			jwtSecret:     "test-secret",
			expectedError: true,
		},
		{
			name: "Invalid Password",
			mockRepo: &MockUserRepository{
				FindUserByEmailFunc: func(ctx context.Context, email string) (*model.User, string, error) {
					return &model.User{
						ID:    "123e4567-e89b-12d3-a456-426614174000",
						Email: email,
						Role:  model.EmployeeRole,
					}, string(validPasswordHash), nil
				},
			},
			input: dto.LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			jwtSecret:     "test-secret",
			expectedError: true,
		},
	}

	for ttNum, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := auth.NewAuthService(tt.mockRepo, tt.jwtSecret)
			_, err := s.Login(context.Background(), tt.input)

			if (err != nil) != tt.expectedError {
				t.Errorf("Test %v: AuthService.Login() error = %v, expectedError %v", ttNum, err, tt.expectedError)
			}
		})
	}
}
