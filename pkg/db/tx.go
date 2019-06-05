package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type TransactionFunc func(tx *sqlx.Tx) error

func WithTx(db *sqlx.DB, fn TransactionFunc) (err error) {
	var tx *sqlx.Tx
	if tx, err = db.Beginx(); err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and repanic
			if err := tx.Rollback(); err != nil {
				panic(fmt.Sprintf("Recover: %s, error on rollback: %s", p, err))
			}
			panic(p)
		} else if err != nil {
			// something went wrong, rollback
			if err := tx.Rollback(); err != nil {
				panic(err)
			}
		} else {
			// all good, commit
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return err
}
