package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/Arlandaren/easyfund/internal/models"
	"github.com/Arlandaren/easyfund/internal/repos"
)

type LoanService interface {
	CreateLoan(ctx context.Context, loan *models.Loan, splits []map[int16]string) (*models.LoanDetailDTO, error)
	GetLoanDetail(ctx context.Context, loanID int64) (*models.LoanDetailDTO, error)
	ListUserLoans(ctx context.Context, userID uuid.UUID) ([]models.Loan, error)
	MakePayment(ctx context.Context, payment *models.LoanPayment, allocations []models.PaymentAllocation) error
	GetTotalDebt(ctx context.Context, userID uuid.UUID) (string, error)
}

type loanServiceImpl struct {
	loanRepo    repos.LoanRepository
	accountRepo repos.UserBankAccountRepository
	paymentRepo repos.LoanPaymentRepository
}

func NewLoanService(loanRepo repos.LoanRepository, accountRepo repos.UserBankAccountRepository, paymentRepo repos.LoanPaymentRepository) LoanService {
	return &loanServiceImpl{
		loanRepo:    loanRepo,
		accountRepo: accountRepo,
		paymentRepo: paymentRepo,
	}
}

func (s *loanServiceImpl) CreateLoan(ctx context.Context, loan *models.Loan, splits []map[int16]string) (*models.LoanDetailDTO, error) {
	// Создаём кредит
	loanID, err := s.loanRepo.CreateLoan(ctx, loan)
	if err != nil {
		return nil, fmt.Errorf("failed to create loan: %w", err)
	}

	loan.LoanID = loanID

	// Создаём splits для каждого банка
	var loanSplits []models.LoanSplit
	for _, splitData := range splits {
		for bankID, amount := range splitData {
			split := models.LoanSplit{
				LoanID:             loanID,
				BankID:             bankID,
				SplitAmount:        amount,
				RemainingPrincipal: amount,
			}
			err := s.loanRepo.CreateLoanSplit(ctx, &split)
			if err != nil {
				return nil, fmt.Errorf("failed to create loan split: %w", err)
			}
			loanSplits = append(loanSplits, split)
		}
	}

	return &models.LoanDetailDTO{
		Loan:           loan,
		Splits:         loanSplits,
		PercentPaid:    0.0,
		RemainingDebt:  loan.OriginalAmount,
		PaymentHistory: []models.LoanPayment{},
	}, nil
}

func (s *loanServiceImpl) GetLoanDetail(ctx context.Context, loanID int64) (*models.LoanDetailDTO, error) {
	loan, err := s.loanRepo.GetLoanByID(ctx, loanID)
	if err != nil {
		return nil, fmt.Errorf("failed to get loan: %w", err)
	}

	splits, err := s.loanRepo.GetLoanSplits(ctx, loanID)
	if err != nil {
		return nil, fmt.Errorf("failed to get loan splits: %w", err)
	}

	payments, err := s.paymentRepo.ListLoanPayments(ctx, loanID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payments: %w", err)
	}

	// Вычисляем процент выплаченного и остаток долга
	// TODO: добавить правильный расчет на основе principal/interest
	percentPaid := 0.0
	remainingDebt := loan.OriginalAmount

	return &models.LoanDetailDTO{
		Loan:           loan,
		Splits:         splits,
		PercentPaid:    percentPaid,
		RemainingDebt:  remainingDebt,
		PaymentHistory: payments,
	}, nil
}

func (s *loanServiceImpl) ListUserLoans(ctx context.Context, userID uuid.UUID) ([]models.Loan, error) {
	return s.loanRepo.ListUserLoans(ctx, userID)
}

func (s *loanServiceImpl) MakePayment(ctx context.Context, payment *models.LoanPayment, allocations []models.PaymentAllocation) error {
	// Создаём платеж
	paymentID, err := s.paymentRepo.CreatePayment(ctx, payment)
	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}

	// Создаём allocations и обновляем remaining_principal
	for _, alloc := range allocations {
		alloc.PaymentID = paymentID
		err := s.paymentRepo.CreatePaymentAllocation(ctx, &alloc)
		if err != nil {
			return fmt.Errorf("failed to create payment allocation: %w", err)
		}

		// Обновляем remaining_principal в loan_split
		// TODO: вычислить новый remaining_principal
		// err = s.loanRepo.UpdateLoanSplitPrincipal(ctx, alloc.SplitID, newRemaining)
	}

	return nil
}

func (s *loanServiceImpl) GetTotalDebt(ctx context.Context, userID uuid.UUID) (string, error) {
	// SELECT SUM(remaining_principal) FROM loan_splits WHERE loan_id IN (SELECT loan_id FROM loans WHERE user_id = $1 AND status = 'ACTIVE')
	// TODO: реализовать правильный запрос
	return "0.00", nil
}
