package service

import (
	"context"
	"errors"
	"time"

	"finance-tracker/internal/model"
	"finance-tracker/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepository *repository.UserRepository
	jwtSecret      string
}

func NewAuthService(userRepository *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepository: userRepository,
		jwtSecret:      jwtSecret,
	}
}

func (s *AuthService) Register(ctx context.Context, name, email, password string) (*model.User, error) {
	_, err := s.userRepository.GetByEmail(ctx, email)
	if err == nil {
		return nil, errors.New("user already exists")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(passwordHash),
	}

	if err := s.userRepository.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
