package payment

import "context"

const (
	Incoming = "incoming"
	Outgoing = "outgoing"
)

// Payment represents an payment
type Payment struct {
	Account     string  `json:"account"`
	Amount      float64 `json:"amount"`
	FromAccount string  `json:"from_account"`
	Direction   string  `json:"direction"`
}

// Repository describes the persistence on payment model
type Repository interface {
	Create(ctx context.Context, payment Payment) error
	Get(ctx context.Context) ([]Payment, error)
}
