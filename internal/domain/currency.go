package domain

type Currency struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
	Sign string `json:"sign"`
}

type CurencyExchange struct {
	BaseCurrency    Currency `json:"baseCurrency"`
	TargetCurrency  Currency `json:"targetCurrency"`
	Rate            float64  `json:"rate"`
	Amount          float64  `json:"amount"`
	ConvertedAmount float64  `json:"convertedAmount"`
}

type CurrencyResponse struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Sign string `json:"sign"`
}
