package account

import "github.com/kenanya/fin_coins/common"

// var ErrInvalidArgument = errors.New("invalid argument")

type Service interface {
	CreateAccount(id string, balance float32, currency string) (Account, error)
	GetAllAccount() ([]Account, error)
	GetAccountByID(id string) (Account, error)
}

type service struct {
	accountRepo Repository
}

func (s *service) CreateAccount(id string, balance float32, currency string) (Account, error) {
	var accountPass = Account{}
	if id == "" || balance < 0 || currency == "" {
		return accountPass, common.ErrInvalidArgument
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

func (s *service) GetAllAccount() ([]Account, error) {

	accounts, err := s.accountRepo.GetAllAccount()
	if err != nil {
		return accounts, err
	}

	return accounts, nil
}

func (s *service) GetAccountByID(id string) (Account, error) {
	var account = Account{}
	if id == "" {
		return account, common.ErrInvalidArgument
	}
	account, err := s.accountRepo.GetAccountByID(id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return account, common.ErrUnknownAccount
		}
		return account, err
	}

	return account, nil
}

func NewService(accountRepo Repository) Service {
	return &service{
		accountRepo: accountRepo,
	}
}
