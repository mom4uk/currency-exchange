package controllers

import (
	"currency-exchange/internal/domain"
	"currency-exchange/internal/repositories"
	"currency-exchange/internal/services"
	"currency-exchange/internal/utilities"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
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
		c.getExchangeRates(w)
	case "POST":
		c.addExchangeRates(w, r)
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
	r.ParseForm()

	currency := domain.Currency{
		Name: r.FormValue("name"),
		Code: r.FormValue("code"),
		Sign: r.FormValue("sign"),
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
	currency, err := c.repository.GetCurrencyByCode(currencyCode)

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

func (c *CurrencyController) addExchangeRates(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	rate, err := strconv.ParseFloat(r.FormValue("rate"), 64)

	if err != nil {
		http.Error(w, "Invalid rate", http.StatusInternalServerError)
		return
	}

	exchangeRate := domain.AddExchangeRateRequest{
		BaseCurrencyCode:   r.FormValue("baseCurrencyCode"),
		TargetCurrencyCode: r.FormValue("targetCurrencyCode"),
		Rate:               rate,
	}

	res, err := c.repository.AddExchangeRates(exchangeRate)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(res)
}

func (c *CurrencyController) GetExchange(w http.ResponseWriter, r *http.Request) {

}

func (c *CurrencyController) getExchangeRates(w http.ResponseWriter) {
	exchangeRates, err := c.repository.GetExchangeRates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(exchangeRates)
}

func (c *CurrencyController) getExchangeRate(w http.ResponseWriter, r *http.Request) {

}

func (c *CurrencyController) patchExchangeRate(w http.ResponseWriter, r *http.Request) {

}
