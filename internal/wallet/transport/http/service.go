package http

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/GGGY/wallet/internal/wallet/transport"
	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"net/http"
)

var (
	ErrBadRouting = errors.New("bad routing")
)

type errorer interface {
	error() error
}

func NewService(endpoints transport.Endpoints, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	options := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	r.Methods("POST").Path("/send").Handler(kithttp.NewServer(
		endpoints.SendPayment,
		decodeSendRequest,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/payments").Handler(kithttp.NewServer(
		endpoints.GetPayments,
		decodeGetPaymentsRequest,
		encodeResponse,
		options...,
	))

	r.Methods("GET").Path("/accounts").Handler(kithttp.NewServer(
		endpoints.GetAccounts,
		decodeGetAccountsRequest,
		encodeResponse,
		options...,
	))

	return r
}

func decodeSendRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	//todo: logic

	return nil, nil
}

func decodeGetPaymentsRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	// todo: logic

	return nil, nil
}

func decodeGetAccountsRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	// todo: logic

	return nil, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case errors.New("some text"): //todo: send errors
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
