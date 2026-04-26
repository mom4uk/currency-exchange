package controllers

import (
	"currency-exchange/internal/utilities"
	"net/http"
)

func RegisterCurrencyRoutes(mux *http.ServeMux, c *CurrencyController) {
	mux.Handle(
		"GET /currencies",
		utilities.JSON(http.HandlerFunc(c.GetCurrencies)),
	)
	mux.Handle(
		"POST /currencies",
		utilities.JSON(http.HandlerFunc(c.AddCurrency)),
	)

	mux.Handle(
		"GET /currency/",
		utilities.JSON(http.HandlerFunc(c.GetCurrency)),
	)
}

func RegisterExchangeRateRoutes(mux *http.ServeMux, c *ExchangeRateController) {
	mux.Handle(
		"GET /exchangeRates",
		utilities.JSON(http.HandlerFunc(c.GetExchangeRates)),
	)
	mux.Handle(
		"POST /exchangeRates",
		utilities.JSON(http.HandlerFunc(c.AddExchangeRates)),
	)

	mux.Handle(
		"GET /exchangeRate/",
		utilities.JSON(http.HandlerFunc(c.GetExchangeRate)),
	)
	mux.Handle(
		"PATCH /exchangeRate/",
		utilities.JSON(http.HandlerFunc(c.UpdateExchangeRate)),
	)
}

func RegisterExchangeRoutes(mux *http.ServeMux, c *ExchangeController) {
	mux.Handle(
		"GET /exchange",
		utilities.JSON(http.HandlerFunc(c.GetExchange)),
	)
}
