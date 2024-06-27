package server

import (
	"net/http"
	"time"

	"github.com/GrayC9/URL-Shortener/internal/config"
	"github.com/GrayC9/URL-Shortener/internal/handlers"
)

type Server struct {
	srv *http.Server
	//router *handlers.Router
	cnf *config.Config
}

func New() *Server {
	return &Server{
		srv: &http.Server{},
		//router: handlers.NewRouter(),
		cnf: config.NewConfig(),
	}
}

func (s *Server) Run(router *handlers.Router) error {
	router.InitRoutes()

	s.srv = &http.Server{
		Addr:           s.cnf.Addr_Port,
		Handler:        router.R,
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}
	return s.srv.ListenAndServe()
}
