package models

type Bank struct {
	BankID int16  `json:"bank_id" db:"bank_id"`
	Code   string `json:"code" db:"code"`
	Name   string `json:"name" db:"name"`
}
