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

func (c *CurrencyController) HandleCurrencies(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		c.getCurrencies(w)
	case "POST":
		c.addCurrency(w, r)
	default:
		http.Error(w, "This method is not allowed", http.StatusMethodNotAllowed)
	}
}

func (c *CurrencyController) getCurrencies(w http.ResponseWriter) {
	currencies, err := c.service.GetCurrencies()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(currencies); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *CurrencyController) addCurrency(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	name := r.FormValue("name")
	code := r.FormValue("code")
	sign := r.FormValue("sign")

	req := dto.CurrencyRequest{
		Name: name,
		Code: code,
		Sign: sign,
	}

	if err := dto.ValidateCurrencyFields(req); err != nil {
		utilities.HandleError(w, err)
		return
	}

	currency := domain.Currency{
		Name: req.Name,
		Code: req.Code,
		Sign: req.Sign,
	}

	res, err := c.service.AddCurrency(currency)

	if err != nil {
		utilities.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *CurrencyController) GetCurrency(w http.ResponseWriter, r *http.Request) {
	currencyCode, err := utilities.GetLastPathSegment(r.URL.Path)
	if err != nil {
		utilities.HandleError(w, err)
		return
	}

	currency, err := c.service.GetCurrencyByCode(currencyCode)

	if err != nil {
		utilities.HandleError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(currency); err != nil {
		utilities.HandleError(w, err)
		return
	}
}
