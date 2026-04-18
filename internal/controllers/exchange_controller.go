package controllers

import (
	"currency-exchange/internal/services"
	"encoding/json"
	"net/http"
	"strconv"
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

	amount, err := strconv.ParseFloat(r.URL.Query().Get("amount"), 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	exchange, err := e.service.GetExchange(baseCurrency, targetCurrency, amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(exchange); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
