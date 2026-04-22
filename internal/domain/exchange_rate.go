package domain

import "math/big"

type ExchangeRate struct {
	ID               int
	BaseCurrencyId   int
	TargetCurrencyId int
	Rate             *big.Rat
}
