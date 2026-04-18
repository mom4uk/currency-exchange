package services

import (
	"currency-exchange/internal/domain"
	"currency-exchange/internal/repositories"
)

type CurrencyService struct {
	exchangeRateRepository *repositories.ExchangeRateRepository
	currencyRepository     *repositories.CurrencyRepository
}

func CurrencyServiceNew(
	exchangeRateRepository *repositories.ExchangeRateRepository,
	currencyRepository *repositories.CurrencyRepository,
) *CurrencyService {
	return &CurrencyService{
		exchangeRateRepository: exchangeRateRepository,
		currencyRepository:     currencyRepository,
	}
}

func (c *CurrencyService) GetCurrencies() ([]domain.Currency, error) {
	return c.currencyRepository.GetCurrencies()
}

func (c *CurrencyService) AddCurrency(currency domain.Currency) (domain.Currency, error) {
	return c.currencyRepository.AddCurrency(currency)
}

func (c *CurrencyService) GetCurrencyByCode(code string) (domain.Currency, error) {
	return c.currencyRepository.GetCurrencyByCode(code)
}
