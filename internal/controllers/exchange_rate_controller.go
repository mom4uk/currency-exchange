package controllers

import (
	"currency-exchange/internal/dto"
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

	req := dto.AddExchangeRateRequest{
		BaseCurrencyCode:   r.FormValue("baseCurrencyCode"),
		TargetCurrencyCode: r.FormValue("targetCurrencyCode"),
		Rate:               r.FormValue("rate"),
	}

	if err := dto.ValidateExchangeRateFields(req); err != nil {
		utilities.HandleError(w, err)
		return
	}

	rate, err := strconv.ParseFloat(r.FormValue("rate"), 64)
	if err != nil {
		http.Error(w, "Invalid rate", http.StatusInternalServerError)
		return
	}

	exchangeRate := dto.ExchangeRate{
		BaseCurrencyCode:   req.BaseCurrencyCode,
		TargetCurrencyCode: req.TargetCurrencyCode,
		Rate:               rate,
	}

	result, err := e.service.AddExchangeRates(exchangeRate)
	if err != nil {
		utilities.HandleError(w, err)
		return
	}

	res, err := e.service.GetExchangeRateResponse(result)
	if err != nil {
		utilities.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(res)
}

func (e *ExchangeRateController) getExchangeRates(w http.ResponseWriter) {
	rates, err := e.service.GetExchangeRates()
	if err != nil {
		utilities.HandleError(w, err)
		return
	}

	response, err := e.service.GetExchangeRatesResponse(rates)

	json.NewEncoder(w).Encode(response)
}

func (e *ExchangeRateController) getExchangeRate(w http.ResponseWriter, r *http.Request) {
	baseCurrencyCode, targetCurrencyCode, err := utilities.GetCurrencyCodes(r.URL.Path)
	if err != nil {
		utilities.HandleError(w, err)
		return
	}

	rate, err := e.service.GetExchangeRateByCodes(baseCurrencyCode, targetCurrencyCode)
	if err != nil {
		utilities.HandleError(w, err)
		return
	}

	response, err := e.service.GetExchangeRateResponse(rate)
	if err != nil {
		utilities.HandleError(w, err)
		return
	}

	json.NewEncoder(w).Encode(response)
}

func (e *ExchangeRateController) patchExchangeRate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	baseCurrencyCode, targetCurrencyCode, err := utilities.GetCurrencyCodes(r.URL.Path)
	if err != nil {
		utilities.HandleError(w, err)
		return
	}

	rateValue, err := strconv.ParseFloat(r.FormValue("rate"), 64)

	rate, err := e.service.UpdateExchangeRate(baseCurrencyCode, targetCurrencyCode, rateValue)
	if err != nil {
		utilities.HandleError(w, err)
		return
	}

	json.NewEncoder(w).Encode(rate)
}
