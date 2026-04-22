package domain

import "math/big"

type CurencyExchange struct {
	BaseCurrency    Currency `json:"baseCurrency"`
	TargetCurrency  Currency `json:"targetCurrency"`
	Rate            *big.Rat `json:"rate"`
	Amount          *big.Rat `json:"amount"`
	ConvertedAmount *big.Rat `json:"convertedAmount"`
}

type Exchange struct {
	BaseCurrencyCode   string
	TargetCurrencyCode string
	Rate               *big.Rat
}
