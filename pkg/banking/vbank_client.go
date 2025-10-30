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
    clientID     string // teamXXX (для X-Requesting-Bank)
    clientSecret string
    httpClient   *http.Client

    // Token caching
    bankToken   string
    tokenExpiry time.Time
}

func NewVBankClient(baseURL, clientID, clientSecret string) *VBankClient {
    return &VBankClient{
        baseURL:      trimRightSlash(baseURL),
        clientID:     clientID,
        clientSecret: clientSecret,
        httpClient:   &http.Client{Timeout: 30 * time.Second},
    }
}

func trimRightSlash(s string) string {
    if len(s) > 0 && s[len(s)-1] == '/' {
        return s[:len(s)-1]
    }
    return s
}

// ensureBankToken гарантирует валидный bank token в памяти.
func (c *VBankClient) ensureBankToken() (string, error) {
    if c.bankToken != "" && time.Now().Before(c.tokenExpiry) {
        return c.bankToken, nil
    }
    return c.refreshBankToken()
}

func (c *VBankClient) refreshBankToken() (string, error) {
    url := fmt.Sprintf("%s/auth/bank-token?client_id=%s&client_secret=%s", c.baseURL, c.clientID, c.clientSecret)

    req, err := http.NewRequest("POST", url, nil)
    if err != nil {
        return "", fmt.Errorf("create bank-token request: %w", err)
    }
    req.Header.Set("User-Agent", "easyfund-vbank-client/1.0")
    req.Header.Set("Accept", "application/json")

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return "", fmt.Errorf("do bank-token request: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        b, _ := io.ReadAll(resp.Body)
        return "", fmt.Errorf("bank-token HTTP %d: %s", resp.StatusCode, string(b))
    }

    var tok BankTokenResponse
    if err := json.NewDecoder(resp.Body).Decode(&tok); err != nil {
        return "", fmt.Errorf("decode bank-token: %w", err)
    }

    // Кэшируем
    c.bankToken = tok.AccessToken
    // Небольшой буфер безопасности, чтобы не упасть ровно на границе истечения
    c.tokenExpiry = time.Now().Add(time.Duration(tok.ExpiresIn-60) * time.Second)
    if tok.ExpiresIn < 120 {
        // если тестовый токен короткий — не уходим в отрицательные значения
        c.tokenExpiry = time.Now().Add(time.Duration(tok.ExpiresIn) * time.Second)
    }

    return c.bankToken, nil
}

// Публичный метод получения токена (если нужен)
func (c *VBankClient) GetBankToken() (string, error) {
    return c.ensureBankToken()
}

// Вспомогательный метод отправки HTTP с авто-ретраем по 401/403 (обновление токена)
func (c *VBankClient) doWithBankAuth(req *http.Request) (*http.Response, error) {
    // Первый проход — с текущим/гарантированным токеном
    token, err := c.ensureBankToken()
    if err != nil {
        return nil, err
    }
    req.Header.Set("Authorization", "Bearer "+token)
    req.Header.Set("User-Agent", "easyfund-vbank-client/1.0")
    req.Header.Set("Accept", "application/json")

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }

    if resp.StatusCode != http.StatusUnauthorized && resp.StatusCode != http.StatusForbidden {
        return resp, nil
    }

    // Если токен истёк/невалиден — один раз обновим и повторим
    _ = resp.Body.Close()
    if _, err := c.refreshBankToken(); err != nil {
        return nil, fmt.Errorf("refresh token after %d: %w", resp.StatusCode, err)
    }

    // Пересобираем запрос (тело можно повторно читать только если оно было без тела или заранее буферизовано)
    // Здесь у нас все GET/POST с заранее буферизованным body, поэтому вызывающие методы должны передавать новый req.
    // Проще — вызывающие методы перед ретраем сами пересоздадут req. Для простоты — создадим копию без тела.
    req2, err := http.NewRequest(req.Method, req.URL.String(), req.Body)
    if err != nil {
        return nil, fmt.Errorf("recreate request: %w", err)
    }
    // Скопируем важные заголовки, кроме Authorization (мы поставим новый)
    for k, vals := range req.Header {
        if k == "Authorization" {
            continue
        }
        for _, v := range vals {
            req2.Header.Add(k, v)
        }
    }
    // Новый токен
    req2.Header.Set("Authorization", "Bearer "+c.bankToken)
    req2.Header.Set("User-Agent", "easyfund-vbank-client/1.0")
    req2.Header.Set("Accept", "application/json")

    return c.httpClient.Do(req2)
}

// ============= Demo helpers =============

func (c *VBankClient) GetRandomDemoClient() (*RandomClientResponse, error) {
    url := fmt.Sprintf("%s/auth/random-demo-client", c.baseURL)
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("create random-demo-client request: %w", err)
    }

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("do random-demo-client: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("random-demo-client HTTP %d: %s", resp.StatusCode, string(body))
    }

    var client RandomClientResponse
    if err := json.NewDecoder(resp.Body).Decode(&client); err != nil {
        return nil, fmt.Errorf("decode random-demo-client: %w", err)
    }
    return &client, nil
}

func (c *VBankClient) LoginClient(username, password string) (*ClientLoginResponse, error) {
    url := fmt.Sprintf("%s/auth/login", c.baseURL)
    payload := ClientLoginRequest{Username: username, Password: password}
    body, err := json.Marshal(payload)
    if err != nil {
        return nil, fmt.Errorf("marshal login: %w", err)
    }
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
    if err != nil {
        return nil, fmt.Errorf("create login request: %w", err)
    }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("User-Agent", "easyfund-vbank-client/1.0")
    req.Header.Set("Accept", "application/json")

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("do login: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        b, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("login HTTP %d: %s", resp.StatusCode, string(b))
    }
    var out ClientLoginResponse
    if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
        return nil, fmt.Errorf("decode login: %w", err)
    }
    return &out, nil
}

// ============= Consent =============

func (c *VBankClient) CreateConsent(clientID, requestingBank, requestingBankName string) (*ConsentResponse, error) {
    if requestingBank == "" {
        requestingBank = c.clientID
    }

    url := fmt.Sprintf("%s/account-consents/request", c.baseURL)
    payload := ConsentRequest{
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
    body, err := json.Marshal(payload)
    if err != nil {
        return nil, fmt.Errorf("marshal consent: %w", err)
    }

    // Первый запрос
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
    if err != nil {
        return nil, fmt.Errorf("create consent request: %w", err)
    }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("X-Requesting-Bank", requestingBank)

    resp, err := c.doWithBankAuth(req)
    if err != nil {
        return nil, fmt.Errorf("do consent: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        b, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("consent HTTP %d: %s", resp.StatusCode, string(b))
    }

    var out ConsentResponse
    if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
        return nil, fmt.Errorf("decode consent: %w", err)
    }
    return &out, nil
}

// ============= Accounts =============

func (c *VBankClient) GetAccountsWithConsent(clientID, consentID string) (*AccountsResponse, error) {
    url := fmt.Sprintf("%s/accounts?client_id=%s", c.baseURL, clientID)

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("create accounts request: %w", err)
    }
    req.Header.Set("x-consent-id", consentID)
    req.Header.Set("X-Requesting-Bank", c.clientID)

    resp, err := c.doWithBankAuth(req)
    if err != nil {
        return nil, fmt.Errorf("do accounts: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        b, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("accounts HTTP %d: %s", resp.StatusCode, string(b))
    }

    var out AccountsResponse
    if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
        return nil, fmt.Errorf("decode accounts: %w", err)
    }
    return &out, nil
}

// ============= Transactions =============

func (c *VBankClient) GetTransactionsWithConsent(accountID, clientID, consentID string, page, limit int) (*TransactionsResponse, error) {
    qs := fmt.Sprintf("client_id=%s", clientID)
    if page > 0 {
        qs += fmt.Sprintf("&page=%d", page)
    }
    if limit > 0 {
        qs += fmt.Sprintf("&limit=%d", limit)
    }
    url := fmt.Sprintf("%s/accounts/%s/transactions?%s", c.baseURL, accountID, qs)

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("create transactions request: %w", err)
    }
    req.Header.Set("x-consent-id", consentID)
    req.Header.Set("X-Requesting-Bank", c.clientID)

    resp, err := c.doWithBankAuth(req)
    if err != nil {
        return nil, fmt.Errorf("do transactions: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        b, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("transactions HTTP %d: %s", resp.StatusCode, string(b))
    }

    var out TransactionsResponse
    if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
        return nil, fmt.Errorf("decode transactions: %w", err)
    }
    return &out, nil
}

// ============= Balances =============

func (c *VBankClient) GetBalancesWithConsent(accountID, clientID, consentID string) (*BalancesResponse, error) {
    url := fmt.Sprintf("%s/accounts/%s/balances?client_id=%s", c.baseURL, accountID, clientID)

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("create balances request: %w", err)
    }
    req.Header.Set("x-consent-id", consentID)
    req.Header.Set("X-Requesting-Bank", c.clientID)

    resp, err := c.doWithBankAuth(req)
    if err != nil {
        return nil, fmt.Errorf("do balances: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        b, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("balances HTTP %d: %s", resp.StatusCode, string(b))
    }

    var out BalancesResponse
    if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
        return nil, fmt.Errorf("decode balances: %w", err)
    }
    return &out, nil
}
