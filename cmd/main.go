package main

import (
	"currency-exchange/db"
	"currency-exchange/internal/controllers"
	"currency-exchange/internal/repositories"
	"currency-exchange/internal/server"
	"currency-exchange/internal/services"
	"log"
)

func main() {
	db := db.InitDb()

	srv := server.New()

	CurrencyRepository := repositories.CurrencyRepositoryNew(db)
	ExchangeRateRepository := repositories.ExchangeRateRepositoryNew(db)

	currencyService := services.CurrencyServiceNew(ExchangeRateRepository, CurrencyRepository)
	exchangeService := services.ExchangeRateServiceNew(ExchangeRateRepository, CurrencyRepository)

	currencyController := controllers.NewController(currencyService)
	exchangeRateController := controllers.NewExchangeRateController(exchangeService)

	controllers.RegisterCurrencyRoutes(srv.GetMux(), currencyController)
	controllers.RegisterExchangeRoutes(srv.GetMux(), exchangeRateController)

	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
