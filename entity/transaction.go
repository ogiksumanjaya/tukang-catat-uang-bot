package entity

import "time"

type Transaction struct {
	ID          int       `db:"id" json:"id"`
	Username    string    `db:"username" json:"username"`
	AccountId   int       `db:"account_id" json:"account_id"`
	CategoryId  int       `db:"category_id" json:"category_id"`
	Amount      float64   `db:"amount" json:"amount"`
	Description string    `db:"description" json:"description"`
	Type        string    `db:"type" json:"type"` // income, expense
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}
