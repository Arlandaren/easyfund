package banking

// BankingClient интерфейс для клиента банка
type BankingClient interface {
	GetRandomDemoClient() (*RandomClientResponse, error)
	GetBankToken() (string, error)
	LoginClient(username, password string) (*ClientLoginResponse, error)
	GetAccounts(clientToken string) (*AccountsResponse, error)
	GetTransactions(clientToken, accountID string, page, limit int) (*TransactionsResponse, error)
	GetBalances(clientToken, accountID string) (*BalancesResponse, error)
	CreateConsent(clientID, requestingBank, requestingBankName string) (*ConsentResponse, error)
}
