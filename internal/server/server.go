package server

import (
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
	return http.ListenAndServe(os.Getenv("ADDR"), s.mux)
}
