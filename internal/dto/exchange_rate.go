package dto

import "currency-exchange/internal/domain"

type AddExchangeRateRequest struct {
	BaseCurrencyCode   string
	TargetCurrencyCode string
	Rate               string
}

type ExchangeRate struct {
	BaseCurrencyCode   string
	TargetCurrencyCode string
	Rate               float64
}

type ExchangeRateResponse struct {
	ID             int             `json:"id"`
	BaseCurrency   domain.Currency `json:"baseCurrency"`
	TargetCurrency domain.Currency `json:"targetCurrency"`
	Rate           float64         `json:"rate"`
}

type UpdateExchangeRateRequest struct {
	Rate string
}

func ValidateExchangeRateFields(req AddExchangeRateRequest) error {
	if req.BaseCurrencyCode == "" || req.TargetCurrencyCode == "" || req.Rate == "" {
		return domain.ErrAbsenceOfExchangeRateField
	}
	return nil
}

func ValidateExchangeRateFieldsForUpdate(req UpdateExchangeRateRequest) error {
	if req.Rate == "" {
		return domain.ErrAbsenceOfExchangeRateFieldForUpdate
	}
	return nil
}
