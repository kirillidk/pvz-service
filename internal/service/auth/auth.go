package auth

import (
	"errors"
	"fmt"

	"github.com/kirillidk/pvz-service/internal/dto"
	"github.com/kirillidk/pvz-service/internal/model"
	"github.com/kirillidk/pvz-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepository *repository.UserRepository
	jwtSecret      string
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepository: userRepo,
		jwtSecret:      jwtSecret,
	}
}

func (s *AuthService) Register(registerReq dto.RegisterRequest) (*model.User, error) {
	exists, err := s.userRepository.UserExists(registerReq.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check if user exists: %w", err)
	}

	if exists {
		return nil, errors.New("user with this email already exists")
	}

	user, err := s.userRepository.CreateUser(registerReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *AuthService) Login(loginReq dto.LoginRequest) (string, error) {
	user, passwordHash, err := s.userRepository.FindUserByEmail(loginReq.Email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(loginReq.Password))
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	token, err := GenerateToken(user.Role, s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}
