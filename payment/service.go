package payment

import (
	"github.com/kenanya/fin_coins/account"
	"github.com/kenanya/fin_coins/common"
)

type Service interface {
	SendPayment(accountID string, amount float32, toAccount string) error
	GetAllPayment() ([]Payment, error)
}

type service struct {
	paymentRepo Repository
	accountRepo account.Repository
}

func (s *service) SendPayment(accountID string, amount float32, toAccount string) error {

	if accountID == "" || amount < 0 || toAccount == "" {
		return common.ErrInvalidArgument
	}

	fromAccountData, _ := s.accountRepo.GetAccountByID(accountID)
	toAccountData, _ := s.accountRepo.GetAccountByID(toAccount)

	if fromAccountData.ID == "" || toAccountData.ID == "" {
		return common.ErrAccountNotRegistered
	}

	if fromAccountData.Currency != toAccountData.Currency {
		return common.ErrDifferentCurrency
	}

	if fromAccountData.Balance < amount || fromAccountData.Balance <= 0 {
		return common.ErrInsufficientBalance
	}

	err := s.paymentRepo.SendPayment(accountID, amount, toAccount)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) GetAllPayment() ([]Payment, error) {

	payments, err := s.paymentRepo.GetAllPayment()
	if err != nil {
		return payments, err
	}

	return payments, nil
}

func NewService(paymentRepo Repository, accountRepo account.Repository) Service {
	return &service{
		paymentRepo: paymentRepo,
		accountRepo: accountRepo,
	}
}
