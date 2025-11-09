package services

import (
	"context"

	"github.com/Arlandaren/easyfund/internal/models"
	"github.com/Arlandaren/easyfund/internal/repos"
)

type UserBankAccountService interface {
	GetTotalBalance(ctx context.Context, userID int64) (string, error)
	GetUserAccounts(ctx context.Context, userID int64) ([]models.UserBankAccount, error)
	CreateAccount(ctx context.Context, account *models.UserBankAccount) error
}

type userBankAccountServiceImpl struct {
	repo repos.UserBankAccountRepository
}

func NewUserBankAccountService(repo repos.UserBankAccountRepository) UserBankAccountService {
	return &userBankAccountServiceImpl{repo: repo}
}

func (s *userBankAccountServiceImpl) GetTotalBalance(ctx context.Context, userID int64) (string, error) {
	return s.repo.GetTotalBalance(ctx, userID)
}

func (s *userBankAccountServiceImpl) GetUserAccounts(ctx context.Context, userID int64) ([]models.UserBankAccount, error) {
	return s.repo.GetUserAccounts(ctx, userID)
}

func (s *userBankAccountServiceImpl) CreateAccount(ctx context.Context, account *models.UserBankAccount) error {
	return s.repo.CreateAccount(ctx, account)
}