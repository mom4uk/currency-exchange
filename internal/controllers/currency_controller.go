package controllers

import (
	"currency-exchange/internal/apierrors"
	"currency-exchange/internal/domain"
	"currency-exchange/internal/dto"
	"currency-exchange/internal/httputil"
	"currency-exchange/internal/services"
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
		apierrors.HandleError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(currencies); err != nil {
		apierrors.WriteError(w, "Json convertation error", http.StatusInternalServerError)
		return
	}
}

func (c *CurrencyController) AddCurrency(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		apierrors.WriteError(w, "Parse form error", http.StatusBadRequest)
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
		apierrors.HandleError(w, err)
		return
	}

	currency := domain.Currency{
		Name: req.Name,
		Code: req.Code,
		Sign: req.Sign,
	}

	res, err := c.service.AddCurrency(currency)

	if err != nil {
		apierrors.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		apierrors.WriteError(w, "Json convertation error", http.StatusInternalServerError)
		return
	}
}

func (c *CurrencyController) GetCurrency(w http.ResponseWriter, r *http.Request) {
	currencyCode, err := httputil.GetLastPathSegment(r.URL.Path)
	if err != nil {
		apierrors.HandleError(w, err)
		return
	}

	currency, err := c.service.GetCurrencyByCode(currencyCode)

	if err != nil {
		apierrors.HandleError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(currency); err != nil {
		apierrors.WriteError(w, "Json convertation error", http.StatusInternalServerError)
		return
	}
}
