package models

type UserStatsDTO struct {
	User         *User  `json:"user"`
	TotalBalance string `json:"total_balance"`
	TotalDebt    string `json:"total_debt"`
	CreditRating string `json:"credit_rating"`
}

type LoanDetailDTO struct {
	Loan           *Loan       `json:"loan"`
	Splits         []LoanSplit `json:"splits"`
	PercentPaid    float64     `json:"percent_paid"`
	RemainingDebt  string      `json:"remaining_debt"`
	PaymentHistory []LoanPayment `json:"payment_history"`
}

type TransactionHistoryDTO struct {
	Transactions []Transaction `json:"transactions"`
	TotalSpent   string        `json:"total_spent"`
}
