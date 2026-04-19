package services

import (
	"currency-exchange/internal/domain"
	"currency-exchange/internal/dto"
	"currency-exchange/internal/repositories"
)

type ExchangeRateService struct {
	exchangeRateRepository *repositories.ExchangeRateRepository
	currencyRepository     *repositories.CurrencyRepository
}

func ExchangeRateServiceNew(
	exchangeRateRepository *repositories.ExchangeRateRepository,
	currencyRepository *repositories.CurrencyRepository,
) *ExchangeRateService {
	return &ExchangeRateService{
		exchangeRateRepository: exchangeRateRepository,
		currencyRepository:     currencyRepository,
	}
}

func (s *ExchangeRateService) UpdateExchangeRate(baseCurrencyCode string, targetCurrencyCode string, rate float64) (dto.ExchangeRateResponse, error) {
	baseCurrency, err := s.currencyRepository.GetCurrencyByCode(baseCurrencyCode)
	if err != nil {
		return dto.ExchangeRateResponse{}, err
	}

	targetCurrency, err := s.currencyRepository.GetCurrencyByCode(targetCurrencyCode)
	if err != nil {
		return dto.ExchangeRateResponse{}, err
	}

	result, err := s.exchangeRateRepository.UpdateExchangeRate(baseCurrency, targetCurrency, rate)
	if err != nil {
		return dto.ExchangeRateResponse{}, err
	}

	return dto.ExchangeRateResponse{
		ID:             result.ID,
		BaseCurrency:   baseCurrency,
		TargetCurrency: targetCurrency,
		Rate:           result.Rate,
	}, nil
}

func (s *ExchangeRateService) GetExchangeRateByCodes(baseCurrencyCode string, targetCurrencyCode string) (domain.ExchangeRate, error) {
	baseCurrency, err := s.currencyRepository.GetCurrencyByCode(baseCurrencyCode)
	if err != nil {
		return domain.ExchangeRate{}, err
	}

	targetCurrency, err := s.currencyRepository.GetCurrencyByCode(targetCurrencyCode)
	if err != nil {
		return domain.ExchangeRate{}, err
	}

	return s.exchangeRateRepository.GetExchangeRateByCodes(baseCurrency, targetCurrency)
}

func (s *ExchangeRateService) GetExchangeRates() ([]domain.ExchangeRate, error) {
	return s.exchangeRateRepository.GetExchangeRates()
}

func (s *ExchangeRateService) AddExchangeRates(req dto.ExchangeRate) (domain.ExchangeRate, error) {
	baseCurrency, err := s.currencyRepository.GetCurrencyByCode(req.BaseCurrencyCode)
	if err != nil {
		return domain.ExchangeRate{}, err
	}

	targetCurrency, err := s.currencyRepository.GetCurrencyByCode(req.TargetCurrencyCode)
	if err != nil {
		return domain.ExchangeRate{}, err
	}

	return s.exchangeRateRepository.AddExchangeRates(baseCurrency, targetCurrency, req.Rate)
}

func (s *ExchangeRateService) GetExchangeRateResponse(rate domain.ExchangeRate) (dto.ExchangeRateResponse, error) {
	baseCurrency, err := s.currencyRepository.GetCurrencyById(rate.BaseCurrencyId)
	if err != nil {
		return dto.ExchangeRateResponse{}, err
	}

	targetCurrency, err := s.currencyRepository.GetCurrencyById(rate.TargetCurrencyId)
	if err != nil {
		return dto.ExchangeRateResponse{}, err
	}

	return dto.ExchangeRateResponse{
		ID:             rate.ID,
		BaseCurrency:   baseCurrency,
		TargetCurrency: targetCurrency,
		Rate:           rate.Rate,
	}, err
}

func (s *ExchangeRateService) GetExchangeRatesResponse(rates []domain.ExchangeRate) ([]dto.ExchangeRateResponse, error) {
	var result []dto.ExchangeRateResponse

	for _, rate := range rates {
		response, err := s.GetExchangeRateResponse(rate)
		if err != nil {
			return nil, err
		}
		result = append(result, response)
	}
	return result, nil
}
