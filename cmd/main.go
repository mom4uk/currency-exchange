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

	CurrencyRepository := repositories.CurrencyRepositoryNew(db)
	ExchangeRateRepository := repositories.ExchangeRateRepositoryNew(db)

	currencyService := services.CurrencyServiceNew(CurrencyRepository)
	exchangeService := services.ExchangeRateServiceNew(ExchangeRateRepository)

	currencyController := controllers.NewController(currencyService)
	exchangeRateController := controllers.NewController(exchangeService)

	srv := server.New()
	srv.RegisterRoutes(currencyController, exchangeRateController)

	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
