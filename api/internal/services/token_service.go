package services

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/Arlandaren/easyfund/internal/config"
	"github.com/Arlandaren/easyfund/internal/middleware"
	"github.com/Arlandaren/easyfund/internal/models"
)

type TokenService interface {
	GenerateToken(user *models.User) (string, error)
	ValidateToken(tokenString string) (*middleware.CustomClaims, error)
	RefreshToken(oldToken string) (string, error)
}

type tokenServiceImpl struct {
	config *config.JWTConfig
}

func NewTokenService(cfg *config.JWTConfig) TokenService {
	return &tokenServiceImpl{
		config: cfg,
	}
}

// GenerateToken генерирует JWT токен для пользователя
func (s *tokenServiceImpl) GenerateToken(user *models.User) (string, error) {
	now := time.Now()
	expirationTime := now.Add(s.config.Expiry)

	claims := &middleware.CustomClaims{
		UserID: user.UserID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "your_project_api",
			Subject:   user.UserID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.config.Secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken валидирует JWT токен и возвращает claims
func (s *tokenServiceImpl) ValidateToken(tokenString string) (*middleware.CustomClaims, error) {
	claims := &middleware.CustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Проверяем алгоритм подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	// Проверяем, что claims извлечены правильно
	parsedClaims, ok := token.Claims.(*middleware.CustomClaims)
	if !ok {
		return nil, fmt.Errorf("could not extract claims from token")
	}

	return parsedClaims, nil
}

// RefreshToken обновляет токен (выдаёт новый с новым временем expiry)
func (s *tokenServiceImpl) RefreshToken(oldToken string) (string, error) {
	claims, err := s.ValidateToken(oldToken)
	if err != nil {
		return "", fmt.Errorf("failed to validate old token: %w", err)
	}

	// Создаём новый токен с обновленным временем истечения
	user := &models.User{
		UserID: claims.UserID,
		Email:  claims.Email,
	}

	return s.GenerateToken(user)
}
