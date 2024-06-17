package server

import "net/http"

type Server struct {
	srv *http.Server
}

func New() *Server {
	return &Server{}
}

func (s *Server) Run() error {
	s.srv = &http.Server{
		Addr: ":10000",
	}
	return s.srv.ListenAndServe()
}
