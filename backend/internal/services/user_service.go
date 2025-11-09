package services

import (
	"context"

	"github.com/Arlandaren/easyfund/internal/models"
	"github.com/Arlandaren/easyfund/internal/repos"
)

type UserService interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, userID int64) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetRandomUser(ctx context.Context) (*models.User, error)
	ListUsers(ctx context.Context, limit int) ([]models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, userID int64) error
	UpdatePasswordHash(ctx context.Context, userID int64, hash string) error
}

type userServiceImpl struct {
	repo repos.UserRepository
}

func NewUserService(repo repos.UserRepository) UserService {
	return &userServiceImpl{repo: repo}
}

func (s *userServiceImpl) CreateUser(ctx context.Context, user *models.User) error {
	return s.repo.CreateUser(ctx, user)
}

func (s *userServiceImpl) GetUserByID(ctx context.Context, userID int64) (*models.User, error) {
	return s.repo.GetUserByID(ctx, userID)
}

func (s *userServiceImpl) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return s.repo.GetUserByEmail(ctx, email)
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

func (s *userServiceImpl) DeleteUser(ctx context.Context, userID int64) error {
	return s.repo.DeleteUser(ctx, userID)
}

func (s *userServiceImpl) UpdatePasswordHash(ctx context.Context, userID int64, hash string) error {
	return s.repo.UpdatePasswordHash(ctx, userID, hash)
}