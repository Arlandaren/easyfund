package models

import (
	"time"
)

type Transaction struct {
	TransactionID int64     `json:"transaction_id" db:"transaction_id"`
	UserID        int64     `json:"user_id" db:"user_id"`
	BankID        int16     `json:"bank_id" db:"bank_id"`
	OccurredAt    time.Time `json:"occurred_at" db:"occurred_at"`
	Amount        string    `json:"amount" db:"amount"`
	Category      string    `json:"category" db:"category"`
	Description   string    `json:"description" db:"description"`
}