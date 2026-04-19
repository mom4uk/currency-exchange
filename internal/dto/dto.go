package dto

import "currency-exchange/internal/domain"

type CurrencyResponse struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Sign string `json:"sign"`
}

type AddExchangeRateRequest struct {
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
