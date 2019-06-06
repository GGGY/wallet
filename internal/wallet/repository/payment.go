package repository

import (
	"context"
	"database/sql"
	"github.com/GGGY/wallet/internal/wallet/domain/payment"
	"github.com/jmoiron/sqlx"
)

type paymentPg struct {
	db *sqlx.DB
}

func NewPaymentPg(db *sqlx.DB) payment.Repository {
	return &paymentPg{
		db: db,
	}
}

func (r *paymentPg) Create(ctx context.Context, tx *sqlx.Tx, payment payment.Payment) error {
	_, err := tx.Exec(`
		INSERT INTO payment(account, amount, from_account, to_account, direction) VALUES ($1, $2, $3, $4, $5)
	`, payment.Account, payment.Amount, payment.FromAccount, payment.ToAccount, payment.Direction)

	return err
}

func (r *paymentPg) Get(ctx context.Context) ([]payment.Payment, error) {
	rows, err := r.db.Queryx(`
		SELECT
			account,
		    amount,
		    from_account,
		    to_account,
		    direction
		FROM payment
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	payments := make([]payment.Payment, 0)
	for rows.Next() {
		var current payment.Payment
		if err = rows.StructScan(&current); err != nil {
			return nil, err
		}

		payments = append(payments, current)
	}

	if err = rows.Err(); err != nil {
		if err == sql.ErrNoRows {
			return payments, payment.ErrPaymentNotFound
		}
		return nil, err
	}

	return payments, nil
}
