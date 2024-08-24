package repository

import (
	"context"
	"database/sql"
	"github.com/ogiksumanjaya/entity"
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

func (a *AccountRepository) GetAccountList(ctx context.Context, username string) ([]entity.Account, error) {
	var accounts []entity.Account

	query := `SELECT * FROM account WHERE username = $1`

	stmt, err := a.db.Prepare(query)
	if err != nil {
		return []entity.Account{}, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var account entity.Account
		err := rows.Scan(&account.ID, &account.Username, &account.BankName, &account.Balance)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (a *AccountRepository) GetAccountByName(ctx context.Context, account entity.Account) (entity.Account, error) {
	query := `SELECT * FROM account WHERE username = $1 AND bank_name = $2`

	stmt, err := a.db.Prepare(query)
	if err != nil {
		return entity.Account{}, err
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, account.Username, account.BankName).Scan(&account.ID, &account.Username, &account.BankName, &account.Balance)
	if err != nil {
		return entity.Account{}, err
	}

	return account, nil
}

func (a *AccountRepository) UpdateBalance(ctx context.Context, account entity.Account) error {
	query := `UPDATE account SET balance = $1 WHERE id = $2`

	stmt, err := a.db.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, account.Balance, account.ID)

	if err != nil {
		return err
	}

	return nil
}
