// main.go
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

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./internal/web"))))
	//r.HandleFunc("/api/shorten", handlers.CreateShortURLHandler(db, urlCache)).Methods("POST")
	r.HandleFunc("/{shortCode}", handlers.RedirectHandler(db, urlCache)).Methods("GET")
	r.HandleFunc("/", handlers.WebInterfaceHandler(db, urlCache)).Methods("GET", "POST")
	// проверка
	//fmt.Println("Получение записей из кэша...")
	//entry, exists := urlCache.GetEntry("short1")
	//if exists {
	//	fmt.Printf("Запись найдена для short1: %+v\n", entry)
	//} else {
	//	fmt.Println("Запись не найдена для short1")
	//}
	//
	//entry, exists = urlCache.GetEntry("short2")
	//if exists {
	//	fmt.Printf("Запись найдена для short2: %+v\n", entry)
	//} else {
	//	fmt.Println("Запись не найдена для short2")
	//}
	//
	//fmt.Println("Удаление записи short2 из кэша...")
	//urlCache.DeleteEntry("short2")
	//entry, exists = urlCache.GetEntry("short2")
	//if exists {
	//	fmt.Printf("Запись все еще существует для short2: %+v\n", entry)
	//} else {
	//	fmt.Println("Запись успешно удалена для short2")
	//}
	// проверка
	r.HandleFunc("/register", auth.SignUp(db)).Methods("POST")
	r.HandleFunc("/login", auth.AuthMiddleware((auth.Login(db)))).Methods("POST")

	log.Fatal(http.ListenAndServe(cfg.Server.Address, r))
}
