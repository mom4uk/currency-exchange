package services

import "currency-exchange/internal/domain"

type CurrencyRepository interface {
	GetCurrencies() ([]domain.Currency, error)
	AddCurrency(currency domain.Currency) (domain.Currency, error)
	GetCurrencyByCode(code string) (domain.Currency, error)
	GetCurrencyById(id int) (domain.Currency, error)
}
