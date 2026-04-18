package controllers

import (
	"currency-exchange/internal/domain"
	"currency-exchange/internal/services"
	"currency-exchange/internal/utilities"
	"database/sql"
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

	currency := domain.Currency{
		Name: r.FormValue("name"),
		Code: r.FormValue("code"),
		Sign: r.FormValue("sign"),
	}

	res, err := c.service.AddCurrency(currency)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *CurrencyController) GetCurrency(w http.ResponseWriter, r *http.Request) {
	currencyCode := utilities.GetLastPathSegment(r.URL.Path)
	currency, err := c.service.GetCurrencyByCode(currencyCode)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(currency); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
