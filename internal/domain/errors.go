package domain

import "errors"

type ErrorResponse struct {
	Message string
}

var (
	ErrAbsenceOfCode              = errors.New("currency code missing")
	ErrIncorrectLengthOfCode      = errors.New("incorrect length of code in url")
	ErrCurrencyNotFound           = errors.New("currency not found")
	ErrAbsenceOfCurrencyField     = errors.New("necessary field were not provided")
	ErrCurrencyAlreadyExists      = errors.New("currency already exists")
	ErrExchangeRateNotFound       = errors.New("exchange rate not found")
	ErrExchangeRateAlreadyExists  = errors.New("exchange rate already exists")
	ErrAbsenceOfExchangeRateField = errors.New("necessary field were not provided")
)
