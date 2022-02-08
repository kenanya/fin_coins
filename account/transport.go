package account

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"

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
	// loadCargoHandler := kithttp.NewServer(
	// 	makeLoadCargoEndpoint(bs),
	// 	decodeLoadCargoRequest,
	// 	encodeResponse,
	// 	opts...,
	// )
	// requestRoutesHandler := kithttp.NewServer(
	// 	makeRequestRoutesEndpoint(bs),
	// 	decodeRequestRoutesRequest,
	// 	encodeResponse,
	// 	opts...,
	// )
	// assignToRouteHandler := kithttp.NewServer(
	// 	makeAssignToRouteEndpoint(bs),
	// 	decodeAssignToRouteRequest,
	// 	encodeResponse,
	// 	opts...,
	// )
	// changeDestinationHandler := kithttp.NewServer(
	// 	makeChangeDestinationEndpoint(bs),
	// 	decodeChangeDestinationRequest,
	// 	encodeResponse,
	// 	opts...,
	// )
	// listCargosHandler := kithttp.NewServer(
	// 	makeListCargosEndpoint(bs),
	// 	decodeListCargosRequest,
	// 	encodeResponse,
	// 	opts...,
	// )
	// listLocationsHandler := kithttp.NewServer(
	// 	makeListLocationsEndpoint(bs),
	// 	decodeListLocationsRequest,
	// 	encodeResponse,
	// 	opts...,
	// )

	r := mux.NewRouter()

	r.Handle("/account/v1/account", createAccountHandler).Methods("POST")
	// r.Handle("/booking/v1/cargos", listCargosHandler).Methods("GET")
	// r.Handle("/booking/v1/cargos/{id}", loadCargoHandler).Methods("GET")
	// r.Handle("/booking/v1/cargos/{id}/request_routes", requestRoutesHandler).Methods("GET")
	// r.Handle("/booking/v1/cargos/{id}/assign_to_route", assignToRouteHandler).Methods("POST")
	// r.Handle("/booking/v1/cargos/{id}/change_destination", changeDestinationHandler).Methods("POST")
	// r.Handle("/booking/v1/locations", listLocationsHandler).Methods("GET")

	return r
}

var errBadRoute = errors.New("bad route")

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
	case ErrUnknown:
		w.WriteHeader(http.StatusNotFound)
	case ErrInvalidArgument:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
