package account

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type createAccountRequest struct {
	ID       string
	Balance  float32
	Currency string
}

type createAccountResponse struct {
	Account Account `json:"account,omitempty"`
	Err     error   `json:"error,omitempty"`
}

func (r createAccountResponse) error() error { return r.Err }

func makeCreateAccountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createAccountRequest)
		account, err := s.CreateAccount(req.ID, req.Balance, req.Currency)
		return createAccountResponse{Account: account, Err: err}, nil
	}
}
