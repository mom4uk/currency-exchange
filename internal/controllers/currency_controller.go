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
	service services.CurrencyService
}

func NewController(service *services.CurrencyService) *CurrencyController {
	return &CurrencyController{
		service: *service,
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

// func (c *CurrencyController) HandleExchangeRates(w http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case "GET":
// 		c.getExchangeRates(w)
// 	case "POST":
// 		c.addExchangeRates(w, r)
// 	default:
// 		http.Error(w, "This method is not allowed", http.StatusMethodNotAllowed)
// 	}
// }

// func (c *CurrencyController) HandleExchangeRate(w http.ResponseWriter, r *http.Request) {
// 	switch r.Method {
// 	case "GET":
// 		c.getExchangeRate(w, r)
// 	case "PATCH":
// 		c.patchExchangeRate(w, r)
// 	default:
// 		http.Error(w, "This method is not allowed", http.StatusMethodNotAllowed)
// 	}
// }

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

// func (c *CurrencyController) addExchangeRates(w http.ResponseWriter, r *http.Request) {
// 	r.ParseForm()

// 	rate, err := strconv.ParseFloat(r.FormValue("rate"), 64)

// 	if err != nil {
// 		http.Error(w, "Invalid rate", http.StatusInternalServerError)
// 		return
// 	}

// 	req := domain.AddExchangeRateRequest{
// 		BaseCurrencyCode:   r.FormValue("baseCurrencyCode"),
// 		TargetCurrencyCode: r.FormValue("targetCurrencyCode"),
// 		Rate:               rate,
// 	}

// 	exchangeRate, err := c.repository.AddExchangeRates(req)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	res, err := c.exchangeRateService.GetExchangeRateResponse(exchangeRate)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	json.NewEncoder(w).Encode(res)
// }

// func (c *CurrencyController) getExchangeRates(w http.ResponseWriter) {
// 	rates, err := c.repository.GetExchangeRates()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	response, err := c.exchangeRateService.GetExchangeRatesResponse(rates)

// 	json.NewEncoder(w).Encode(response)
// }

// func (c *CurrencyController) getExchangeRate(w http.ResponseWriter, r *http.Request) {
// 	baseCurrencyCode, targetCurrencyCode, err := utilities.GetCurrencyCodes(r.URL.Path)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	rate, err := c.repository.GetExchangeRatesByCodes(baseCurrencyCode, targetCurrencyCode)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	response, err := c.exchangeRateService.GetExchangeRateResponse(rate)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	json.NewEncoder(w).Encode(response)
// }

// func (c *CurrencyController) patchExchangeRate(w http.ResponseWriter, r *http.Request) {
// 	r.ParseForm()
// 	baseCurrencyCode, targetCurrencyCode, err := utilities.GetCurrencyCodes(r.URL.Path)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	rateValue, err := strconv.ParseFloat(r.FormValue("rate"), 64)

// 	rate, err := c.repository.UpdateExchangeRate(baseCurrencyCode, targetCurrencyCode, rateValue)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	response, err := c.exchangeRateService.GetExchangeRateResponse(rate)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	json.NewEncoder(w).Encode(response)
// }

func (c *CurrencyController) GetExchange(w http.ResponseWriter, r *http.Request) {

}
