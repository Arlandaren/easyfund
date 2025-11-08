package repos

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/Arlandaren/easyfund/internal/models"
)

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, tx *models.Transaction) error
	GetTransactionByID(ctx context.Context, id int64) (*models.Transaction, error)
	ListUserTransactions(ctx context.Context, userID uuid.UUID) ([]models.Transaction, error)
	ListBankTransactions(ctx context.Context, userID uuid.UUID, bankID int16) ([]models.Transaction, error)
	GetUserTotalSpent(ctx context.Context, userID uuid.UUID) (string, error)
}

type transactionRepositoryImpl struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) TransactionRepository {
	return &transactionRepositoryImpl{db: db}
}

func (r *transactionRepositoryImpl) CreateTransaction(ctx context.Context, t *models.Transaction) error {
	const q = `
		INSERT INTO transactions (user_id, bank_id, occurred_at, amount, category, description)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, q, t.UserID, t.BankID, t.OccurredAt, t.Amount, t.Category, t.Description)
	return err
}

func (r *transactionRepositoryImpl) GetTransactionByID(ctx context.Context, id int64) (*models.Transaction, error) {
	const q = `
		SELECT transaction_id, user_id, bank_id, occurred_at, amount, category, description
		FROM transactions WHERE transaction_id = $1
	`
	t := &models.Transaction{}
	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&t.TransactionID, &t.UserID, &t.BankID, &t.OccurredAt, &t.Amount, &t.Category, &t.Description,
	)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (r *transactionRepositoryImpl) ListUserTransactions(ctx context.Context, userID uuid.UUID) ([]models.Transaction, error) {
	const q = `
		SELECT transaction_id, user_id, bank_id, occurred_at, amount, category, description
		FROM transactions WHERE user_id = $1
		ORDER BY occurred_at DESC
	`
	rows, err := r.db.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []models.Transaction
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.TransactionID, &t.UserID, &t.BankID, &t.OccurredAt, &t.Amount, &t.Category, &t.Description); err != nil {
			return nil, err
		}
		res = append(res, t)
	}
	return res, rows.Err()
}

func (r *transactionRepositoryImpl) ListBankTransactions(ctx context.Context, userID uuid.UUID, bankID int16) ([]models.Transaction, error) {
	const q = `
		SELECT transaction_id, user_id, bank_id, occurred_at, amount, category, description
		FROM transactions WHERE user_id = $1 AND bank_id = $2
		ORDER BY occurred_at DESC
	`
	rows, err := r.db.QueryContext(ctx, q, userID, bankID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []models.Transaction
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.TransactionID, &t.UserID, &t.BankID, &t.OccurredAt, &t.Amount, &t.Category, &t.Description); err != nil {
			return nil, err
		}
		res = append(res, t)
	}
	return res, rows.Err()
}

func (r *transactionRepositoryImpl) GetUserTotalSpent(ctx context.Context, userID uuid.UUID) (string, error) {
	const q = `SELECT COALESCE(TO_CHAR(SUM(amount), 'FM9999999999990.00'), '0.00') FROM transactions WHERE user_id = $1`
	var total string
	err := r.db.QueryRowContext(ctx, q, userID).Scan(&total)
	return total, err
}
