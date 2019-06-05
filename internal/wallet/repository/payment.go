package repository

import (
	"context"
	"errors"
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
		INSERT INTO payment(account, amount, from_account, direction) VALUES ($1, $2, $3, $4)
	`, payment.Account, payment.Amount, payment.FromAccount, payment.Direction)

	return err
}

func (r *paymentPg) Get(ctx context.Context) ([]payment.Payment, error) {
	rows, err := r.db.Queryx(`
		SELECT
			payment,
		    amount,
		    from_account,
		    direction
		FROM payment
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var payments []payment.Payment
	for rows.Next() {
		var current payment.Payment
		if err = rows.StructScan(&current); err != nil {
			return nil, errors.New("can't scan account")
		}

		payments = append(payments, current)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return payments, nil
}
