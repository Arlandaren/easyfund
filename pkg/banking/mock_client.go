package banking

import (
	"errors"
	"time"
)

// Мок-реализация интерфейса BankingClient для тестов
type MockClient struct{}

// GetAccountsWithConsent implements VBankAPI.
func (m *MockClient) GetAccountsWithConsent(clientID string, consentID string) (*AccountsResponse, error) {
	panic("unimplemented")
}

// GetBalancesWithConsent implements VBankAPI.
func (m *MockClient) GetBalancesWithConsent(accountID string, clientID string, consentID string) (*BalancesResponse, error) {
	panic("unimplemented")
}

// GetTransactionsWithConsent implements VBankAPI.
func (m *MockClient) GetTransactionsWithConsent(accountID string, clientID string, consentID string, page int, limit int) (*TransactionsResponse, error) {
	panic("unimplemented")
}

func NewMockClient() *MockClient {
	return &MockClient{}
}

func (m *MockClient) GetRandomDemoClient() (*RandomClientResponse, error) {
	return &RandomClientResponse{
		PersonID: "mock-pid",
		FullName: "Mock User",
		Password: "password",
	}, nil
}

func (m *MockClient) GetBankToken() (string, error) {
	return "mock-bank-token", nil
}

func (m *MockClient) LoginClient(username, password string) (*ClientLoginResponse, error) {
	if username == "mock-pid" && password == "password" {
		return &ClientLoginResponse{
			AccessToken: "mock-access-token",
			TokenType:   "bearer",
			ClientID:    "mock-client",
		}, nil
	}
	return nil, errors.New("invalid credentials in mock")
}

func (m *MockClient) GetAccounts(clientToken string) (*AccountsResponse, error) {
	return &AccountsResponse{
		Accounts: []Account{
			{
				AccountID:   "mock-acc-1",
				Currency:    "RUB",
				AccountType: "Personal",
				Nickname:    "Mock Account",
			},
		},
	}, nil
}

func (m *MockClient) GetTransactions(clientToken, accountID string, page, limit int) (*TransactionsResponse, error) {
	return &TransactionsResponse{
		Data: []Transaction{
			{
				TransactionID:        "mock-tx-1",
				Amount:               "1000.00",
				Currency:             "RUB",
				CreditDebitIndicator: "credit",
				Status:               "completed",
				BookingDateTime:      time.Now(),
			},
		},
	}, nil
}

func (m *MockClient) GetBalances(clientToken, accountID string) (*BalancesResponse, error) {
	return &BalancesResponse{
		Data: []Balance{
			{
				Amount:   "5000.00",
				Currency: "RUB",
			},
		},
	}, nil
}

func (m *MockClient) CreateConsent(clientID, requestingBank, requestingBankName string) (*ConsentResponse, error) {
	return &ConsentResponse{
		Status:       "approved",
		ConsentID:    "mock-consent-1",
		AutoApproved: true,
	}, nil
}
