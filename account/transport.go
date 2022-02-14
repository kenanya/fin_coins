package account

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kenanya/fin_coins/common"

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
)

// MakeHandler returns a handler for the account service.
func MakeHandler(ac Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
	}

	createAccountHandler := kithttp.NewServer(
		makeCreateAccountEndpoint(ac),
		decodeCreateAccountRequest,
		encodeResponse,
		opts...,
	)
	getAllAccountHandler := kithttp.NewServer(
		makeGetAllAccountEndpoint(ac),
		decodeGetAllAccountRequest,
		encodeResponse,
		opts...,
	)
	getAccountByIDHandler := kithttp.NewServer(
		makeGetAccountByIDEndpoint(ac),
		decodeGetAccountByIDRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()

	r.Handle("/account/v1/account", createAccountHandler).Methods("POST")
	r.Handle("/account/v1/account", getAllAccountHandler).Methods("GET")
	r.Handle("/account/v1/account/{id}", getAccountByIDHandler).Methods("GET")

	return r
}

func decodeCreateAccountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		ID       string  `json:"id"`
		Balance  float32 `json:"balance"`
		Currency string  `json:"currency"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	return createAccountRequest{
		ID:       body.ID,
		Balance:  body.Balance,
		Currency: body.Currency,
	}, nil
}

func decodeGetAccountByIDRequest(_ context.Context, r *http.Request) (interface{}, error) {

	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, common.ErrBadRoute
	}
	return getAccountByIDRequest{ID: id}, nil
}

func decodeGetAllAccountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return getAllAccountRequest{}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}

// encode errors from business-logic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case common.ErrUnknownAccount:
		w.WriteHeader(http.StatusNotFound)
	case common.ErrInvalidArgument:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
