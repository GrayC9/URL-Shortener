package main

import (
	"log"
	"net/http"

	"url_shortener/internal/config"
	"url_shortener/internal/handlers"
	"url_shortener/internal/storage"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig()

	r := mux.NewRouter()
	db, err := storage.NewMariaDBStorage(cfg.DB.DSN)
	if err != nil {
		r.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
			handlers.WriteJSON(w, r, http.StatusInternalServerError, err.Error())
		})
	}

	r.HandleFunc("/shorten", handlers.CreateShortURLHandler(db)).Methods("POST")
	r.HandleFunc("/{shortCode}", handlers.RedirectHandler(db)).Methods("GET")

	log.Fatal(http.ListenAndServe(cfg.Server.Address, r))
}
