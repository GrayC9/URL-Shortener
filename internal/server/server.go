package server

import (
	"context"
	"net/http"
	"time"

	"github.com/GrayC9/URL-Shortener/internal/config"
	"github.com/GrayC9/URL-Shortener/internal/handlers"
	"github.com/sirupsen/logrus"
)

type Server struct {
	srv  *http.Server
	cnf  *config.Config
	logg *logrus.Logger
}

func New() *Server {
	return &Server{
		srv:  &http.Server{},
		cnf:  config.NewConfig(),
		logg: logrus.New(),
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
	s.logg.Infoln("listening --> " + s.cnf.Addr_Port)
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
