package payment

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

// MakeHandler returns a handler for the payment service.
func MakeHandler(ac Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
	}

	sendPaymentHandler := kithttp.NewServer(
		makeSendPaymentEndpoint(ac),
		decodeSendPaymentRequest,
		encodeResponse,
		opts...,
	)
	// loadCargoHandler := kithttp.NewServer(
	// 	makeLoadCargoEndpoint(bs),
	// 	decodeLoadCargoRequest,
	// 	encodeResponse,
	// 	opts...,
	// )

	r := mux.NewRouter()

	r.Handle("/payment/v1/payment", sendPaymentHandler).Methods("POST")
	// r.Handle("/booking/v1/cargos", listCargosHandler).Methods("GET")
	// r.Handle("/booking/v1/cargos/{id}", loadCargoHandler).Methods("GET")

	return r
}

func decodeSendPaymentRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		AccountID string  `json:"account_id"`
		Amount    float32 `json:"amount"`
		ToAccount string  `json:"to_account"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	return sendPaymentRequest{
		AccountID: body.AccountID,
		Amount:    body.Amount,
		ToAccount: body.ToAccount,
	}, nil
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
	case common.ErrUnknown:
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
