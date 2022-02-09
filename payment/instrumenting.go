package payment

import (
	"time"

	"github.com/go-kit/kit/metrics"
)

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	Service
}

// NewInstrumentingService returns an instance of an instrumenting Service.
func NewInstrumentingService(counter metrics.Counter, latency metrics.Histogram, s Service) Service {
	return &instrumentingService{
		requestCount:   counter,
		requestLatency: latency,
		Service:        s,
	}
}

func (s *instrumentingService) SendPayment(accountID string, amount float32, toAccount string) error {
	defer func(begin time.Time) {
		s.requestCount.With("method", "send payment").Add(1)
		s.requestLatency.With("method", "send payment").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.SendPayment(accountID, amount, toAccount)
}

func (s *instrumentingService) GetAllPayment() ([]Payment, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "get all payment").Add(1)
		s.requestLatency.With("method", "get all payment").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.Service.GetAllPayment()
}
