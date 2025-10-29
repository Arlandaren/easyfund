package usecase

import (
    "github.com/Arlandaren/easyfund/pkg/banking"
    "fmt"
    "strconv"
)

type BankingUsecase struct {
    vbankClient *banking.VBankClient
}

// Создаем usecase с клиентом VBank
func NewBankingUsecase(vbankClient *banking.VBankClient) *BankingUsecase {
    return &BankingUsecase{
        vbankClient: vbankClient,
    }
}

// Теперь токен берется из claims.VBankAccessToken, пароль не нужен
func (b *BankingUsecase) GetAccounts(claims *banking.JWTClaims) (*banking.AccountsResponse, error) {
    if claims.Mode != "demo" || claims.PersonID == "" || claims.VBankAccessToken == "" {
        return nil, fmt.Errorf("only demo clients with access token can access VBank data")
    }

    return b.vbankClient.GetAccounts(claims.VBankAccessToken)
}

func (b *BankingUsecase) GetTransactions(claims *banking.JWTClaims, accountID string, page, limit int) (*banking.TransactionsResponse, error) {
    if claims.Mode != "demo" || claims.PersonID == "" || claims.VBankAccessToken == "" {
        return nil, fmt.Errorf("only demo clients with access token can access VBank data")
    }

    return b.vbankClient.GetTransactions(claims.VBankAccessToken, accountID, page, limit)
}

func (b *BankingUsecase) GetBalances(claims *banking.JWTClaims, accountID string) (*banking.BalancesResponse, error) {
    if claims.Mode != "demo" || claims.PersonID == "" || claims.VBankAccessToken == "" {
        return nil, fmt.Errorf("only demo clients with access token can access VBank data")
    }

    return b.vbankClient.GetBalances(claims.VBankAccessToken, accountID)
}

func (b *BankingUsecase) GetFinancialInsights(claims *banking.JWTClaims) (map[string]interface{}, error) {
    if claims.Mode != "demo" || claims.PersonID == "" || claims.VBankAccessToken == "" {
        return nil, fmt.Errorf("only demo clients with access token can access VBank data")
    }

    accounts, err := b.GetAccounts(claims)
    if err != nil {
        return nil, fmt.Errorf("failed to get accounts: %w", err)
    }

    insights := map[string]interface{}{
        "total_accounts": len(accounts.Accounts),
        "currencies":     []string{},
        "account_types":  []string{},
        "total_balance":  0.0,
    }

    currencyMap := make(map[string]bool)
    typeMap := make(map[string]bool)
    totalBalance := 0.0

    for _, account := range accounts.Accounts {
        if !currencyMap[account.Currency] {
            currencyMap[account.Currency] = true
            insights["currencies"] = append(insights["currencies"].([]string), account.Currency)
        }

        if !typeMap[account.AccountType] {
            typeMap[account.AccountType] = true
            insights["account_types"] = append(insights["account_types"].([]string), account.AccountType)
        }

        balances, err := b.GetBalances(claims, account.AccountID)
        if err == nil && len(balances.Data) > 0 {
            for _, balance := range balances.Data {
                if balance.Currency == "RUB" {
                    if amount, err := strconv.ParseFloat(balance.Amount, 64); err == nil {
                        totalBalance += amount
                    }
                }
            }
        }
    }

    insights["total_balance"] = totalBalance

    return insights, nil
}
