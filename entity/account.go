package entity

type Account struct {
	ID       int     `json:"id" db:"id"`
	Username string  `json:"username" db:"username"`
	BankName string  `json:"bank_name" db:"bank_name"`
	Balance  float64 `json:"balance" db:"balance"`
}
