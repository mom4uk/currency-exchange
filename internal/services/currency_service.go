package services

import (
	"currency-exchange/internal/repositories"
)

type CurrencyService struct {
	repository repositories.CurrencyRepository
}

func CurrencyServiceNew(repo *repositories.CurrencyRepository) *CurrencyService {
	return &CurrencyService{
		repository: *repo,
	}
}
