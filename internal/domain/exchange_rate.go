package domain

type AddExchangeRateRequest struct {
	BaseCurrencyCode   string
	TargetCurrencyCode string
	Rate               float64
}

type ExchangeRateResponse struct {
	ID             int
	BaseCurrency   Currency
	TargetCurrency Currency
	Rate           float64
}
