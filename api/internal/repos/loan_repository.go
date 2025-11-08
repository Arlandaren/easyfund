package repos

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/Arlandaren/easyfund/internal/models"
)

type LoanRepository interface {
	CreateLoan(ctx context.Context, loan *models.Loan) (int64, error)
	GetLoanByID(ctx context.Context, loanID int64) (*models.Loan, error)
	ListUserLoans(ctx context.Context, userID uuid.UUID) ([]models.Loan, error)

	CreateLoanSplit(ctx context.Context, split *models.LoanSplit) error
	GetLoanSplits(ctx context.Context, loanID int64) ([]models.LoanSplit, error)
	UpdateLoanSplitPrincipal(ctx context.Context, splitID int64, remainingPrincipal string) error
	UpdateLoanStatus(ctx context.Context, loanID int64, status string) error
}

type loanRepositoryImpl struct {
	db *sql.DB
}

func NewLoanRepository(db *sql.DB) LoanRepository {
	return &loanRepositoryImpl{db: db}
}

func (r *loanRepositoryImpl) CreateLoan(ctx context.Context, loan *models.Loan) (int64, error) {
	const q = `
		INSERT INTO loans (user_id, original_amount, taken_at, interest_rate, status, purpose, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING loan_id
	`
	var id int64
	err := r.db.QueryRowContext(ctx, q,
		loan.UserID, loan.OriginalAmount, loan.TakenAt, loan.InterestRate, loan.Status, loan.Purpose, loan.CreatedAt,
	).Scan(&id)
	return id, err
}

func (r *loanRepositoryImpl) GetLoanByID(ctx context.Context, loanID int64) (*models.Loan, error) {
	const q = `
		SELECT loan_id, user_id, original_amount, taken_at, interest_rate, status, purpose, created_at
		FROM loans WHERE loan_id = $1
	`
	l := &models.Loan{}
	err := r.db.QueryRowContext(ctx, q, loanID).Scan(
		&l.LoanID, &l.UserID, &l.OriginalAmount, &l.TakenAt, &l.InterestRate, &l.Status, &l.Purpose, &l.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func (r *loanRepositoryImpl) ListUserLoans(ctx context.Context, userID uuid.UUID) ([]models.Loan, error) {
	const q = `
		SELECT loan_id, user_id, original_amount, taken_at, interest_rate, status, purpose, created_at
		FROM loans WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var loans []models.Loan
	for rows.Next() {
		var l models.Loan
		if err := rows.Scan(&l.LoanID, &l.UserID, &l.OriginalAmount, &l.TakenAt, &l.InterestRate, &l.Status, &l.Purpose, &l.CreatedAt); err != nil {
			return nil, err
		}
		loans = append(loans, l)
	}
	return loans, rows.Err()
}

func (r *loanRepositoryImpl) CreateLoanSplit(ctx context.Context, split *models.LoanSplit) error {
	const q = `
		INSERT INTO loan_splits (loan_id, bank_id, split_amount, remaining_principal)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.ExecContext(ctx, q, split.LoanID, split.BankID, split.SplitAmount, split.RemainingPrincipal)
	return err
}

func (r *loanRepositoryImpl) GetLoanSplits(ctx context.Context, loanID int64) ([]models.LoanSplit, error) {
	const q = `
		SELECT split_id, loan_id, bank_id, split_amount, remaining_principal
		FROM loan_splits WHERE loan_id = $1
		ORDER BY split_id
	`
	rows, err := r.db.QueryContext(ctx, q, loanID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var splits []models.LoanSplit
	for rows.Next() {
		var s models.LoanSplit
		if err := rows.Scan(&s.SplitID, &s.LoanID, &s.BankID, &s.SplitAmount, &s.RemainingPrincipal); err != nil {
			return nil, err
		}
		splits = append(splits, s)
	}
	return splits, rows.Err()
}

func (r *loanRepositoryImpl) UpdateLoanSplitPrincipal(ctx context.Context, splitID int64, remainingPrincipal string) error {
	const q = `UPDATE loan_splits SET remaining_principal = $1 WHERE split_id = $2`
	_, err := r.db.ExecContext(ctx, q, remainingPrincipal, splitID)
	return err
}

func (r *loanRepositoryImpl) UpdateLoanStatus(ctx context.Context, loanID int64, status string) error {
	const q = `UPDATE loans SET status = $1 WHERE loan_id = $2`
	_, err := r.db.ExecContext(ctx, q, status, loanID)
	return err
}
