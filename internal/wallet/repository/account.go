package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/GGGY/wallet/internal/wallet/domain/account"
	"github.com/jmoiron/sqlx"
)

type accountPg struct {
	db *sqlx.DB
}

func NewAccountPg(db *sqlx.DB) account.Repository {
	return &accountPg{
		db: db,
	}
}

func (r *accountPg) Create(ctx context.Context, account account.Account) error {
	_, err := r.db.Exec(`
		INSERT INTO account(id, balance, currency) VALUES ($1, $2, $3)
	`, account.ID, account.Balance, account.Currency)

	return err
}

func (r *accountPg) Update(ctx context.Context, tx *sqlx.Tx, account *account.Account) error {
	_, err := tx.Exec(`
		UPDATE account
		SET    balance = $1
		WHERE  id = $2
	`, account.Balance, account.ID)

	return err
}

func (r *accountPg) Get(ctx context.Context) ([]account.Account, error) {
	rows, err := r.db.Queryx(`
		SELECT
			id,
			balance,
			currency
		FROM account
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var accounts []account.Account
	for rows.Next() {
		var current account.Account
		if err = rows.StructScan(&current); err != nil {
			return nil, errors.New("can't scan account")
		}

		accounts = append(accounts, current)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

func (r *accountPg) GetByID(ctx context.Context, tx *sqlx.Tx, id string) (*account.Account, error) {
	var a account.Account
	err := tx.Get(&a, `
		SELECT
			id,
			balance,
			currency
		FROM account
		WHERE id = $1
		FOR UPDATE 
	`, id)

	if err == sql.ErrNoRows {
		return nil, account.ErrAccountNotFound
	}

	return &a, nil
}
