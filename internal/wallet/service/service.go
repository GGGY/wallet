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

func NewService(db *sqlx.DB, accountRepo account.Repository, paymentRepo payment.Repository, logger log.Logger) Service {
	return &wallet{
		db:          db,
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
	w.logger.Log("from", from)
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
	payments, err := w.paymentRepo.Get(ctx)
	if err != nil {
		if err == payment.ErrPaymentNotFound {
			return payments, nil
		}
		w.logger.Log("can't get payments from db: ", err)
		return payments, err
	}

	return payments, nil
}

func (w *wallet) GetAccounts(ctx context.Context) ([]account.Account, error) {
	accounts, err := w.accountRepo.Get(ctx)
	if err != nil {
		if err == account.ErrAccountNotFound {
			return accounts, nil
		}
		w.logger.Log("can't get accounts from db: ", err)
		return accounts, err
	}

	return accounts, nil
}
