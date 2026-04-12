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

	repo := repositories.CurrencyRepositoryNew(db)
	service := services.CurrencyServiceNew(repo)
	controller := controllers.NewController(repo, service)

	srv := server.New()
	srv.RegisterRoutes(controller)

	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
