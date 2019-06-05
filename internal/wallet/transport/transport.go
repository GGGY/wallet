package transport

import (
	"context"
	"errors"
	"github.com/GGGY/wallet/internal/wallet/domain/account"
	"github.com/GGGY/wallet/internal/wallet/domain/payment"
	"github.com/GGGY/wallet/internal/wallet/service"
	"github.com/go-kit/kit/endpoint"
)

type TransferRequest struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

type TransferResponse struct {
	Success bool  `json:"success"`
	Err     error `json:"err,omitempty"`
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
	Transfer    endpoint.Endpoint
	GetAccounts endpoint.Endpoint
	GetPayments endpoint.Endpoint
}

func MakeEndpoints(s service.Service) Endpoints {
	return Endpoints{
		Transfer:    makeTransferEndpoint(s),
		GetAccounts: makeGetAccountsEndpoint(s),
		GetPayments: makeGetPaymentsEndpoint(s),
	}
}

func makeTransferEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(TransferRequest)
		if !ok {
			return nil, errors.New("invalid request")
		}
		//todo: check amount is positive
		err = s.Transfer(ctx, req.From, req.To, req.Amount)
		var success bool
		if err == nil {
			success = true
		}

		return TransferResponse{Success: success}, err
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
