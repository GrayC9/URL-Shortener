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
	cfg := config.LoadConfig()        // Функция для загрузки конфигурации
	db := storage.NewDatabase(cfg.DB) // Инициализация базы данных

	r := mux.NewRouter()
	r.HandleFunc("/shorten", handlers.CreateShortURLHandler(db)).Methods("POST")
	r.HandleFunc("/{shortCode}", handlers.RedirectHandler(db)).Methods("GET")

	log.Fatal(http.ListenAndServe(cfg.Server.Address, r))
}
