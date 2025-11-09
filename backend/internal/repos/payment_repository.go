package repos

import (
	"context"
	"database/sql"

	"github.com/Arlandaren/easyfund/internal/models"
)

type LoanPaymentRepository interface {
	CreatePayment(ctx context.Context, payment *models.LoanPayment) (int64, error)
	GetPaymentByID(ctx context.Context, paymentID int64) (*models.LoanPayment, error)
	ListLoanPayments(ctx context.Context, loanID int64) ([]models.LoanPayment, error)
	ListUserPayments(ctx context.Context, userID int64) ([]models.LoanPayment, error)

	CreatePaymentAllocation(ctx context.Context, alloc *models.PaymentAllocation) error
	GetPaymentAllocations(ctx context.Context, paymentID int64) ([]models.PaymentAllocation, error)
}

type loanPaymentRepositoryImpl struct {
	db *sql.DB
}

func NewLoanPaymentRepository(db *sql.DB) LoanPaymentRepository {
	return &loanPaymentRepositoryImpl{db: db}
}

func (r *loanPaymentRepositoryImpl) CreatePayment(ctx context.Context, p *models.LoanPayment) (int64, error) {
	const q = `
		INSERT INTO loan_payments (loan_id, user_id, paid_at, total_amount, comment)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING payment_id
	`
	var id int64
	err := r.db.QueryRowContext(ctx, q, p.LoanID, p.UserID, p.PaidAt, p.TotalAmount, p.Comment).Scan(&id)
	return id, err
}

func (r *loanPaymentRepositoryImpl) GetPaymentByID(ctx context.Context, id int64) (*models.LoanPayment, error) {
	const q = `
		SELECT payment_id, loan_id, user_id, paid_at, total_amount, comment
		FROM loan_payments WHERE payment_id = $1
	`
	p := &models.LoanPayment{}
	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&p.PaymentID, &p.LoanID, &p.UserID, &p.PaidAt, &p.TotalAmount, &p.Comment,
	)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *loanPaymentRepositoryImpl) ListLoanPayments(ctx context.Context, loanID int64) ([]models.LoanPayment, error) {
	const q = `
		SELECT payment_id, loan_id, user_id, paid_at, total_amount, comment
		FROM loan_payments WHERE loan_id = $1
		ORDER BY paid_at DESC
	`
	rows, err := r.db.QueryContext(ctx, q, loanID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []models.LoanPayment
	for rows.Next() {
		var p models.LoanPayment
		if err := rows.Scan(&p.PaymentID, &p.LoanID, &p.UserID, &p.PaidAt, &p.TotalAmount, &p.Comment); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, rows.Err()
}

func (r *loanPaymentRepositoryImpl) ListUserPayments(ctx context.Context, userID int64) ([]models.LoanPayment, error) {
	const q = `
		SELECT payment_id, loan_id, user_id, paid_at, total_amount, comment
		FROM loan_payments WHERE user_id = $1
		ORDER BY paid_at DESC
	`
	rows, err := r.db.QueryContext(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []models.LoanPayment
	for rows.Next() {
		var p models.LoanPayment
		if err := rows.Scan(&p.PaymentID, &p.LoanID, &p.UserID, &p.PaidAt, &p.TotalAmount, &p.Comment); err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	return res, rows.Err()
}

func (r *loanPaymentRepositoryImpl) CreatePaymentAllocation(ctx context.Context, a *models.PaymentAllocation) error {
	const q = `
		INSERT INTO payment_allocations (payment_id, split_id, principal_paid, interest_paid)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.ExecContext(ctx, q, a.PaymentID, a.SplitID, a.PrincipalPaid, a.InterestPaid)
	return err
}

func (r *loanPaymentRepositoryImpl) GetPaymentAllocations(ctx context.Context, paymentID int64) ([]models.PaymentAllocation, error) {
	const q = `
		SELECT allocation_id, payment_id, split_id, principal_paid, interest_paid
		FROM payment_allocations WHERE payment_id = $1
		ORDER BY allocation_id
	`
	rows, err := r.db.QueryContext(ctx, q, paymentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var res []models.PaymentAllocation
	for rows.Next() {
		var a models.PaymentAllocation
		if err := rows.Scan(&a.AllocationID, &a.PaymentID, &a.SplitID, &a.PrincipalPaid, &a.InterestPaid); err != nil {
			return nil, err
		}
		res = append(res, a)
	}
	return res, rows.Err()
}