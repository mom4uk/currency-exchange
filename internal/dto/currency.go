package dto

import (
	"currency-exchange/internal/domain"
	"regexp"
)

type CurrencyResponse struct {
	ID   string `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
	Sign string `json:"sign"`
}

type CurrencyRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Sign string `json:"sign"`
}

var latinOnly = regexp.MustCompile(`^[A-Za-z]+$`)

func ValidateCurrencyFields(req CurrencyRequest) error {
	if req.Code == "" || req.Name == "" || req.Sign == "" {
		return domain.ErrAbsenceOfCurrencyField
	}

	if len(req.Sign) > 3 {
		return domain.ErrInvalidCurrencySign
	}

	if !latinOnly.MatchString(req.Code) ||
		!latinOnly.MatchString(req.Name) ||
		!latinOnly.MatchString(req.Sign) {
		return domain.ErrInvalidCurrencyField
	}

	return nil
}
