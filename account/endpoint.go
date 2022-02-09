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

type getAllAccountRequest struct {
}

type getAllAccountResponse struct {
	Accounts []Account `json:"accounts,omitempty"`
	Err      error     `json:"error,omitempty"`
}

func (r getAllAccountResponse) error() error { return r.Err }

func makeGetAllAccountEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(getAllAccountRequest)
		accounts, err := s.GetAllAccount()
		return getAllAccountResponse{Accounts: accounts, Err: err}, nil
	}
}

type getAccountByIDRequest struct {
	ID string
}

type getAccountByIDResponse struct {
	Account Account `json:"account,omitempty"`
	Err     error   `json:"error,omitempty"`
}

func (r getAccountByIDResponse) error() error { return r.Err }

func makeGetAccountByIDEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getAccountByIDRequest)
		account, err := s.GetAccountByID(req.ID)
		return getAccountByIDResponse{Account: account, Err: err}, nil
	}
}
