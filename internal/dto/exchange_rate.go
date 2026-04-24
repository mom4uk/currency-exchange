package dto

import (
	"currency-exchange/internal/domain"
	"strings"
)

type AddExchangeRateRequest struct {
	BaseCurrencyCode   string
	TargetCurrencyCode string
	Rate               string
}

type ExchangeRateResponse struct {
	ID             int              `json:"id"`
	BaseCurrency   CurrencyResponse `json:"baseCurrency"`
	TargetCurrency CurrencyResponse `json:"targetCurrency"`
	Rate           string           `json:"rate"`
}

type UpdateExchangeRateRequest struct {
	Rate string
}

func ValidateExchangeRateFields(req AddExchangeRateRequest) error {
	base := strings.TrimSpace(req.BaseCurrencyCode)
	target := strings.TrimSpace(req.TargetCurrencyCode)
	rate := strings.TrimSpace(req.Rate)

	if base == "" || target == "" || rate == "" {
		return domain.ErrAbsenceOfExchangeRateField
	}

	return nil
}

func ValidateExchangeRateFieldsForUpdate(req UpdateExchangeRateRequest) error {
	rate := strings.TrimSpace(req.Rate)
	if rate == "" {
		return domain.ErrAbsenceOfExchangeRateFieldForUpdate
	}
	return nil
}
