package usecase

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/Arlandaren/easyfund/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type AuthUseCase struct {
	repo      *repository.UserRepository
	jwtSecret string
}

func NewAuthUseCase(repo *repository.UserRepository, jwtSecret string) *AuthUseCase {
	return &AuthUseCase{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

func (uc *AuthUseCase) Register(ctx context.Context, email, password, role string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = uc.repo.Create(ctx, email, string(hash), role)
	return err
}

func (uc *AuthUseCase) Login(ctx context.Context, email, password string) (string, error) {
	user, err := uc.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	return token.SignedString([]byte(uc.jwtSecret))
}
