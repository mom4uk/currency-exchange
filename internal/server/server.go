package server

import (
	"currency-exchange/internal/controllers"
	"currency-exchange/internal/utilities"
	"net/http"
)

type Server struct {
	mux *http.ServeMux
}

func New() *Server {
	return &Server{
		mux: http.NewServeMux(),
	}
}

func (s *Server) RegisterRoutes(controller *controllers.CurrencyController) {
	s.mux.Handle("/currencies",
		utilities.JSON(http.HandlerFunc(controller.HandleCurrencies)),
	)
	s.mux.Handle("/currency/",
		utilities.JSON(http.HandlerFunc(controller.GetCurrency)),
	)
	s.mux.Handle("/exchangeRates",
		utilities.JSON(http.HandlerFunc(controller.HandleExchangeRates)),
	)
	s.mux.Handle("/exchangeRate/",
		utilities.JSON(http.HandlerFunc(controller.HandleExchangeRate)),
	)
	s.mux.Handle("/exchange",
		utilities.JSON(http.HandlerFunc(controller.GetExchange)),
	)
}

func (s *Server) Start() error {
	return http.ListenAndServe(":8080", s.mux)
}
