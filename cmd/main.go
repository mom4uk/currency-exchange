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

	currencyService := services.CurrencyServiceNew(CurrencyRepository)
	exchangeRateService := services.ExchangeRateServiceNew(ExchangeRateRepository, CurrencyRepository)
	exchangeService := services.ExchangeServiceNew(ExchangeRateRepository, CurrencyRepository)

	currencyController := controllers.NewController(currencyService)
	exchangeRateController := controllers.NewExchangeRateController(exchangeRateService)
	exchangeController := controllers.NewExchangeController(exchangeService)

	controllers.RegisterCurrencyRoutes(srv.GetMux(), currencyController)
	controllers.RegisterExchangeRateRoutes(srv.GetMux(), exchangeRateController)
	controllers.RegisterExchangeRoutes(srv.GetMux(), exchangeController)

	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
