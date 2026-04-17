package controllers

import (
	"currency-exchange/internal/utilities"
	"net/http"
)

func RegisterCurrencyRoutes(mux *http.ServeMux, c *CurrencyController) {
	mux.Handle("/currencies", utilities.JSON(http.HandlerFunc(c.HandleCurrencies)))
	mux.Handle("/currency/", utilities.JSON(http.HandlerFunc(c.GetCurrency)))
	mux.Handle("/exchange", utilities.JSON(http.HandlerFunc(c.GetExchange)))
}

func RegisterExchangeRoutes(mux *http.ServeMux, c *ExchangeRateController) {
	mux.Handle("/exchangeRates", utilities.JSON(http.HandlerFunc(c.HandleExchangeRates)))
	mux.Handle("/exchangeRate/", utilities.JSON(http.HandlerFunc(c.HandleExchangeRate)))
}
