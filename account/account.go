package account

import (
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

type Repository interface {
	CreateAccount(account Account) (Account, error)
	GetAllAccount() ([]Account, error)
	GetAccountByID(id string) (Account, error)
	DeleteAccount(id string) error
}
