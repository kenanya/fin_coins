package payment

import (
	"time"
)

type Payment struct {
	ID            string    `json:"id"`
	AccountID     string    `json:"account_id"`
	TransactionID string    `json:"transaction_id"`
	Amount        float32   `json:"amount"`
	ToAccount     string    `json:"to_account"`
	FromAccount   string    `json:"from_account"`
	Direction     string    `json:"direction"`
	CreatedAt     time.Time `json:"created_at"`
}

type Repository interface {
	SendPayment(accountID string, amount float32, toAccount string) error
	GetPaymentByDirection(direction string) ([]Payment, error)
	GetAllPayment() ([]Payment, error)
}
