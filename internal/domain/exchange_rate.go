package domain

type ExchangeRate struct {
	ID               int
	BaseCurrencyId   int
	TargetCurrencyId int
	Rate             float64
}
