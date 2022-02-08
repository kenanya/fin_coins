package account

import (
	"context"
	"errors"
	"time"
)

// type AccountID string

type Account struct {
	ID        string    `json:"id"`
	Balance   float32   `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var ErrUnknown = errors.New("unknown account")

type Repository interface {
	CreateAccount(account Account) (Account, error)
	GetAllAcount(ctx context.Context) ([]Account, error)
	GetAccountByID(ctx context.Context, id string) (Account, error)
	UpdateBalance(ctx context.Context, amount float32, accountID string) (Account, error)
}
