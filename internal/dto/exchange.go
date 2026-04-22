package dto

type CurencyExchangeResponse struct {
	BaseCurrency    CurrencyResponse `json:"baseCurrency"`
	TargetCurrency  CurrencyResponse `json:"targetCurrency"`
	Rate            string           `json:"rate"`
	Amount          string           `json:"amount"`
	ConvertedAmount string           `json:"convertedAmount"`
}
