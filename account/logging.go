package account

import (
	"time"

	"github.com/go-kit/kit/log"
)

type loggingService struct {
	logger log.Logger
	Service
}

// NewLoggingService returns a new instance of a logging Service.
func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) CreateAccount(id string, balance float32, currency string) (accountRes Account, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "create account",
			"id", id,
			"balance", balance,
			"currency", currency,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.CreateAccount(id, balance, currency)
}

func (s *loggingService) GetAllAccount() (accounts []Account, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "get all account",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.GetAllAccount()
}

func (s *loggingService) GetAccountByID(id string) (account Account, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "get account by id",
			"id", id,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.GetAccountByID(id)
}
