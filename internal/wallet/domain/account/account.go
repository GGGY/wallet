package account

import "context"

// Account represents an account
type Account struct {
	ID       string  `json:"id"`
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
}

// Repository describes the persistence on account model
type Repository interface {
	Create(ctx context.Context, account Account) error
	ChangeBalance(ctx context.Context, id string, amount float64) error
	Get(ctx context.Context) ([]Account, error)
}
