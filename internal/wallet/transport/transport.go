package transport

import (
	"context"
	"errors"
	"github.com/GGGY/wallet/internal/wallet/domain/account"
	"github.com/GGGY/wallet/internal/wallet/domain/payment"
	"github.com/GGGY/wallet/internal/wallet/service"
	"github.com/go-kit/kit/endpoint"
)

type SendPaymentRequest struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

type SendPaymentResponse struct {
	Err error `json:"err,omitempty"`
}

type GetPaymentsRequest struct{}

type GetPaymentsResponse struct {
	Payments []payment.Payment
	Err      error `json:"err,omitempty"`
}

type GetAccountsRequest struct{}

type GetAccountsResponse struct {
	Accounts []account.Account `json:"accounts"`
	Err      error             `json:"err,omitempty"`
}

type Endpoints struct {
	SendPayment endpoint.Endpoint
	GetAccounts endpoint.Endpoint
	GetPayments endpoint.Endpoint
}

func MakeEndpoints(s service.Service) Endpoints {
	return Endpoints{
		SendPayment: makeSendPaymentEndpoint(s),
		GetAccounts: makeGetAccountsEndpoint(s),
		GetPayments: makeGetPaymentsEndpoint(s),
	}
}

func makeSendPaymentEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(SendPaymentRequest)
		if !ok {
			return nil, errors.New("invalid request")
		}
		//todo: dont work with default payment
		_, err = s.SendPayment(ctx, payment.Payment{
			Account:     req.From,
			FromAccount: req.From,
			Amount:      req.Amount,
		})
		return SendPaymentResponse{}, err
	}
}

func makeGetAccountsEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, error error) {
		//req := request.(GetAccountsRequest)
		accounts, err := s.GetAccounts(ctx)
		return accounts, err
	}
}

func makeGetPaymentsEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		payments, err := s.GetPayments(ctx)
		return payments, err
	}
}
