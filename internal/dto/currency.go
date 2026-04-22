package dto

import "currency-exchange/internal/domain"

type CurrencyResponse struct {
	ID   string `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
	Sign string `json:"sign"`
}

type CurrencyRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Sign string `json:"sign"`
}

func ValidateCurrencyFields(req CurrencyRequest) error {
	if req.Code == "" || req.Name == "" || req.Sign == "" {
		return domain.ErrAbsenceOfCurrencyField
	}
	return nil
}
