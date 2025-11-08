package repos

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/Arlandaren/easyfund/internal/models"
)

type CreditApplicationRepository interface {
	CreateApplication(ctx context.Context, app *models.CreditApplication) (int64, error)
	GetApplicationByID(ctx context.Context, appID int64) (*models.CreditApplication, error)
	ListUserApplications(ctx context.Context, userID uuid.UUID) ([]models.CreditApplication, error)
	UpdateApplicationStatus(ctx context.Context, appID int64, status string) error
	UpdateApplicationLoanID(ctx context.Context, appID int64, loanID int64) error
}

type creditApplicationRepositoryImpl struct {
	db *sql.DB
}

func NewCreditApplicationRepository(db *sql.DB) CreditApplicationRepository {
	return &creditApplicationRepositoryImpl{db: db}
}

func (r *creditApplicationRepositoryImpl) CreateApplication(ctx context.Context, app *models.CreditApplication) (int64, error) {
	const q = `
		INSERT INTO credit_applications (user_id, bank_id, type_code, status_code, requested_amount, loan_id, submitted_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING application_id
	`
	var id int64
	err := r.db.QueryRowContext(ctx, q,
		app.UserID, app.BankID, app.TypeCode, app.StatusCode, app.RequestedAmount, app.LoanID, app.SubmittedAt, app.UpdatedAt,
	).Scan(&id)
	return id, err
}

func (r *creditApplicationRepositoryImpl) GetApplicationByID(ctx context.Context, appID int64) (*models.CreditApplication, error) {
	const q = `
		SELECT application_id, user_id, bank_id, type_code, status_code, requested_amount, loan_id, submitted_at, updated_at
		FROM credit_applications WHERE application_id = $1
	`
	a := &models.CreditApplication{}
	err := r.db.QueryRowContext(ctx, q, appID).Scan(
		&a.ApplicationID, &a.UserID, &a.BankID, &a.TypeCode, &a.StatusCode, &a.RequestedAmount, &a.LoanID, &a.SubmittedAt, &a.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (r *creditApplicationRepositoryImpl) ListUserApplications(ctx context.Context, userID uuid.UUID) ([]models.CreditApplication, error) {
	const q = `
		SELECT application_id, user_id, bank_id, type_code, status_code, requested_amount, loan_id, submitted_at, updated_at
		FROM credit_applications WHERE user_id = $1
		ORDER BY submitted_at DESC
	`
	rows, err := r.db.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []models.CreditApplication
	for rows.Next() {
		var a models.CreditApplication
		if err := rows.Scan(
			&a.ApplicationID, &a.UserID, &a.BankID, &a.TypeCode, &a.StatusCode, &a.RequestedAmount, &a.LoanID, &a.SubmittedAt, &a.UpdatedAt,
		); err != nil {
			return nil, err
		}
		res = append(res, a)
	}
	return res, rows.Err()
}

func (r *creditApplicationRepositoryImpl) UpdateApplicationStatus(ctx context.Context, appID int64, status string) error {
	const q = `UPDATE credit_applications SET status_code = $1, updated_at = NOW() WHERE application_id = $2`
	_, err := r.db.ExecContext(ctx, q, status, appID)
	return err
}

func (r *creditApplicationRepositoryImpl) UpdateApplicationLoanID(ctx context.Context, appID int64, loanID int64) error {
	const q = `UPDATE credit_applications SET loan_id = $1, updated_at = NOW() WHERE application_id = $2`
	_, err := r.db.ExecContext(ctx, q, loanID, appID)
	return err
}
