package controllers

import (
	"currency-exchange/internal/dto"
	"currency-exchange/internal/services"
	"currency-exchange/internal/utilities"
	"encoding/json"
	"math/big"
	"net/http"
)

type ExchangeController struct {
	service services.ExchangeService
}

func NewExchangeController(srv *services.ExchangeService) *ExchangeController {
	return &ExchangeController{
		service: *srv,
	}
}

func (e *ExchangeController) GetExchange(w http.ResponseWriter, r *http.Request) {
	baseCurrencyCode := r.URL.Query().Get("from")

	targetCurrencyCode := r.URL.Query().Get("to")

	amountStr := r.URL.Query().Get("amount")
	if err := dto.ValidateExchangeFields(baseCurrencyCode, targetCurrencyCode, amountStr); err != nil {
		utilities.HandleError(w, err)
		return
	}

	amountValue := new(big.Rat)
	_, ok := amountValue.SetString(amountStr)
	if !ok {
		http.Error(w, "Ошибка в amount", http.StatusInternalServerError)
		return
	}

	exchange, err := e.service.GetExchange(baseCurrencyCode, targetCurrencyCode, amountValue)
	if err != nil {
		utilities.HandleError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(exchange); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
