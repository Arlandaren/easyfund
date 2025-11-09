package models

import (
	"time"

)

type LoanPayment struct {
	PaymentID   int64     `json:"payment_id" db:"payment_id"`
	LoanID      int64     `json:"loan_id" db:"loan_id"`
	UserID      int64 `json:"user_id" db:"user_id"`
	PaidAt      time.Time `json:"paid_at" db:"paid_at"`
	TotalAmount string    `json:"total_amount" db:"total_amount"`
	Comment     string    `json:"comment" db:"comment"`
}

type PaymentAllocation struct {
	AllocationID  int64  `json:"allocation_id" db:"allocation_id"`
	PaymentID     int64  `json:"payment_id" db:"payment_id"`
	SplitID       int64  `json:"split_id" db:"split_id"`
	PrincipalPaid string `json:"principal_paid" db:"principal_paid"`
	InterestPaid  string `json:"interest_paid" db:"interest_paid"`
}
