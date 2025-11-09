package models

import (
	"time"

)

type UserBankAccount struct {
	AccountID int64     `json:"account_id" db:"account_id"`
	UserID    int64 `json:"user_id" db:"user_id"`
	BankID    int16     `json:"bank_id" db:"bank_id"`
	Balance   string    `json:"balance" db:"balance"` // numeric(18,2)
	Currency  string    `json:"currency" db:"currency"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
