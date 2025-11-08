package models

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	TransactionID int64     `json:"transaction_id" db:"transaction_id"`
	UserID        uuid.UUID `json:"user_id" db:"user_id"`
	BankID        int16     `json:"bank_id" db:"bank_id"`
	OccurredAt    time.Time `json:"occurred_at" db:"occurred_at"`
	Amount        string    `json:"amount" db:"amount"` // numeric(18,2), расход — отрицательное
	Category      string    `json:"category" db:"category"`
	Description   string    `json:"description" db:"description"`
}
