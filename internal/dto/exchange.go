package dto

import (
	"currency-exchange/internal/domain"
	"strconv"
	"strings"
)

type CurencyExchangeResponse struct {
	BaseCurrency    CurrencyResponse `json:"baseCurrency"`
	TargetCurrency  CurrencyResponse `json:"targetCurrency"`
	Rate            float64          `json:"rate"`
	Amount          float64          `json:"amount"`
	ConvertedAmount float64          `json:"convertedAmount"`
}

func ValidateExchangeFields(from, to, amount string) error {
	from = strings.TrimSpace(from)
	to = strings.TrimSpace(to)
	amount = strings.TrimSpace(amount)

	if from == "" {
		return domain.ErrMissingFromCurrency
	}

	if to == "" {
		return domain.ErrMissingToCurrency
	}

	if amount == "" {
		return domain.ErrMissingAmount
	}

	if _, err := strconv.ParseFloat(amount, 64); err != nil {
		return domain.ErrAmountFormatIncorrect
	}

	return nil
}
