package models

import (
	"time"

	"github.com/google/uuid"
)

type Loan struct {
	LoanID         int64     `json:"loan_id" db:"loan_id"`
	UserID         uuid.UUID `json:"user_id" db:"user_id"`
	OriginalAmount string    `json:"original_amount" db:"original_amount"` // numeric
	TakenAt        time.Time `json:"taken_at" db:"taken_at"`
	InterestRate   string    `json:"interest_rate" db:"interest_rate"`     // numeric(5,2)
	Status         string    `json:"status" db:"status"`
	Purpose        string    `json:"purpose" db:"purpose"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

type LoanSplit struct {
	SplitID            int64  `json:"split_id" db:"split_id"`
	LoanID             int64  `json:"loan_id" db:"loan_id"`
	BankID             int16  `json:"bank_id" db:"bank_id"`
	SplitAmount        string `json:"split_amount" db:"split_amount"`
	RemainingPrincipal string `json:"remaining_principal" db:"remaining_principal"`
}
