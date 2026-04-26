package services

import (
	"currency-exchange/internal/domain"
	"currency-exchange/internal/dto"
	"math/big"
)

type ExchangeService struct {
	exchangeRateRepository ExchangeRateRepository
	currencyRepository     CurrencyRepository
}

func ExchangeServiceNew(
	exchangeRateRepository ExchangeRateRepository,
	currencyRepository CurrencyRepository,
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

	targetCurrency, err := e.currencyRepository.GetCurrencyByCode(targetCurrencyCode)
	if err != nil {
		return dto.CurencyExchangeResponse{}, err
	}

	baseDTO := dto.CurrencyResponse{
		ID:   baseCurrency.ID,
		Code: baseCurrency.Code,
		Name: baseCurrency.Name,
		Sign: baseCurrency.Sign,
	}

	targetDTO := dto.CurrencyResponse{
		ID:   targetCurrency.ID,
		Code: targetCurrency.Code,
		Name: targetCurrency.Name,
		Sign: targetCurrency.Sign,
	}

	if rate, found, err := e.exchangeRateRepository.GetExchangeRate(baseCurrency.ID, targetCurrency.ID); err == nil && found {

		converted := new(big.Rat).Mul(amount, rate.Rate)

		return dto.CurencyExchangeResponse{
			BaseCurrency:    baseDTO,
			TargetCurrency:  targetDTO,
			Rate:            rate.Rate.FloatString(2),
			Amount:          amount.FloatString(2),
			ConvertedAmount: converted.FloatString(2),
		}, nil
	}

	if rate, found, err := e.exchangeRateRepository.GetExchangeRate(targetCurrency.ID, baseCurrency.ID); err == nil && found {

		one := big.NewRat(1, 1)
		inverted := new(big.Rat).Quo(one, rate.Rate)

		converted := new(big.Rat).Mul(amount, inverted)

		return dto.CurencyExchangeResponse{
			BaseCurrency:    baseDTO,
			TargetCurrency:  targetDTO,
			Rate:            inverted.FloatString(2),
			Amount:          amount.FloatString(2),
			ConvertedAmount: converted.FloatString(2),
		}, nil
	}

	usd, err := e.currencyRepository.GetCurrencyByCode("USD")
	if err != nil {
		return dto.CurencyExchangeResponse{}, err
	}

	baseToUSD, baseFound, err := e.exchangeRateRepository.GetExchangeRate(usd.ID, baseCurrency.ID)
	if err != nil {
		return dto.CurencyExchangeResponse{}, err
	}

	targetToUSD, targetFound, err := e.exchangeRateRepository.GetExchangeRate(usd.ID, targetCurrency.ID)
	if err != nil {
		return dto.CurencyExchangeResponse{}, err
	}

	if baseFound && targetFound {

		crossRate := new(big.Rat).Quo(targetToUSD.Rate, baseToUSD.Rate)
		converted := new(big.Rat).Mul(amount, crossRate)

		return dto.CurencyExchangeResponse{
			BaseCurrency:    baseDTO,
			TargetCurrency:  targetDTO,
			Rate:            crossRate.FloatString(2),
			Amount:          amount.FloatString(2),
			ConvertedAmount: converted.FloatString(2),
		}, nil
	}

	return dto.CurencyExchangeResponse{}, domain.ErrExchangeRateNotFound
}
