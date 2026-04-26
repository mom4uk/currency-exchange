package services

import (
	"currency-exchange/internal/domain"
)

type CurrencyService struct {
	currencyRepository CurrencyRepository
}

func CurrencyServiceNew(
	currencyRepository CurrencyRepository,
) *CurrencyService {
	return &CurrencyService{
		currencyRepository: currencyRepository,
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
