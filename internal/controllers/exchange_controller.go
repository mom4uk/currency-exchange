package controllers

import (
	"currency-exchange/internal/services"
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
	baseCurrency := r.URL.Query().Get("from")

	targetCurrency := r.URL.Query().Get("to")

	amountStr := r.URL.Query().Get("amount")
	amountValue := new(big.Rat)
	_, ok := amountValue.SetString(amountStr)
	if !ok {
		http.Error(w, "Ошибка в amount", http.StatusInternalServerError)
		return
	}

	exchange, err := e.service.GetExchange(baseCurrency, targetCurrency, amountValue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(exchange); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
