package models

import (
	"time"
)

type CreditApplication struct {
	ApplicationID   int64     `json:"application_id" db:"application_id"`
	UserID          int64 `json:"user_id" db:"user_id"`
	BankID          int16     `json:"bank_id" db:"bank_id"`
	TypeCode        string    `json:"type_code" db:"type_code"`
	StatusCode      string    `json:"status_code" db:"status_code"`
	RequestedAmount string    `json:"requested_amount" db:"requested_amount"`
	LoanID          *int64    `json:"loan_id" db:"loan_id"`
	SubmittedAt     time.Time `json:"submitted_at" db:"submitted_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}
