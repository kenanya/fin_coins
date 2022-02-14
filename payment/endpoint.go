package payment

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type sendPaymentRequest struct {
	AccountID string  `json:"account_id"`
	Amount    float32 `json:"amount"`
	ToAccount string  `json:"to_account"`
}

type sendPaymentResponse struct {
	Err error `json:"error,omitempty"`
}

func (r sendPaymentResponse) error() error { return r.Err }

func makeSendPaymentEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(sendPaymentRequest)
		err := s.SendPayment(req.AccountID, req.Amount, req.ToAccount)
		return sendPaymentResponse{Err: err}, nil
	}
}

type getAllPaymentRequest struct {
}

type getAllPaymentResponse struct {
	Payments []Payment `json:"payments,omitempty"`
	Err      error     `json:"error,omitempty"`
}

func (r getAllPaymentResponse) error() error { return r.Err }

func makeGetAllPaymentEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(getAllPaymentRequest)
		payments, err := s.GetAllPayment()
		return getAllPaymentResponse{Payments: payments, Err: err}, nil
	}
}
