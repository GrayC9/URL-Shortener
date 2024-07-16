package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"

	"url_shortener/internal/auth"

	"url_shortener/internal/config"
	"url_shortener/internal/handlers"
	"url_shortener/internal/storage"

	"github.com/gorilla/mux"
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
	r.HandleFunc("/api/shorten", handlers.CreateShortURLHandler(db)).Methods("POST")
	r.HandleFunc("/{shortCode}", handlers.RedirectHandler(db)).Methods("GET")
	r.HandleFunc("/", handlers.WebInterfaceHandler(db)).Methods("GET", "POST")
	r.HandleFunc("/register", auth.SignUp(db)).Methods("POST")

	//subrouter := r.HandleFunc("/login", auth.Login(db)).Subrouter()
	//subrouter.Use(auth.AuthMiddleware)
	//subrouter.Methods("POST")

	r.HandleFunc("/login", auth.AuthMiddleware((auth.Login(db)))).Methods("POST")

	log.Fatal(http.ListenAndServe(cfg.Server.Address, r))
}
