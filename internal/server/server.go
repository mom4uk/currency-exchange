package server

import (
	"currency-exchange/internal/controllers"
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
	s.mux.HandleFunc("/currencies", controller.HandleCurrencies)
	s.mux.HandleFunc("/currency/", controller.GetCurrency)
	s.mux.HandleFunc("/exchangeRates", controller.HandleExchangeRates)
	s.mux.HandleFunc("/exchangeRate/", controller.HandleExchangeRate)
	s.mux.HandleFunc("/exchange", controller.GetExchange)
}

func (s *Server) Start() error {
	return http.ListenAndServe(":8080", s.mux)
}
