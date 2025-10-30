package usecase

import (
    "context"

    "github.com/Arlandaren/easyfund/pkg/banking"
)

type BankingUsecase struct {
    client banking.VBankAPI
    repo   ConsentRepository
}

type ConsentRepository interface {
    GetConsentID(ctx context.Context, clientID string) (string, error)
    SaveConsentID(ctx context.Context, clientID, consentID string) error
}

func NewBankingUsecase(client banking.VBankAPI) *BankingUsecase {
    return &BankingUsecase{client: client}
}

func (u *BankingUsecase) GetConsentID(ctx context.Context, clientID string) (string, error) {
    if u.repo == nil {
        return "", nil
    }
    return u.repo.GetConsentID(ctx, clientID)
}

func (u *BankingUsecase) SaveConsentID(ctx context.Context, clientID, consentID string) error {
    if u.repo == nil {
        return nil
    }
    return u.repo.SaveConsentID(ctx, clientID, consentID)
}

func (u *BankingUsecase) GetFinancialInsights(claims *banking.JWTClaims) (interface{}, error) {
    // оставьте вашу логику
    return map[string]any{"status": "ok"}, nil
}
