package service

import (
	"context"
	"errors"
	"github.com/GGGY/wallet/internal/wallet/domain/account"
	"github.com/GGGY/wallet/internal/wallet/domain/payment"
	"github.com/GGGY/wallet/pkg/db"
	"github.com/go-kit/kit/log"
	"github.com/jmoiron/sqlx"
)

type Service interface {
	Transfer(ctx context.Context, from string, to string, amount float64) error
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

func (w *wallet) Transfer(ctx context.Context, from string, to string, amount float64) error {
	//todo: don't use float for math
	if err := db.WithTx(w.db, func(tx *sqlx.Tx) error {
		balanceFrom, err := w.accountRepo.GetByID(ctx, tx, from)
		if err != nil {
			return err
		}

		restBalance := balanceFrom.Balance - amount
		if restBalance < 0 {
			return errors.New("balance is not enough")
		}

		balanceTo, err := w.accountRepo.GetByID(ctx, tx, to)
		if err != nil {
			return err
		}

		if balanceFrom.Currency != balanceTo.Currency {
			return errors.New("different currencies")
		}

		balanceFrom.Balance -= amount
		balanceTo.Balance += amount

		outgoingPayment := payment.Payment{
			Account:   balanceFrom.ID,
			Amount:    amount,
			ToAccount: balanceTo.ID,
			Direction: payment.Outgoing,
		}

		incomingPayment := payment.Payment{
			Account:     balanceTo.ID,
			Amount:      amount,
			FromAccount: balanceFrom.ID,
			Direction:   payment.Incoming,
		}

		err = w.accountRepo.Update(ctx, tx, balanceFrom)
		if err != nil {
			return err
		}

		err = w.paymentRepo.Create(ctx, tx, outgoingPayment)
		if err != nil {
			return err
		}

		err = w.accountRepo.Update(ctx, tx, balanceTo)
		if err != nil {
			return err
		}

		err = w.paymentRepo.Create(ctx, tx, incomingPayment)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (w *wallet) GetPayments(ctx context.Context) ([]payment.Payment, error) {
	payments := make([]payment.Payment, 0)

	return payments, nil
}

func (w *wallet) GetAccounts(ctx context.Context) ([]account.Account, error) {
	accounts := make([]account.Account, 0)

	return accounts, nil
}
