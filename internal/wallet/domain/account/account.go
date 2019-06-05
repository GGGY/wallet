package account

import (
	"context"
	"errors"
	"github.com/jmoiron/sqlx"
)

var (
	ErrAccountNotFound = errors.New("balance not found")
)

// Account represents an account
type Account struct {
	ID       string  `json:"id"`
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
}

// Repository describes the persistence on account model
type Repository interface {
	Create(ctx context.Context, account Account) error
	Update(ctx context.Context, tx *sqlx.Tx, account *Account) error
	Get(ctx context.Context) ([]Account, error)
	GetByID(ctx context.Context, tx *sqlx.Tx, id string) (*Account, error)
}
