package banking

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type VBankClient struct {
	baseURL      string
	clientID     string
	clientSecret string
	httpClient   *http.Client
	
	// Token caching
	bankToken   string
	tokenExpiry time.Time
}

func NewVBankClient(baseURL, clientID, clientSecret string) *VBankClient {
	return &VBankClient{
		baseURL:      baseURL,
		clientID:     clientID,
		clientSecret: clientSecret,
		httpClient:   &http.Client{Timeout: 30 * time.Second},
	}
}

// GetRandomDemoClient получает случайного тестового клиента
func (c *VBankClient) GetRandomDemoClient() (*RandomClientResponse, error) {
	url := fmt.Sprintf("%s/auth/random-demo-client", c.baseURL)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var client RandomClientResponse
	if err := json.NewDecoder(resp.Body).Decode(&client); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &client, nil
}

// GetBankToken получает токен банка для межбанковых запросов
func (c *VBankClient) GetBankToken() (string, error) {
	// Проверяем кэш токена
	if c.bankToken != "" && time.Now().Before(c.tokenExpiry) {
		return c.bankToken, nil
	}

	url := fmt.Sprintf("%s/auth/bank-token?client_id=%s&client_secret=%s", 
		c.baseURL, c.clientID, c.clientSecret)
	
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp BankTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// Кэшируем токен
	c.bankToken = tokenResp.AccessToken
	c.tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	return c.bankToken, nil
}

// LoginClient авторизует клиента в VBank
func (c *VBankClient) LoginClient(username, password string) (*ClientLoginResponse, error) {
	url := fmt.Sprintf("%s/auth/login", c.baseURL)
	
	loginReq := ClientLoginRequest{
		Username: username,
		Password: password,
	}

	jsonData, err := json.Marshal(loginReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var loginResp ClientLoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &loginResp, nil
}

// GetAccounts получает список счетов клиента
func (c *VBankClient) GetAccounts(clientToken string) (*AccountsResponse, error) {
	url := fmt.Sprintf("%s/accounts", c.baseURL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+clientToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var accountsResp AccountsResponse
	if err := json.NewDecoder(resp.Body).Decode(&accountsResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &accountsResp, nil
}

// GetTransactions получает транзакции по счету
func (c *VBankClient) GetTransactions(clientToken, accountID string, page, limit int) (*TransactionsResponse, error) {
	url := fmt.Sprintf("%s/accounts/%s/transactions", c.baseURL, accountID)
	if page > 0 {
		url += fmt.Sprintf("?page=%d", page)
	}
	if limit > 0 {
		if page > 0 {
			url += fmt.Sprintf("&limit=%d", limit)
		} else {
			url += fmt.Sprintf("?limit=%d", limit)
		}
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+clientToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var transactionsResp TransactionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&transactionsResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &transactionsResp, nil
}

// GetBalances получает баланс счета
func (c *VBankClient) GetBalances(clientToken, accountID string) (*BalancesResponse, error) {
	url := fmt.Sprintf("%s/accounts/%s/balances", c.baseURL, accountID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+clientToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var balancesResp BalancesResponse
	if err := json.NewDecoder(resp.Body).Decode(&balancesResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &balancesResp, nil
}

// CreateConsent создает согласие для межбанковых запросов
func (c *VBankClient) CreateConsent(clientID, requestingBank, requestingBankName string) (*ConsentResponse, error) {
	bankToken, err := c.GetBankToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get bank token: %w", err)
	}

	url := fmt.Sprintf("%s/account-consents/request", c.baseURL)
	
	consentReq := ConsentRequest{
		ClientID: clientID,
		Permissions: []string{
			"ReadAccountsDetail",
			"ReadBalances",
			"ReadTransactionsDetail",
		},
		Reason:             "Агрегация банковских данных для консорциального кредитования",
		RequestingBank:     requestingBank,
		RequestingBankName: requestingBankName,
	}

	jsonData, err := json.Marshal(consentReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+bankToken)
	req.Header.Set("X-Requesting-Bank", requestingBank)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	var consentResp ConsentResponse
	if err := json.NewDecoder(resp.Body).Decode(&consentResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &consentResp, nil
}
