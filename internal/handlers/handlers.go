package handlers

import (
	"net/http"

	"github.com/GrayC9/URL-Shortener/internal/config"
	"github.com/gorilla/mux"
)

type Router struct {
	R    *mux.Router
	conf *config.Config
}

func NewRouter() *Router {
	return &Router{
		R:    mux.NewRouter(),
		conf: config.NewConfig(),
	}
}

func (r *Router) InitRoutes() {
	r.R.NewRoute()
	s := r.R.Host(r.conf.Addr_Port).Subrouter()
	s.HandleFunc("/", homepage).Methods("GET")
}

func homepage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello World!"))
}
