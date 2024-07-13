package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"url_shortener/internal/auth"
	"url_shortener/internal/config"
	"url_shortener/internal/handlers"
	"url_shortener/internal/storage"
)

func main() {
	cfg := config.LoadConfig()
	auth.JWTSecretKey = []byte(cfg.Server.JWTSecret)

	r := mux.NewRouter()
	db, err := storage.NewMariaDBStorage(cfg.DB.DSN)
	if err != nil {
		r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			handlers.WriteJSON(w, r, http.StatusInternalServerError, err.Error())
		})
	}

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./internal/web"))))
	r.HandleFunc("/shorten", handlers.CreateShortURLHandler(db)).Methods("POST")
	r.HandleFunc("/{shortCode}", handlers.RedirectHandler(db)).Methods("GET")
	r.HandleFunc("/", handlers.WebInterfaceHandler(db)).Methods("GET", "POST")

	// Пример использования middleware.AuthMiddleware
	r.Use(auth.AuthMiddleware)

	log.Fatal(http.ListenAndServe(cfg.Server.Address, r))
}
