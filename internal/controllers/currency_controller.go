package controllers

import (
	"currency-exchange/internal/domain"
	"currency-exchange/internal/dto"
	"currency-exchange/internal/services"
	"currency-exchange/internal/utilities"
	"encoding/json"
	"net/http"
)

type CurrencyController struct {
	service *services.CurrencyService
}

func NewController(service *services.CurrencyService) *CurrencyController {
	return &CurrencyController{
		service: service,
	}
}

func (c *CurrencyController) GetCurrencies(w http.ResponseWriter, r *http.Request) {
	currencies, err := c.service.GetCurrencies()

	if err != nil {
		HandleError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(currencies); err != nil {
		WriteError(w, "Json convertation error", http.StatusInternalServerError)
		return
	}
}

func (c *CurrencyController) AddCurrency(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		WriteError(w, "Parse form error", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	code := r.FormValue("code")
	sign := r.FormValue("sign")

	req := dto.CurrencyRequest{
		Name: name,
		Code: code,
		Sign: sign,
	}

	if err := dto.ValidateCurrencyFields(req); err != nil {
		HandleError(w, err)
		return
	}

	currency := domain.Currency{
		Name: req.Name,
		Code: req.Code,
		Sign: req.Sign,
	}

	res, err := c.service.AddCurrency(currency)

	if err != nil {
		HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		WriteError(w, "Json convertation error", http.StatusInternalServerError)
		return
	}
}

func (c *CurrencyController) GetCurrency(w http.ResponseWriter, r *http.Request) {
	currencyCode, err := utilities.GetLastPathSegment(r.URL.Path)
	if err != nil {
		HandleError(w, err)
		return
	}

	currency, err := c.service.GetCurrencyByCode(currencyCode)

	if err != nil {
		HandleError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(currency); err != nil {
		WriteError(w, "Json convertation error", http.StatusInternalServerError)
		return
	}
}
