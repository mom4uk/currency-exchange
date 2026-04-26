package controllers

import (
	"currency-exchange/internal/domain"
	"currency-exchange/internal/dto"
	"currency-exchange/internal/services"
	"currency-exchange/internal/utilities"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
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
		e.updateExchangeRate(w, r)
	default:
		http.Error(w, "This method is not allowed", http.StatusMethodNotAllowed)
	}
}

func (e *ExchangeRateController) addExchangeRates(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req := dto.AddExchangeRateRequest{
		BaseCurrencyCode:   r.FormValue("baseCurrencyCode"),
		TargetCurrencyCode: r.FormValue("targetCurrencyCode"),
		Rate:               r.FormValue("rate"),
	}

	if err := dto.ValidateExchangeRateFields(req); err != nil {
		utilities.HandleError(w, err)
		return
	}

	rateStr := r.FormValue("rate")

	rate := new(big.Rat)

	_, ok := rate.SetString(rateStr)
	if !ok {
		utilities.WriteError(
			w,
			fmt.Sprintf("Неверный формат суммы: %s", rateStr),
			http.StatusBadRequest,
		)
		return
	}

	exchangeRate := domain.Exchange{
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

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (e *ExchangeRateController) getExchangeRates(w http.ResponseWriter) {
	rates, err := e.service.GetExchangeRates()
	if err != nil {
		utilities.HandleError(w, err)
		return
	}

	response, err := e.service.GetExchangeRatesResponse(rates)
	if err != nil {
		utilities.HandleError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (e *ExchangeRateController) updateExchangeRate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	baseCurrencyCode, targetCurrencyCode, err := utilities.GetCurrencyCodes(r.URL.Path)
	if err != nil {
		utilities.HandleError(w, err)
		return
	}

	req := dto.UpdateExchangeRateRequest{
		Rate: r.FormValue("rate"),
	}

	if err := dto.ValidateExchangeRateFieldsForUpdate(req); err != nil {
		utilities.HandleError(w, err)
		return
	}
	rateStr := r.FormValue("rate")
	rateValue := new(big.Rat)
	_, ok := rateValue.SetString(rateStr)
	if !ok {
		utilities.HandleError(w, domain.ErrRateConvertaion)
		return
	}

	rate, err := e.service.UpdateExchangeRate(baseCurrencyCode, targetCurrencyCode, rateValue)
	if err != nil {
		utilities.HandleError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(rate); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
