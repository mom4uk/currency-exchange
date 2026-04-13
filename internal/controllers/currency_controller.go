package controllers

import (
	"currency-exchange/internal/domain"
	"currency-exchange/internal/repositories"
	"currency-exchange/internal/services"
	"currency-exchange/internal/utilities"
	"database/sql"
	"encoding/json"
	"net/http"
)

type CurrencyController struct {
	repository repositories.CurrencyRepository
	service    services.CurrencyService
}

func NewController(repo *repositories.CurrencyRepository, service *services.CurrencyService) *CurrencyController {
	return &CurrencyController{
		repository: *repo,
		service:    *service,
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

func (c *CurrencyController) HandleExchangeRates(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		c.getExchangeRates(w, r)
	case "POST":
		c.postExchangeRates(w, r)
	default:
		http.Error(w, "This method is not allowed", http.StatusMethodNotAllowed)
	}
}

func (c *CurrencyController) HandleExchangeRate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		c.getExchangeRate(w, r)
	case "PATCH":
		c.patchExchangeRate(w, r)
	default:
		http.Error(w, "This method is not allowed", http.StatusMethodNotAllowed)
	}
}

func (c *CurrencyController) getCurrencies(w http.ResponseWriter) {
	currencies, err := c.repository.GetCurrencies()

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
	var currency domain.Currency

	if err := json.NewDecoder(r.Body).Decode(&currency); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := c.repository.AddCurrency(currency)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(currency); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *CurrencyController) GetCurrency(w http.ResponseWriter, r *http.Request) {
	currencyCode := utilities.GetCurrencyCode(r.URL.Path)
	currency, err := c.repository.GetCurrency(currencyCode)

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

func (c *CurrencyController) GetExchange(w http.ResponseWriter, r *http.Request) {

}

func (c *CurrencyController) getExchangeRates(w http.ResponseWriter, r *http.Request) {

}

func (c *CurrencyController) postExchangeRates(w http.ResponseWriter, r *http.Request) {

}

func (c *CurrencyController) getExchangeRate(w http.ResponseWriter, r *http.Request) {

}

func (c *CurrencyController) patchExchangeRate(w http.ResponseWriter, r *http.Request) {

}
