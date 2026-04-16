package services

import (
	"currency-exchange/internal/domain"
	"currency-exchange/internal/repositories"
)

type CurrencyService struct {
	repository *repositories.CurrencyRepository
}

func CurrencyServiceNew(repo *repositories.CurrencyRepository) *CurrencyService {
	return &CurrencyService{
		repository: repo,
	}
}

func (c *CurrencyService) GetCurrencies() ([]domain.Currency, error) {
	return c.repository.GetCurrencies()
}

func (c *CurrencyService) AddCurrency(currency domain.Currency) (domain.Currency, error) {
	return c.repository.AddCurrency(currency)
}

func (c *CurrencyService) GetCurrencyByCode(code string) (domain.Currency, error) {
	return c.repository.GetCurrencyByCode(code)
}
