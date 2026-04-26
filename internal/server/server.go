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
	err := http.ListenAndServe(os.Getenv("ADDR"), s.mux)
	if err != nil {
		log.Fatal(err)
	}
	return err
}
