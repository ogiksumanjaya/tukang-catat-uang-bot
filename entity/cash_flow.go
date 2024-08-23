package entity

type CashFlow struct {
	ID          int64   `db:"id" json:"id"`
	Amount      float64 `db:"amount" json:"amount"`
	Description string  `db:"description" json:"description"`
	Category    string  `db:"category" json:"category"`
	Type        string  `db:"type" json:"type"`       // in or out
	Account     string  `db:"account" json:"account"` // cash, bank, etc
	CreatedAt   string  `db:"created_at" json:"created_at"`
	CreatedBy   string  `db:"created_by" json:"created_by"`
}
