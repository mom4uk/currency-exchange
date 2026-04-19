package dto

import "currency-exchange/internal/domain"

type CurrencyResponse struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Sign string `json:"sign"`
}

type CurrencyRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Sign string `json:"sign"`
}

func ValidateFields(req CurrencyRequest) error {
	if req.Code == "" || req.Name == "" || req.Sign == "" {
		return domain.ErrAbsenceOfField
	}
	return nil
}
