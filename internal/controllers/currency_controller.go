package controllers

import (
	"currency-exchange/internal/repositories"
	"currency-exchange/internal/services"
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
		c.getCurrencies(w, r)
	case "POST":
		c.postCurrencies(w, r)
	default:
		http.Error(w, "This method is not allowed", 405)
	}
}

func (c *CurrencyController) HandleExchangeRates(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		c.getExchangeRates(w, r)
	case "POST":
		c.postExchangeRates(w, r)
	default:
		http.Error(w, "This method is not allowed", 405)
	}
}

func (c *CurrencyController) HandleExchangeRate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		c.getExchangeRate(w, r)
	case "PATCH":
		c.patchExchangeRate(w, r)
	default:
		http.Error(w, "This method is not allowed", 405)
	}
}

func (c *CurrencyController) GetCurrency(w http.ResponseWriter, r *http.Request) {

}

func (c *CurrencyController) GetExchange(w http.ResponseWriter, r *http.Request) {

}

func (c *CurrencyController) getCurrencies(w http.ResponseWriter, r *http.Request) {
	// some logic
}

func (c *CurrencyController) postCurrencies(w http.ResponseWriter, r *http.Request) {
	// some logic
}

func (c *CurrencyController) getExchangeRates(w http.ResponseWriter, r *http.Request) {

}

func (c *CurrencyController) postExchangeRates(w http.ResponseWriter, r *http.Request) {

}

func (c *CurrencyController) getExchangeRate(w http.ResponseWriter, r *http.Request) {

}

func (c *CurrencyController) patchExchangeRate(w http.ResponseWriter, r *http.Request) {

}
