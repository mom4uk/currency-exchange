package services

import (
	"currency-exchange/internal/dto"
	"currency-exchange/internal/repositories"
	"math/big"
	"strconv"
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
	amount *big.Rat,
) (dto.CurencyExchangeResponse, error) {

	baseCurrency, err := e.currencyRepository.GetCurrencyByCode(baseCurrencyCode)
	if err != nil {
		return dto.CurencyExchangeResponse{}, err
	}
	baseCurrencyResponce := dto.CurrencyResponse{
		ID:   strconv.Itoa(baseCurrency.ID),
		Code: baseCurrency.Code,
		Name: baseCurrency.Name,
		Sign: baseCurrency.Sign,
	}

	targetCurrency, err := e.currencyRepository.GetCurrencyByCode(targetCurrencyCode)
	if err != nil {
		return dto.CurencyExchangeResponse{}, err
	}
	targetCurrencyResponce := dto.CurrencyResponse{
		ID:   strconv.Itoa(targetCurrency.ID),
		Code: targetCurrency.Code,
		Name: targetCurrency.Name,
		Sign: targetCurrency.Sign,
	}

	rate, found, err := e.exchangeRateRepository.GetExchangeRate(baseCurrency.ID, targetCurrency.ID)
	if err != nil {
		return dto.CurencyExchangeResponse{}, err
	}
	if found {
		return dto.CurencyExchangeResponse{
			BaseCurrency:    baseCurrencyResponce,
			TargetCurrency:  targetCurrencyResponce,
			Rate:            rate.Rate.FloatString(2),
			Amount:          amount.FloatString(2),
			ConvertedAmount: new(big.Rat).Mul(amount, rate.Rate).FloatString(2),
		}, nil
	}

	rate, found, err = e.exchangeRateRepository.GetExchangeRate(targetCurrency.ID, baseCurrency.ID)
	if err != nil {
		return dto.CurencyExchangeResponse{}, err
	}
	if found {
		one := big.NewRat(1, 1)

		invertedRate := new(big.Rat).Quo(one, rate.Rate)

		return dto.CurencyExchangeResponse{
			BaseCurrency:    baseCurrencyResponce,
			TargetCurrency:  targetCurrencyResponce,
			Rate:            invertedRate.FloatString(2),
			Amount:          amount.FloatString(2),
			ConvertedAmount: new(big.Rat).Mul(amount, invertedRate).FloatString(2),
		}, nil
	}

	usdCurrency, err := e.currencyRepository.GetCurrencyByCode("USD")
	if err != nil {
		return dto.CurencyExchangeResponse{}, err
	}

	baseToUsd, baseFound, err := e.exchangeRateRepository.GetExchangeRate(usdCurrency.ID, baseCurrency.ID)
	if err != nil {
		return dto.CurencyExchangeResponse{}, err
	}

	targetToUsd, targetFound, err := e.exchangeRateRepository.GetExchangeRate(usdCurrency.ID, targetCurrency.ID)
	if err != nil {
		return dto.CurencyExchangeResponse{}, err
	}

	if baseFound && targetFound {
		crossRate := new(big.Rat).Quo(targetToUsd.Rate, baseToUsd.Rate)
		converted := new(big.Rat).Mul(amount, crossRate)

		return dto.CurencyExchangeResponse{
			BaseCurrency:    baseCurrencyResponce,
			TargetCurrency:  targetCurrencyResponce,
			Rate:            crossRate.FloatString(2),
			Amount:          amount.FloatString(2),
			ConvertedAmount: converted.FloatString(2),
		}, nil
	}

	return dto.CurencyExchangeResponse{}, err
}
