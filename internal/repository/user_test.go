package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kirillidk/pvz-service/internal/dto"
	"github.com/kirillidk/pvz-service/internal/model"
	"github.com/kirillidk/pvz-service/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	ctx := context.Background()

	tests := []struct {
		name          string
		registerReq   dto.RegisterRequest
		mockBehavior  func()
		expectedUser  *model.User
		expectedError error
	}{
		{
			name: "Success",
			registerReq: dto.RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
				Role:     model.EmployeeRole,
			},
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "email", "role"}).
					AddRow("123e4567-e89b-12d3-a456-426614174000", "test@example.com", "employee")

				mock.ExpectQuery(`INSERT INTO users`).
					WithArgs("test@example.com", sqlmock.AnyArg(), "employee").
					WillReturnRows(rows)
			},
			expectedUser: &model.User{
				ID:    "123e4567-e89b-12d3-a456-426614174000",
				Email: "test@example.com",
				Role:  model.EmployeeRole,
			},
			expectedError: nil,
		},
		{
			name: "DB Error",
			registerReq: dto.RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
				Role:     model.EmployeeRole,
			},
			mockBehavior: func() {
				mock.ExpectQuery(`INSERT INTO users`).
					WithArgs("test@example.com", sqlmock.AnyArg(), "employee").
					WillReturnError(errors.New("db error"))
			},
			expectedUser:  nil,
			expectedError: errors.New("failed to create user: db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			user, err := userRepo.CreateUser(ctx, tt.registerReq)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser, user)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUserRepository_FindUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	ctx := context.Background()

	tests := []struct {
		name             string
		email            string
		mockBehavior     func()
		expectedUser     *model.User
		expectedPassword string
		expectedError    error
	}{
		{
			name:  "Success",
			email: "test@example.com",
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "role"}).
					AddRow("123e4567-e89b-12d3-a456-426614174000", "test@example.com", "hashed_password", "employee")

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, email, password_hash, role FROM users WHERE email = $1`)).
					WithArgs("test@example.com").
					WillReturnRows(rows)
			},
			expectedUser: &model.User{
				ID:    "123e4567-e89b-12d3-a456-426614174000",
				Email: "test@example.com",
				Role:  model.EmployeeRole,
			},
			expectedPassword: "hashed_password",
			expectedError:    nil,
		},
		{
			name:  "User Not Found",
			email: "nonexistent@example.com",
			mockBehavior: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, email, password_hash, role FROM users WHERE email = $1`)).
					WithArgs("nonexistent@example.com").
					WillReturnError(sql.ErrNoRows)
			},
			expectedUser:     nil,
			expectedPassword: "",
			expectedError:    errors.New("user not found"),
		},
		{
			name:  "DB Error",
			email: "test@example.com",
			mockBehavior: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, email, password_hash, role FROM users WHERE email = $1`)).
					WithArgs("test@example.com").
					WillReturnError(errors.New("db error"))
			},
			expectedUser:     nil,
			expectedPassword: "",
			expectedError:    errors.New("failed to get user: db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			user, password, err := userRepo.FindUserByEmail(ctx, tt.email)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				assert.Nil(t, user)
				assert.Empty(t, password)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser, user)
				assert.Equal(t, tt.expectedPassword, password)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestUserRepository_UserExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	ctx := context.Background()

	tests := []struct {
		name           string
		email          string
		mockBehavior   func()
		expectedExists bool
		expectedError  error
	}{
		{
			name:  "User Exists",
			email: "existing@example.com",
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`)).
					WithArgs("existing@example.com").
					WillReturnRows(rows)
			},
			expectedExists: true,
			expectedError:  nil,
		},
		{
			name:  "User Does Not Exist",
			email: "nonexistent@example.com",
			mockBehavior: func() {
				rows := sqlmock.NewRows([]string{"exists"}).AddRow(false)

				mock.ExpectQuery(regexp.QuoteMeta(`SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`)).
					WithArgs("nonexistent@example.com").
					WillReturnRows(rows)
			},
			expectedExists: false,
			expectedError:  nil,
		},
		{
			name:  "DB Error",
			email: "test@example.com",
			mockBehavior: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`)).
					WithArgs("test@example.com").
					WillReturnError(errors.New("db error"))
			},
			expectedExists: false,
			expectedError:  errors.New("failed to check if user exists: db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			exists, err := userRepo.UserExists(ctx, tt.email)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				assert.False(t, exists)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedExists, exists)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
