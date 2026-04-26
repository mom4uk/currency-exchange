package controllers

import (
	"currency-exchange/internal/apierrors"
	"currency-exchange/internal/domain"
	"currency-exchange/internal/dto"
	"currency-exchange/internal/httputil"
	"currency-exchange/internal/services"
	"encoding/json"
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

func (e *ExchangeRateController) AddExchangeRates(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		apierrors.WriteError(w, "Parse form error", http.StatusBadRequest)
		return
	}

	req := dto.AddExchangeRateRequest{
		BaseCurrencyCode:   r.FormValue("baseCurrencyCode"),
		TargetCurrencyCode: r.FormValue("targetCurrencyCode"),
		Rate:               r.FormValue("rate"),
	}

	if err := dto.ValidateExchangeRateFields(req); err != nil {
		apierrors.HandleError(w, err)
		return
	}

	rateStr := r.FormValue("rate")

	rate := new(big.Rat)

	_, ok := rate.SetString(rateStr)
	if !ok {
		apierrors.HandleError(w, domain.ErrRateConvertaion)
		return
	}

	exchangeRate := domain.Exchange{
		BaseCurrencyCode:   req.BaseCurrencyCode,
		TargetCurrencyCode: req.TargetCurrencyCode,
		Rate:               rate,
	}

	result, err := e.service.AddExchangeRates(exchangeRate)
	if err != nil {
		apierrors.HandleError(w, err)
		return
	}

	res, err := e.service.GetExchangeRateResponse(result)
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

func (e *ExchangeRateController) GetExchangeRates(w http.ResponseWriter, r *http.Request) {
	rates, err := e.service.GetExchangeRates()
	if err != nil {
		apierrors.HandleError(w, err)
		return
	}

	response, err := e.service.GetExchangeRatesResponse(rates)
	if err != nil {
		apierrors.HandleError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		apierrors.WriteError(w, "Json convertation error", http.StatusInternalServerError)
		return
	}
}

func (e *ExchangeRateController) GetExchangeRate(w http.ResponseWriter, r *http.Request) {
	baseCurrencyCode, targetCurrencyCode, err := httputil.GetCurrencyCodes(r.URL.Path)
	if err != nil {
		apierrors.HandleError(w, err)
		return
	}

	rate, err := e.service.GetExchangeRateByCodes(baseCurrencyCode, targetCurrencyCode)
	if err != nil {
		apierrors.HandleError(w, err)
		return
	}

	response, err := e.service.GetExchangeRateResponse(rate)
	if err != nil {
		apierrors.HandleError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		apierrors.WriteError(w, "Json convertation error", http.StatusInternalServerError)
		return
	}
}

func (e *ExchangeRateController) UpdateExchangeRate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		apierrors.WriteError(w, "Parse form error", http.StatusBadRequest)
		return
	}

	baseCurrencyCode, targetCurrencyCode, err := httputil.GetCurrencyCodes(r.URL.Path)
	if err != nil {
		apierrors.HandleError(w, err)
		return
	}

	req := dto.UpdateExchangeRateRequest{
		Rate: r.FormValue("rate"),
	}

	if err := dto.ValidateExchangeRateFieldsForUpdate(req); err != nil {
		apierrors.HandleError(w, err)
		return
	}
	rateStr := r.FormValue("rate")
	rateValue := new(big.Rat)
	_, ok := rateValue.SetString(rateStr)
	if !ok {
		apierrors.HandleError(w, domain.ErrRateConvertaion)
		return
	}

	rate, err := e.service.UpdateExchangeRate(baseCurrencyCode, targetCurrencyCode, rateValue)
	if err != nil {
		apierrors.HandleError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(rate); err != nil {
		apierrors.WriteError(w, "Json convertation error", http.StatusInternalServerError)
		return
	}
}
