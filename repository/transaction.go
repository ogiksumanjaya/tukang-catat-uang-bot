package repository

import (
	"context"
	"database/sql"
	"github.com/ogiksumanjaya/helpers"
	"time"
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

func (t *TransactionRepository) GetTransactionList(ctx context.Context, username string, startDate time.Time, endDate time.Time) ([]helpers.ReportTransaction, error) {
	var transactions []helpers.ReportTransaction
	query := `SELECT
            ROW_NUMBER() OVER (ORDER BY t.created_at) AS No,  -- Nomor urut
            t.created_at AS Tanggal,                         -- Tanggal transaksi
            a.bank_name AS "Nama Akun",                   -- Nama akun bank
            c.name AS "Kategori",                       -- Nama kategori
            t.description AS Keterangan,                    -- Deskripsi transaksi
            CASE WHEN t.transaction_type = 'INCOME' THEN t.amount ELSE 0 END AS "Pemasukan",  -- Jumlah pendapatan
            CASE WHEN t.transaction_type = 'EXPENSE' THEN t.amount ELSE 0 END AS "Pengeluaran" -- Jumlah pengeluaran
				FROM
					Transaction t
						JOIN
					Account a ON t.account_id = a.id                -- Menghubungkan transaksi dengan akun
						JOIN
					Category c ON t.category_id = c.id              -- Menghubungkan transaksi dengan kategori
				WHERE t.username = $1 AND t.created_at BETWEEN $2 AND $3
				      ORDER BY t.created_at ASC
	`

	stmt, err := t.db.Prepare(query)
	if err != nil {
		return []helpers.ReportTransaction{}, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, username, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var transaction helpers.ReportTransaction
		err := rows.Scan(
			&transaction.No,
			&transaction.Date,
			&transaction.Account,
			&transaction.Category,
			&transaction.Description,
			&transaction.IncomeAmount,
			&transaction.ExpenseAmount,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
