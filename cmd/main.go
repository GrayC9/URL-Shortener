package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"url_shortener/internal/config"
	"url_shortener/internal/handlers"
	"url_shortener/internal/storage"
)

func main() {
	cfg := config.LoadConfig()

	db, err := storage.NewMariaDBStorage(cfg.DB.DSN)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./internal/web"))))
	r.HandleFunc("/shorten", handlers.CreateShortURLHandler(db)).Methods("POST")
	r.HandleFunc("/{shortCode}", handlers.RedirectHandler(db)).Methods("GET")
	r.HandleFunc("/", handlers.WebInterfaceHandler(db)).Methods("GET", "POST")

	log.Fatal(http.ListenAndServe(cfg.Server.Address, r))
}
