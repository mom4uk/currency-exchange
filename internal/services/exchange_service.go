package services

import (
	"currency-exchange/internal/domain"
	"currency-exchange/internal/dto"
	"currency-exchange/internal/repositories"
	"currency-exchange/internal/utilities"
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

	targetCurrency, err := e.currencyRepository.GetCurrencyByCode(targetCurrencyCode)
	if err != nil {
		return dto.CurencyExchangeResponse{}, err
	}

	baseDTO := dto.CurrencyResponse{
		ID:   strconv.Itoa(baseCurrency.ID),
		Code: baseCurrency.Code,
		Name: baseCurrency.Name,
		Sign: baseCurrency.Sign,
	}

	targetDTO := dto.CurrencyResponse{
		ID:   strconv.Itoa(targetCurrency.ID),
		Code: targetCurrency.Code,
		Name: targetCurrency.Name,
		Sign: targetCurrency.Sign,
	}

	// 1. direct rate
	if rate, found, err := e.exchangeRateRepository.GetExchangeRate(baseCurrency.ID, targetCurrency.ID); err == nil && found {

		converted := new(big.Rat).Mul(amount, rate.Rate)

		return dto.CurencyExchangeResponse{
			BaseCurrency:    baseDTO,
			TargetCurrency:  targetDTO,
			Rate:            utilities.FormatRat(rate.Rate),
			Amount:          utilities.FormatRat(amount),
			ConvertedAmount: utilities.FormatRat(converted),
		}, nil
	}

	// 2. inverse rate
	if rate, found, err := e.exchangeRateRepository.GetExchangeRate(targetCurrency.ID, baseCurrency.ID); err == nil && found {

		one := big.NewRat(1, 1)
		inverted := new(big.Rat).Quo(one, rate.Rate)

		converted := new(big.Rat).Mul(amount, inverted)

		return dto.CurencyExchangeResponse{
			BaseCurrency:    baseDTO,
			TargetCurrency:  targetDTO,
			Rate:            utilities.FormatRat(inverted),
			Amount:          utilities.FormatRat(amount),
			ConvertedAmount: utilities.FormatRat(converted),
		}, nil
	}

	// 3. cross via USD
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
			Rate:            utilities.FormatRat(crossRate),
			Amount:          utilities.FormatRat(amount),
			ConvertedAmount: utilities.FormatRat(converted),
		}, nil
	}

	// 4. not found
	return dto.CurencyExchangeResponse{}, domain.ErrExchangeRateNotFound
}
