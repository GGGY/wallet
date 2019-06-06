package payment

import (
	"context"
	"github.com/jmoiron/sqlx"
)

// Payment directions
const (
	Incoming = "incoming"
	Outgoing = "outgoing"
)

// Payment represents an payment
type Payment struct {
	Account     string  `json:"account"`
	Amount      float64 `json:"amount"`
	FromAccount string  `json:"from_account,omitempty"`
	ToAccount   string  `json:"to_account,omitempty"`
	Direction   string  `json:"direction"`
}

// Repository describes the persistence on payment model
type Repository interface {
	Create(ctx context.Context, tx *sqlx.Tx, payment Payment) error
	Get(ctx context.Context) ([]Payment, error)
}
