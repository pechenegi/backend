package models

import "time"

type Currency int64

func (c Currency) Float64() float64 {
	return float64(c) / 100
}

func ToCurrency(f float64) Currency {
	return Currency((f * 100) + 0.5)
}

// DebtInfo holds initial information about user's debt to SO
// which is required to make debt calculation.
type DebtInfo struct {
	ID               int
	StartDate        time.Time
	ContractDuration int
	TechBlockPrice   Currency
	ITBlockPrice     Currency
	HadStartUp       bool
	StartUpAmount    Currency
	StartUpDuration  int
}

// DebtStats holds calculated information about user's debt.
type DebtStats struct {
	StudyLoan    LoanStat
	StartUpLoan  LoanStat
	CalculatedAt time.Time
}

// LoanStat holds information about single loan.
type LoanStat struct {
	Loan       float64
	Paid       float64
	DaysTotal  int
	DaysPassed int
}

// User holds user-related information.
type User struct {
	ID       string `json:"id,omitempty"`
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}
