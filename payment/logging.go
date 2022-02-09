package payment

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

func (s *loggingService) SendPayment(accountID string, amount float32, toAccount string) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "send payment",
			"accountID", accountID,
			"amount", amount,
			"toAccount", toAccount,
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.SendPayment(accountID, amount, toAccount)
}

func (s *loggingService) GetAllPayment() (payments []Payment, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "send payment",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.GetAllPayment()
}
