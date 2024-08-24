package repository

import (
	"context"
	"database/sql"
	"github.com/ogiksumanjaya/helpers"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}

func (t *TransactionRepository) CreateNewTransaction(ctx context.Context, value helpers.InputValue) error {
	query := `INSERT INTO transaction (username, category_id, account_id, amount, description, transaction_type) VALUES ($1, $2, $3, $4, $5, $6)`

	stmt, err := t.db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, value.Username, value.CategoryID, value.BankID, value.Amount, value.Note, value.Type)

	if err != nil {
		return err
	}

	return nil

}
