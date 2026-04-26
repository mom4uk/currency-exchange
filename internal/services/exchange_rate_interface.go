package services

import (
	"currency-exchange/internal/domain"
	"math/big"
)

type ExchangeRateRepository interface {
	UpdateExchangeRate(baseCurrency domain.Currency, targetCurrency domain.Currency, rate *big.Rat) (domain.ExchangeRate, error)
	GetExchangeRate(baseCurrencyId int, targetCurrencyId int) (domain.ExchangeRate, bool, error)
	AddExchangeRates(baseCurrency domain.Currency, targetCurrency domain.Currency, rate *big.Rat) (domain.ExchangeRate, error)
	GetExchangeRates() ([]domain.ExchangeRate, error)
}
