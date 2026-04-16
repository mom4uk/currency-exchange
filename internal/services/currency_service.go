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

func (c *CurrencyService) GetExchange(
	baseCurrencyCode string,
	targetCurrencyCode string,
	amount float64,
) (domain.CurencyExchange, error) {

	baseCurrency, err := c.currencyRepository.GetCurrencyByCode(baseCurrencyCode)
	if err != nil {
		return domain.CurencyExchange{}, err
	}

	targetCurrency, err := c.currencyRepository.GetCurrencyByCode(targetCurrencyCode)
	if err != nil {
		return domain.CurencyExchange{}, err
	}

	rate, found, err := c.exchangeRateRepository.GetExchangeRate(baseCurrency.ID, targetCurrency.ID)
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

	rate, found, err = c.exchangeRateRepository.GetExchangeRate(targetCurrency.ID, baseCurrency.ID)
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

	usdCurrency, err := c.currencyRepository.GetCurrencyByCode("USD")
	if err != nil {
		return domain.CurencyExchange{}, err
	}

	baseToUsd, baseFound, err := c.exchangeRateRepository.GetExchangeRate(usdCurrency.ID, baseCurrency.ID)
	if err != nil {
		return domain.CurencyExchange{}, err
	}

	targetToUsd, targetFound, err := c.exchangeRateRepository.GetExchangeRate(usdCurrency.ID, targetCurrency.ID)
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
