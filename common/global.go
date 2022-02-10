package common

import "errors"

const (
	ErrPriceInvalid int = 30001
	ErrInternal     int = 50098
	ErrGeneral      int = 50099
)

var (
	ErrInvalidArgument      = errors.New("invalid argument")
	ErrBadRoute             = errors.New("bad route")
	ErrUnknown              = errors.New("unknown account")
	ErrAccountNotRegistered = errors.New("sender or receiver not exists")
	ErrDifferentCurrency    = errors.New("different currency between sender and receiver")
	ErrInsufficientBalance  = errors.New("insufficient balance")
)
