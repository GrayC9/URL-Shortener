package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"url_shortener/internal/auth"
	"url_shortener/internal/cache"
	"url_shortener/internal/config"
	"url_shortener/internal/handlers"
	"url_shortener/internal/storage"
)

const popularURLLimit = 1000

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

	urlCache := cache.NewURLCache()
	preloadCache(db, urlCache)

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./internal/web"))))
	r.HandleFunc("/{shortCode}", handlers.RedirectHandler(db, urlCache)).Methods("GET")
	r.HandleFunc("/", handlers.WebInterfaceHandler(db, urlCache)).Methods("GET", "POST")

	r.HandleFunc("/register", auth.SignUp(db)).Methods("POST")
	r.HandleFunc("/login", auth.AuthMiddleware(auth.Login(db))).Methods("POST")

	log.Fatal(http.ListenAndServe(cfg.Server.Address, r))
}

func preloadCache(db storage.Storage, urlCache *cache.URLCache) {
	popularURLs, err := db.GetPopularURLs(popularURLLimit)
	if err != nil {
		log.Printf("Ошибка при получении популярных URL: %v", err)
		return
	}

	for _, url := range popularURLs {
		urlCache.AddEntry(url.OriginalURL, url.ShortCode)
		log.Printf("URL добавлен в кеш: %s -> %s", url.ShortCode, url.OriginalURL)
	}
}

func logRequest(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Получен запрос: %s %s", r.Method, r.RequestURI)
		handler(w, r)
	}
}
