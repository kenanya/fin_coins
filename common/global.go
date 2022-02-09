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

// var httpErrorMap = map[int][]int{
// 	http.StatusBadRequest:          {30001, 30002, 40002, 40101, 41101},
// 	http.StatusUnauthorized:        {41102},
// 	http.StatusForbidden:           {},
// 	http.StatusNotFound:            {40005, 40006},
// 	http.StatusInternalServerError: {50098, 50099},
// }

// func GetHttpError(errCode int) int {
// 	curHttpError := http.StatusInternalServerError
// 	for k, v := range httpErrorMap {
// 		if helper.FindIntInSlice(v, errCode) {
// 			curHttpError = k
// 			break
// 		}
// 	}
// 	return curHttpError
// }
