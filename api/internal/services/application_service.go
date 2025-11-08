package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/Arlandaren/easyfund/internal/models"
	"github.com/Arlandaren/easyfund/internal/repos"
)

type CreditApplicationService interface {
	SubmitApplication(ctx context.Context, app *models.CreditApplication) (int64, error)
	GetApplications(ctx context.Context, userID uuid.UUID) ([]models.CreditApplication, error)
	ApproveApplication(ctx context.Context, appID int64, splits []map[int16]string) error
	RejectApplication(ctx context.Context, appID int64) error
}

type creditApplicationServiceImpl struct {
	appRepo  repos.CreditApplicationRepository
	loanRepo repos.LoanRepository
}

func NewCreditApplicationService(appRepo repos.CreditApplicationRepository, loanRepo repos.LoanRepository) CreditApplicationService {
	return &creditApplicationServiceImpl{
		appRepo:  appRepo,
		loanRepo: loanRepo,
	}
}

func (s *creditApplicationServiceImpl) SubmitApplication(ctx context.Context, app *models.CreditApplication) (int64, error) {
	return s.appRepo.CreateApplication(ctx, app)
}

func (s *creditApplicationServiceImpl) GetApplications(ctx context.Context, userID uuid.UUID) ([]models.CreditApplication, error) {
	return s.appRepo.ListUserApplications(ctx, userID)
}

func (s *creditApplicationServiceImpl) ApproveApplication(ctx context.Context, appID int64, splits []map[int16]string) error {
	app, err := s.appRepo.GetApplicationByID(ctx, appID)
	if err != nil {
		return fmt.Errorf("failed to get application: %w", err)
	}

	// Создаём кредит на основе заявки
	loan := &models.Loan{
		UserID:        app.UserID,
		OriginalAmount: app.RequestedAmount,
		Status:        "ACTIVE",
		Purpose:       app.TypeCode,
	}

	loanID, err := s.loanRepo.CreateLoan(ctx, loan)
	if err != nil {
		return fmt.Errorf("failed to create loan: %w", err)
	}

	// Создаём splits
	for _, splitData := range splits {
		for bankID, amount := range splitData {
			split := &models.LoanSplit{
				LoanID:             loanID,
				BankID:             bankID,
				SplitAmount:        amount,
				RemainingPrincipal: amount,
			}
			err := s.loanRepo.CreateLoanSplit(ctx, split)
			if err != nil {
				return fmt.Errorf("failed to create loan split: %w", err)
			}
		}
	}

	// Обновляем статус заявки и привязываем кредит
	err = s.appRepo.UpdateApplicationStatus(ctx, appID, "APPROVED")
	if err != nil {
		return fmt.Errorf("failed to update application status: %w", err)
	}

	err = s.appRepo.UpdateApplicationLoanID(ctx, appID, loanID)
	if err != nil {
		return fmt.Errorf("failed to update application loan_id: %w", err)
	}

	return nil
}

func (s *creditApplicationServiceImpl) RejectApplication(ctx context.Context, appID int64) error {
	return s.appRepo.UpdateApplicationStatus(ctx, appID, "REJECTED")
}
