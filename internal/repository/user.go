package repository

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/kirillidk/pvz-service/internal/dto"
	"github.com/kirillidk/pvz-service/internal/model"
	"golang.org/x/crypto/bcrypt"
)

const (
	usertableName = "users"
)

type UserRepositoryInterface interface {
	CreateUser(ctx context.Context, registerReq dto.RegisterRequest) (*model.User, error)
	FindUserByEmail(ctx context.Context, email string) (*model.User, string, error)
	UserExists(ctx context.Context, email string) (bool, error)
}

type UserRepository struct {
	db   *sql.DB
	psql sq.StatementBuilderType
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db:   db,
		psql: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, registerReq dto.RegisterRequest) (*model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerReq.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	query, args, err := r.psql.
		Insert(usertableName).
		Columns("email", "password_hash", "role").
		Values(registerReq.Email, string(hashedPassword), registerReq.Role).
		Suffix("RETURNING id, email, role").
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build sql query: %w", err)
	}

	var user model.User
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.Email, &user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) FindUserByEmail(ctx context.Context, email string) (*model.User, string, error) {
	var user model.User
	var passwordHash string

	query, args, err := r.psql.
		Select("id", "email", "password_hash", "role").
		From(usertableName).
		Where(sq.Eq{"email": email}).
		ToSql()

	if err != nil {
		return nil, "", fmt.Errorf("failed to build sql query: %w", err)
	}

	err = r.db.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.Email, &passwordHash, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, "", fmt.Errorf("user not found")
		}
		return nil, "", fmt.Errorf("failed to get user: %w", err)
	}

	return &user, passwordHash, nil
}

func (r *UserRepository) UserExists(ctx context.Context, email string) (bool, error) {
	var exists bool

	query, _, err := r.psql.
		Select("EXISTS(SELECT 1 FROM users WHERE email = $1)").
		ToSql()

	if err != nil {
		return false, fmt.Errorf("failed to build sql query: %w", err)
	}

	err = r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if user exists: %w", err)
	}

	return exists, nil
}
