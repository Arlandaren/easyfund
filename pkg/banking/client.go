package banking

// VBankAPI — интерфейс клиента Virtual Bank, который используют usecase/handlers.
type VBankAPI interface {
    // Token
    GetBankToken() (string, error)

    // Demo/login
    GetRandomDemoClient() (*RandomClientResponse, error)
    LoginClient(username, password string) (*ClientLoginResponse, error)

    // Consent
    CreateConsent(clientID, requestingBank, requestingBankName string) (*ConsentResponse, error)

    // Data (через consent)
    GetAccountsWithConsent(clientID, consentID string) (*AccountsResponse, error)
    GetTransactionsWithConsent(accountID, clientID, consentID string, page, limit int) (*TransactionsResponse, error)
    GetBalancesWithConsent(accountID, clientID, consentID string) (*BalancesResponse, error)
}

// Гарантируем на этапе компиляции соответствие интерфейсу.
var _ VBankAPI = (*VBankClient)(nil)
var _ VBankAPI = (*MockClient)(nil)
