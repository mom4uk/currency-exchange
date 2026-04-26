package server

import (
	"log"
	"net/http"
	"os"
)

type Server struct {
	mux *http.ServeMux
}

func New() *Server {
	return &Server{
		mux: http.NewServeMux(),
	}
}

func (s *Server) GetMux() *http.ServeMux {
	return s.mux
}

func (s *Server) Start() error {
	addr := os.Getenv("ADDR")
	log.Println("starting server on:", addr)
	if addr == "" {
		addr = ":8080"
	}

	return http.ListenAndServe(addr, s.mux)
}
