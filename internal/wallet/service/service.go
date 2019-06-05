package service

import (
	"context"
	"github.com/GGGY/wallet/internal/wallet/domain/account"
	"github.com/GGGY/wallet/internal/wallet/domain/payment"
	"github.com/go-kit/kit/log"
)

type Service interface {
	SendPayment(ctx context.Context, payment payment.Payment) (string, error)
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
	accountRepo account.Repository
	paymentRepo payment.Repository
	logger      log.Logger
}

func (w *wallet) SendPayment(ctx context.Context, payment payment.Payment) (string, error) {
	//todo:logic
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
