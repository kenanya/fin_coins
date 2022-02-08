package account

import (
	"errors"
)

var ErrInvalidArgument = errors.New("invalid argument")

type Service interface {
	// BookNewCargo(origin location.UNLocode, destination location.UNLocode, deadline time.Time) (cargo.TrackingID, error)
	// LoadCargo(id cargo.TrackingID) (Cargo, error)
	CreateAccount(id string, balance float32, currency string) (Account, error)
}

type service struct {
	accountRepo Repository
	// locations      location.Repository
	// handlingEvents cargo.HandlingEventRepository
	// routingService routing.Service
}

func (s *service) CreateAccount(id string, balance float32, currency string) (Account, error) {
	var accountPass = Account{}
	if id == "" || balance < 0 || currency == "" {
		return accountPass, ErrInvalidArgument
	}
	accountPass = Account{
		ID:       id,
		Balance:  balance,
		Currency: currency,
	}

	accountRes, err := s.accountRepo.CreateAccount(accountPass)
	if err != nil {
		return accountRes, err
	}

	return accountRes, nil
}

// NewService creates a booking service with necessary dependencies.
func NewService(accountRepo Repository) Service {
	return &service{
		accountRepo: accountRepo,
	}
}
