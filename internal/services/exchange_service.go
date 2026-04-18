package services

import (
	"currency-exchange/internal/domain"
	"currency-exchange/internal/repositories"
)

type ExchangeService struct {
	exchangeRateRepository *repositories.ExchangeRateRepository
	currencyRepository     *repositories.CurrencyRepository
}

func ExchangeServiceNew(
	exchangeRateRepository *repositories.ExchangeRateRepository,
	currencyRepository *repositories.CurrencyRepository,
) *ExchangeService {
	return &ExchangeService{
		exchangeRateRepository: exchangeRateRepository,
		currencyRepository:     currencyRepository,
	}
}

func (e *ExchangeService) GetExchange(
	baseCurrencyCode string,
	targetCurrencyCode string,
	amount float64,
) (domain.CurencyExchange, error) {

	baseCurrency, err := e.currencyRepository.GetCurrencyByCode(baseCurrencyCode)
	if err != nil {
		return domain.CurencyExchange{}, err
	}

	targetCurrency, err := e.currencyRepository.GetCurrencyByCode(targetCurrencyCode)
	if err != nil {
		return domain.CurencyExchange{}, err
	}

	rate, found, err := e.exchangeRateRepository.GetExchangeRate(baseCurrency.ID, targetCurrency.ID)
	if err != nil {
		return domain.CurencyExchange{}, err
	}
	if found {
		return domain.CurencyExchange{
			BaseCurrency:    baseCurrency,
			TargetCurrency:  targetCurrency,
			Rate:            rate.Rate,
			Amount:          amount,
			ConvertedAmount: amount * rate.Rate,
		}, nil
	}

	rate, found, err = e.exchangeRateRepository.GetExchangeRate(targetCurrency.ID, baseCurrency.ID)
	if err != nil {
		return domain.CurencyExchange{}, err
	}
	if found {
		invertedRate := 1 / rate.Rate

		return domain.CurencyExchange{
			BaseCurrency:    baseCurrency,
			TargetCurrency:  targetCurrency,
			Rate:            invertedRate,
			Amount:          amount,
			ConvertedAmount: amount * invertedRate,
		}, nil
	}

	usdCurrency, err := e.currencyRepository.GetCurrencyByCode("USD")
	if err != nil {
		return domain.CurencyExchange{}, err
	}

	baseToUsd, baseFound, err := e.exchangeRateRepository.GetExchangeRate(usdCurrency.ID, baseCurrency.ID)
	if err != nil {
		return domain.CurencyExchange{}, err
	}

	targetToUsd, targetFound, err := e.exchangeRateRepository.GetExchangeRate(usdCurrency.ID, targetCurrency.ID)
	if err != nil {
		return domain.CurencyExchange{}, err
	}

	if baseFound && targetFound {
		crossRate := targetToUsd.Rate / baseToUsd.Rate

		return domain.CurencyExchange{
			BaseCurrency:    baseCurrency,
			TargetCurrency:  targetCurrency,
			Rate:            crossRate,
			Amount:          amount,
			ConvertedAmount: amount * crossRate,
		}, nil
	}

	return domain.CurencyExchange{}, err
}
