package payment

import (
	"context"
	"time"
)

type Payment struct {
	ID          string    `json:"id"`
	AccountID   string    `json:"account_id"`
	Amount      float32   `json:"amount"`
	ToAccount   string    `json:"to_account"`
	FromAccount string    `json:"from_account_id"`
	Direction   string    `json:"direction"`
	CreatedAt   time.Time `json:"created_at"`
}

type Repository interface {
	SendPayment(ctx context.Context, payment Payment) error
	GetPaymentByDirection(ctx context.Context, direction string) ([]Payment, error)
	GetAllPayment(ctx context.Context) ([]Payment, error)
}
