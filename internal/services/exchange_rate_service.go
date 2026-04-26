package services

import (
	"currency-exchange/internal/domain"
	"currency-exchange/internal/dto"
	"currency-exchange/internal/utilities"
	"math/big"
)

type ExchangeRateService struct {
	exchangeRateRepository ExchangeRateRepository
	currencyRepository     CurrencyRepository
}

func ExchangeRateServiceNew(
	exchangeRateRepository ExchangeRateRepository,
	currencyRepository CurrencyRepository,
) *ExchangeRateService {
	return &ExchangeRateService{
		exchangeRateRepository: exchangeRateRepository,
		currencyRepository:     currencyRepository,
	}
}

func (s *ExchangeRateService) UpdateExchangeRate(baseCurrencyCode string, targetCurrencyCode string, rate *big.Rat) (dto.ExchangeRateResponse, error) {
	baseCurrency, err := s.currencyRepository.GetCurrencyByCode(baseCurrencyCode)
	if err != nil {
		return dto.ExchangeRateResponse{}, err
	}
	baseCurrencyResponce := dto.CurrencyResponse{
		ID:   baseCurrency.ID,
		Code: baseCurrency.Code,
		Name: baseCurrency.Name,
		Sign: baseCurrency.Sign,
	}

	targetCurrency, err := s.currencyRepository.GetCurrencyByCode(targetCurrencyCode)
	if err != nil {
		return dto.ExchangeRateResponse{}, err
	}
	targetCurrencyResponce := dto.CurrencyResponse{
		ID:   targetCurrency.ID,
		Code: targetCurrency.Code,
		Name: targetCurrency.Name,
		Sign: targetCurrency.Sign,
	}

	result, err := s.exchangeRateRepository.UpdateExchangeRate(baseCurrency, targetCurrency, rate)
	if err != nil {
		return dto.ExchangeRateResponse{}, err
	}

	return dto.ExchangeRateResponse{
		ID:             result.ID,
		BaseCurrency:   baseCurrencyResponce,
		TargetCurrency: targetCurrencyResponce,
		Rate:           utilities.RatToFloat(result.Rate),
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

	rate, found, err := s.exchangeRateRepository.GetExchangeRate(baseCurrency.ID, targetCurrency.ID)
	if err != nil {
		return domain.ExchangeRate{}, err
	}
	if !found {
		return domain.ExchangeRate{}, domain.ErrExchangeRateNotFound
	}

	return rate, nil
}

func (s *ExchangeRateService) GetExchangeRates() ([]domain.ExchangeRate, error) {
	return s.exchangeRateRepository.GetExchangeRates()
}

func (s *ExchangeRateService) AddExchangeRates(req domain.Exchange) (domain.ExchangeRate, error) {
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
	baseCurrencyResponce := dto.CurrencyResponse{
		ID:   baseCurrency.ID,
		Code: baseCurrency.Code,
		Name: baseCurrency.Name,
		Sign: baseCurrency.Sign,
	}

	targetCurrency, err := s.currencyRepository.GetCurrencyById(rate.TargetCurrencyId)
	if err != nil {
		return dto.ExchangeRateResponse{}, err
	}
	targetCurrencyResponce := dto.CurrencyResponse{
		ID:   targetCurrency.ID,
		Code: targetCurrency.Code,
		Name: targetCurrency.Name,
		Sign: targetCurrency.Sign,
	}

	return dto.ExchangeRateResponse{
		ID:             rate.ID,
		BaseCurrency:   baseCurrencyResponce,
		TargetCurrency: targetCurrencyResponce,
		Rate:           utilities.RatToFloat(rate.Rate),
	}, err
}

func (s *ExchangeRateService) GetExchangeRatesResponse(rates []domain.ExchangeRate) ([]dto.ExchangeRateResponse, error) {
	result := []dto.ExchangeRateResponse{}

	for _, rate := range rates {
		response, err := s.GetExchangeRateResponse(rate)
		if err != nil {
			return nil, err
		}
		result = append(result, response)
	}
	return result, nil
}
