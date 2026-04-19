package domain

import "errors"

type ErrorResponse struct {
	Message string
}

var ErrAbsenceOfCode = errors.New("currency code missing")
var ErrIncorrectLengthOfCode = errors.New("incorrect length of code in url")
var ErrCurrencyNotFound = errors.New("currency not found")
