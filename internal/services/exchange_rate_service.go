package services

import (
	"currency-exchange/internal/domain"
	"currency-exchange/internal/repositories"
)

type ExchangeRateService struct {
	repository repositories.CurrencyRepository
}

func ExchangeRateServiceNew(repo *repositories.CurrencyRepository) *CurrencyService {
	return &CurrencyService{
		repository: *repo,
	}
}

func (s *ExchangeRateService) GetExchangeRateResponse(rate domain.ExchangeRate) (domain.ExchangeRateResponse, error) {
	baseCurrency, err := s.repository.GetCurrencyById(rate.BaseCurrencyId)
	if err != nil {
		return domain.ExchangeRateResponse{}, err
	}

	targetCurrency, err := s.repository.GetCurrencyById(rate.TargetCurrencyId)
	if err != nil {
		return domain.ExchangeRateResponse{}, err
	}

	return domain.ExchangeRateResponse{
		ID:             rate.ID,
		BaseCurrency:   baseCurrency,
		TargetCurrency: targetCurrency,
		Rate:           rate.Rate,
	}, err
}

func (s *ExchangeRateService) GetExchangeRatesResponse(rates []domain.ExchangeRate) ([]domain.ExchangeRateResponse, error) {
	var result []domain.ExchangeRateResponse

	for _, rate := range rates {
		response, err := s.GetExchangeRateResponse(rate)
		if err != nil {
			return nil, err
		}
		result = append(result, response)
	}
	return result, nil
}
