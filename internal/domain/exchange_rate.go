package domain

type AddExchangeRateRequest struct {
	BaseCurrencyCode   string
	TargetCurrencyCode string
	Rate               float64
}

type ExchangeRateResponse struct {
	ID             int      `json:"id"`
	BaseCurrency   Currency `json:"baseCurrency"`
	TargetCurrency Currency `json:"targetCurrency"`
	Rate           float64  `json:"rate"`
}

type ExchangeRate struct {
	ID               int
	BaseCurrencyId   int
	TargetCurrencyId int
	Rate             float64
}
