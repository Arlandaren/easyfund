package repos

import (
	"context"
	"database/sql"

	"github.com/Arlandaren/easyfund/internal/models"
)

type UserBankAccountRepository interface {
	CreateAccount(ctx context.Context, account *models.UserBankAccount) error
	GetAccountByID(ctx context.Context, accountID int64) (*models.UserBankAccount, error)
	GetUserAccounts(ctx context.Context, userID int64) ([]models.UserBankAccount, error)
	UpdateBalance(ctx context.Context, accountID int64, balance string) error
	GetTotalBalance(ctx context.Context, userID int64) (string, error)
}

type userBankAccountRepositoryImpl struct {
	db *sql.DB
}

func NewUserBankAccountRepository(db *sql.DB) UserBankAccountRepository {
	return &userBankAccountRepositoryImpl{db: db}
}

func (r *userBankAccountRepositoryImpl) CreateAccount(ctx context.Context, account *models.UserBankAccount) error {
	const q = `
		INSERT INTO user_bank_accounts (user_id, bank_id, balance, currency, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, q, account.UserID, account.BankID, account.Balance, account.Currency, account.CreatedAt)
	return err
}

func (r *userBankAccountRepositoryImpl) GetAccountByID(ctx context.Context, accountID int64) (*models.UserBankAccount, error) {
	const q = `
		SELECT account_id, user_id, bank_id, balance, currency, created_at
		FROM user_bank_accounts WHERE account_id = $1
	`
	acc := &models.UserBankAccount{}
	err := r.db.QueryRowContext(ctx, q, accountID).Scan(
		&acc.AccountID, &acc.UserID, &acc.BankID, &acc.Balance, &acc.Currency, &acc.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return acc, nil
}

func (r *userBankAccountRepositoryImpl) GetUserAccounts(ctx context.Context, userID int64) ([]models.UserBankAccount, error) {
	const q = `
		SELECT account_id, user_id, bank_id, balance, currency, created_at
		FROM user_bank_accounts WHERE user_id = $1
		ORDER BY bank_id
	`
	rows, err := r.db.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []models.UserBankAccount
	for rows.Next() {
		var a models.UserBankAccount
		if err := rows.Scan(&a.AccountID, &a.UserID, &a.BankID, &a.Balance, &a.Currency, &a.CreatedAt); err != nil {
			return nil, err
		}
		accounts = append(accounts, a)
	}
	return accounts, rows.Err()
}

func (r *userBankAccountRepositoryImpl) UpdateBalance(ctx context.Context, accountID int64, balance string) error {
	const q = `UPDATE user_bank_accounts SET balance = $1 WHERE account_id = $2`
	_, err := r.db.ExecContext(ctx, q, balance, accountID)
	return err
}

func (r *userBankAccountRepositoryImpl) GetTotalBalance(ctx context.Context, userID int64) (string, error) {
	const q = `SELECT COALESCE(TO_CHAR(SUM(balance), 'FM9999999999990.00'), '0.00') FROM user_bank_accounts WHERE user_id = $1`
	var total string
	err := r.db.QueryRowContext(ctx, q, userID).Scan(&total)
	return total, err
}