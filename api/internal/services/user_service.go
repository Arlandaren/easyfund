package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/Arlandaren/easyfund/internal/models"
	"github.com/Arlandaren/easyfund/internal/repos"
)

type UserService interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	GetRandomUser(ctx context.Context) (*models.User, error)
	ListUsers(ctx context.Context, limit int) ([]models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error

	// Новый метод
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdatePasswordHash(ctx context.Context, userID uuid.UUID, hash string) error
}

type userServiceImpl struct {
	repo repos.UserRepository
}

func NewUserService(repo repos.UserRepository) UserService {
	return &userServiceImpl{repo: repo}
}

func (s *userServiceImpl) CreateUser(ctx context.Context, user *models.User) error {
	if user.UserID == uuid.Nil {
		user.UserID = uuid.New()
	}
	return s.repo.CreateUser(ctx, user)
}

func (s *userServiceImpl) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return s.repo.GetUserByID(ctx, userID)
}

func (s *userServiceImpl) GetRandomUser(ctx context.Context) (*models.User, error) {
	return s.repo.GetRandomUser(ctx)
}

func (s *userServiceImpl) ListUsers(ctx context.Context, limit int) ([]models.User, error) {
	return s.repo.ListUsers(ctx, limit)
}

func (s *userServiceImpl) UpdateUser(ctx context.Context, user *models.User) error {
	return s.repo.UpdateUser(ctx, user)
}

func (s *userServiceImpl) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	return s.repo.DeleteUser(ctx, userID)
}

func (s *userServiceImpl) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.repo.GetUserByEmail(ctx, email)
}

func (s *userServiceImpl) UpdatePasswordHash(ctx context.Context, userID uuid.UUID, hash string) error {
	return s.repo.UpdatePasswordHash(ctx, userID, hash)
}
