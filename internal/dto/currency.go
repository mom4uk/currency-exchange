package dto

import (
	"currency-exchange/internal/domain"
	"regexp"
	"strings"
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

var codeOnly = regexp.MustCompile(`^[A-Za-z]+$`)

func ValidateCurrencyFields(req CurrencyRequest) error {
	if strings.TrimSpace(req.Code) == "" ||
		strings.TrimSpace(req.Name) == "" ||
		strings.TrimSpace(req.Sign) == "" {
		return domain.ErrAbsenceOfCurrencyField
	}

	if len(req.Sign) > 3 {
		return domain.ErrInvalidCurrencySign
	}

	// только code должен быть латиницей
	if !codeOnly.MatchString(req.Code) {
		return domain.ErrInvalidCurrencyField
	}

	return nil
}
