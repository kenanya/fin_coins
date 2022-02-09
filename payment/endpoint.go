package payment

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type sendPaymentRequest struct {
	AccountID string
	Amount    float32
	ToAccount string
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
