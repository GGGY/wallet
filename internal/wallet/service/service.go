package service

import (
	"context"
	"github.com/GGGY/wallet/internal/wallet/domain/account"
	"github.com/GGGY/wallet/internal/wallet/domain/payment"
	"github.com/GGGY/wallet/pkg/db"
	"github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
)

type Service interface {
	SendPayment(ctx context.Context, transfer Transfer) (string, error)
	GetPayments(ctx context.Context) ([]payment.Payment, error)
	GetAccounts(ctx context.Context) ([]account.Account, error)
}

func NewService(accountRepo account.Repository, paymentRepo payment.Repository, logger log.Logger) Service {
	return &wallet{
		accountRepo: accountRepo,
		paymentRepo: paymentRepo,
		logger:      logger,
	}
}

type wallet struct {
	db          *sqlx.DB
	accountRepo account.Repository
	paymentRepo payment.Repository
	logger      log.Logger
}

func (w *wallet) SendPayment(ctx context.Context, transfer Transfer) (string, error) {
	//todo:logic
	if err := db.WithTx(w.db, func(tx *sqlx.Tx) error {

		return nil
	}); err != nil {
		return "", err
	}

	return "", nil
}

func (w *wallet) GetPayments(ctx context.Context) ([]payment.Payment, error) {
	payments := make([]payment.Payment, 0)

	return payments, nil
}

func (w *wallet) GetAccounts(ctx context.Context) ([]account.Account, error) {
	accounts := make([]account.Account, 0)

	return accounts, nil
}
