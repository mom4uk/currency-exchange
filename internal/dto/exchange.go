package dto

import (
	"currency-exchange/internal/domain"
	"strings"
)

type CurencyExchangeResponse struct {
	BaseCurrency    CurrencyResponse `json:"baseCurrency"`
	TargetCurrency  CurrencyResponse `json:"targetCurrency"`
	Rate            string           `json:"rate"`
	Amount          string           `json:"amount"`
	ConvertedAmount string           `json:"convertedAmount"`
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

	return nil
}
