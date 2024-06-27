package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Router struct {
	R *mux.Router
}

func NewRouter() *Router {
	return &Router{
		R: mux.NewRouter(),
	}
}

func (r *Router) InitRoutes() {
	r.R.HandleFunc("/", homepage).Methods("GET")
	r.R.HandleFunc("/original", getOriginal).Methods("GET")
	r.R.HandleFunc("/short", shortened).Methods("GET")
}

func homepage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello World!"))
}

func getOriginal(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Original"))

}

func shortened(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Shortened"))
}
