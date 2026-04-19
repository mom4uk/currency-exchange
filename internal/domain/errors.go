package domain

import "errors"

type ErrorResponse struct {
	Message string
}

var (
	ErrAbsenceOfCode         = errors.New("currency code missing")
	ErrIncorrectLengthOfCode = errors.New("incorrect length of code in url")
	ErrCurrencyNotFound      = errors.New("currency not found")
	ErrAbsenceOfField        = errors.New("necessary field were not provided")
	ErrCurrencyAlreadyExists = errors.New("currency already exists")
)
