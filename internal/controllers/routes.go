package controllers

import (
	"currency-exchange/internal/middleware"
	"net/http"
)

func RegisterCurrencyRoutes(mux *http.ServeMux, c *CurrencyController) {
	mux.Handle(
		"GET /currencies",
		middleware.JSON(http.HandlerFunc(c.GetCurrencies)),
	)
	mux.Handle(
		"POST /currencies",
		middleware.JSON(http.HandlerFunc(c.AddCurrency)),
	)

	mux.Handle(
		"GET /currency/",
		middleware.JSON(http.HandlerFunc(c.GetCurrency)),
	)
}

func RegisterExchangeRateRoutes(mux *http.ServeMux, c *ExchangeRateController) {
	mux.Handle(
		"GET /exchangeRates",
		middleware.JSON(http.HandlerFunc(c.GetExchangeRates)),
	)
	mux.Handle(
		"POST /exchangeRates",
		middleware.JSON(http.HandlerFunc(c.AddExchangeRates)),
	)

	mux.Handle(
		"GET /exchangeRate/",
		middleware.JSON(http.HandlerFunc(c.GetExchangeRate)),
	)
	mux.Handle(
		"PATCH /exchangeRate/",
		middleware.JSON(http.HandlerFunc(c.UpdateExchangeRate)),
	)
}

func RegisterExchangeRoutes(mux *http.ServeMux, c *ExchangeController) {
	mux.Handle(
		"GET /exchange",
		middleware.JSON(http.HandlerFunc(c.GetExchange)),
	)
}
