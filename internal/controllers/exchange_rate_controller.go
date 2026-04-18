package controllers

import (
	"currency-exchange/internal/domain"
	"currency-exchange/internal/services"
	"currency-exchange/internal/utilities"
	"encoding/json"
	"net/http"
	"strconv"
)

type ExchangeRateController struct {
	service *services.ExchangeRateService
}

func NewExchangeRateController(service *services.ExchangeRateService) *ExchangeRateController {
	return &ExchangeRateController{
		service: service,
	}
}

func (e *ExchangeRateController) HandleExchangeRates(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		e.getExchangeRates(w)
	case "POST":
		e.addExchangeRates(w, r)
	default:
		http.Error(w, "This method is not allowed", http.StatusMethodNotAllowed)
	}
}

func (e *ExchangeRateController) HandleExchangeRate(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		e.getExchangeRate(w, r)
	case "PATCH":
		e.patchExchangeRate(w, r)
	default:
		http.Error(w, "This method is not allowed", http.StatusMethodNotAllowed)
	}
}

func (e *ExchangeRateController) addExchangeRates(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	rate, err := strconv.ParseFloat(r.FormValue("rate"), 64)

	if err != nil {
		http.Error(w, "Invalid rate", http.StatusInternalServerError)
		return
	}

	req := domain.AddExchangeRateRequest{
		BaseCurrencyCode:   r.FormValue("baseCurrencyCode"),
		TargetCurrencyCode: r.FormValue("targetCurrencyCode"),
		Rate:               rate,
	}

	exchangeRate, err := e.service.AddExchangeRates(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := e.service.GetExchangeRateResponse(exchangeRate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(res)
}

func (e *ExchangeRateController) getExchangeRates(w http.ResponseWriter) {
	rates, err := e.service.GetExchangeRates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := e.service.GetExchangeRatesResponse(rates)

	json.NewEncoder(w).Encode(response)
}

func (e *ExchangeRateController) getExchangeRate(w http.ResponseWriter, r *http.Request) {
	baseCurrencyCode, targetCurrencyCode, err := utilities.GetCurrencyCodes(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rate, err := e.service.GetExchangeRateByCodes(baseCurrencyCode, targetCurrencyCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := e.service.GetExchangeRateResponse(rate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(response)
}

func (e *ExchangeRateController) patchExchangeRate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	baseCurrencyCode, targetCurrencyCode, err := utilities.GetCurrencyCodes(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rateValue, err := strconv.ParseFloat(r.FormValue("rate"), 64)

	rate, err := e.service.UpdateExchangeRate(baseCurrencyCode, targetCurrencyCode, rateValue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(rate)
}
