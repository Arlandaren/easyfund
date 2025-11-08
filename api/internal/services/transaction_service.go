package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/Arlandaren/easyfund/internal/models"
	"github.com/Arlandaren/easyfund/internal/repos"
)

type TransactionService interface {
	GetUserTransactionHistory(ctx context.Context, userID uuid.UUID) (*models.TransactionHistoryDTO, error)
	GetBankTransactionHistory(ctx context.Context, userID uuid.UUID, bankID int16) ([]models.Transaction, error)
	CreateTransaction(ctx context.Context, tx *models.Transaction) error
}

type transactionServiceImpl struct {
	repo repos.TransactionRepository
}

func NewTransactionService(repo repos.TransactionRepository) TransactionService {
	return &transactionServiceImpl{repo: repo}
}

func (s *transactionServiceImpl) GetUserTransactionHistory(ctx context.Context, userID uuid.UUID) (*models.TransactionHistoryDTO, error) {
	transactions, err := s.repo.ListUserTransactions(ctx, userID)
	if err != nil {
		return nil, err
	}

	totalSpent, err := s.repo.GetUserTotalSpent(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &models.TransactionHistoryDTO{
		Transactions: transactions,
		TotalSpent:   totalSpent,
	}, nil
}

func (s *transactionServiceImpl) GetBankTransactionHistory(ctx context.Context, userID uuid.UUID, bankID int16) ([]models.Transaction, error) {
	return s.repo.ListBankTransactions(ctx, userID, bankID)
}

func (s *transactionServiceImpl) CreateTransaction(ctx context.Context, tx *models.Transaction) error {
	return s.repo.CreateTransaction(ctx, tx)
}
