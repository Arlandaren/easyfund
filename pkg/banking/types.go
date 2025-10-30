package banking

import "time"

// Virtual Bank API Models
type RandomClientResponse struct {
	PersonID string `json:"person_id"`
	FullName string `json:"full_name"`
	Password string `json:"password"`
}

type BankTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ClientID    string `json:"client_id"`
	ExpiresIn   int    `json:"expires_in"`
}

type ClientLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ClientLoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ClientID    string `json:"client_id"`
}

type Account struct {
	AccountID   string `json:"account_id"`
	Currency    string `json:"currency"`
	AccountType string `json:"account_type"`
	Nickname    string `json:"nickname"`
	Servicer    *struct {
		SchemeName     string `json:"scheme_name"`
		Identification string `json:"identification"`
	} `json:"servicer,omitempty"`
}

type Balance struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
	Type     string `json:"type,omitempty"`
}

type Transaction struct {
	TransactionID          string    `json:"transaction_id"`
	Amount                 string    `json:"amount"`
	Currency               string    `json:"currency"`
	CreditDebitIndicator   string    `json:"credit_debit_indicator"`
	Status                 string    `json:"status"`
	BookingDateTime        time.Time `json:"booking_date_time"`
	ValueDateTime          time.Time `json:"value_date_time"`
	TransactionInformation string    `json:"transaction_information,omitempty"`
	BankTransactionCode    *struct {
		Code        string `json:"code"`
		SubCode     string `json:"sub_code"`
		Description string `json:"description"`
	} `json:"bank_transaction_code,omitempty"`
}

// OpenBanking Standard Responses
type AccountsResponse struct {
	Accounts  []Account `json:"accounts"`
	Links Links     `json:"links"`
	Meta  Meta      `json:"meta"`
}

type TransactionsResponse struct {
	Data  []Transaction `json:"data"`
	Links Links         `json:"links"`
	Meta  Meta          `json:"meta"`
}

type BalancesResponse struct {
	Data  []Balance `json:"data"`
	Links Links     `json:"links"`
	Meta  Meta      `json:"meta"`
}

type Links struct {
	Self  string `json:"self"`
	First string `json:"first,omitempty"`
	Prev  string `json:"prev,omitempty"`
	Next  string `json:"next,omitempty"`
	Last  string `json:"last,omitempty"`
}

type Meta struct {
	TotalPages int `json:"total_pages"`
}

// Consent Management
type ConsentRequest struct {
	ClientID           string   `json:"client_id"`
	Permissions        []string `json:"permissions"`
	Reason             string   `json:"reason"`
	RequestingBank     string   `json:"requesting_bank"`
	RequestingBankName string   `json:"requesting_bank_name"`
}

type ConsentResponse struct {
	Status       string `json:"status"`
	ConsentID    string `json:"consent_id"`
	AutoApproved bool   `json:"auto_approved"`
}

// Application Models
type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	PersonID string `json:"person_id,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Role     string `json:"role"`
}

type AuthResponse struct {
	Token    string `json:"token"`
	TokenType string `json:"token_type"`
	User     User   `json:"user"`
}

type DemoClientLoginRequest struct {
	PersonID string `json:"person_id"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
}

// JWT Claims
type JWTClaims struct {
    // Идентификаторы
    UserID   string `json:"user_id"`              // внутренний ID пользователя в вашем сервисе
    ClientID string `json:"client_id,omitempty"`  // ключевое поле: ID клиента в Virtual Bank
    PersonID string `json:"person_id,omitempty"`  // опционально, если используется отдельно от client_id

    // Контактные/ролевые данные
    Email string `json:"email,omitempty"`
    Mode  string `json:"mode,omitempty"` // "demo" или "regular"
    Role  string `json:"role,omitempty"`

    // Токены и интеграции
    VBankAccessToken string `json:"vbank_token,omitempty"` // если где-то используется клиентский UI-токен VB (обычно не нужен серверу)

    // Стандартные JWT-поля для валидации (опционально, но полезно иметь)
    Issuer   string `json:"iss,omitempty"`
    Subject  string `json:"sub,omitempty"`
    Audience string `json:"aud,omitempty"`
    IssuedAt int64  `json:"iat,omitempty"`
    Expires  int64  `json:"exp,omitempty"`
    NotBefore int64 `json:"nbf,omitempty"`
}

// Error Response
type ErrorResponse struct {
    Error   string `json:"error"`
    Message string `json:"message,omitempty"`
    Code    int    `json:"code"`
}