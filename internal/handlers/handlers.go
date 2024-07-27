package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"url_shortener/internal/cache"
	"url_shortener/internal/shortener"
	"url_shortener/internal/storage"
)

//go:embed web/index.html
var indexHTML string

var tmpl = template.Must(template.New("index").Parse(indexHTML))

type PageData struct {
	OriginalURL string
	ShortURL    string
	Error       string
}

func CreateShortURLHandler(db storage.Storage, urlCache *cache.URLCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		originalURL := r.FormValue("original_url")

		cacheEntry, exists := urlCache.GetEntry(originalURL)
		if exists {
			data := PageData{
				OriginalURL: originalURL,
				ShortURL:    r.Host + "/" + cacheEntry.ShortURL,
			}
			err := tmpl.Execute(w, data)
			if err != nil {
				log.Printf("Ошибка при выполнении шаблона: %v", err)
			}
			return
		}

		var url shortener.URL
		url.OriginalURL = originalURL

		shortCode, err := db.GetShortCode(url.OriginalURL)
		if err == nil {
			url.ShortCode = shortCode
		} else {
			url.ShortCode = shortener.GenerateShortCode()
			err = db.SaveURL(url.ShortCode, url.OriginalURL)
			if err != nil {
				http.Error(w, "Ошибка при сохранении URL", http.StatusInternalServerError)
				return
			}
		}

		urlCache.AddEntry(url.OriginalURL, url.ShortCode)

		data := PageData{
			OriginalURL: url.OriginalURL,
			ShortURL:    r.Host + "/" + url.ShortCode,
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			log.Printf("Ошибка при выполнении шаблона: %v", err)
		}
	}
}

func RedirectHandler(db storage.Storage, urlCache *cache.URLCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		shortCode := vars["shortCode"]

		if cacheEntry, ok := urlCache.GetEntry(shortCode); ok {
			urlCache.IncrementCount(shortCode)
			http.Redirect(w, r, cacheEntry.OriginalURL, http.StatusFound)
			return
		}

		originalURL, err := db.GetURL(shortCode)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		err = db.IncrementClickCount(shortCode)
		if err != nil {
			http.Error(w, "Ошибка при обновлении счетчика кликов", http.StatusInternalServerError)
			return
		}

		err = db.UpdateLastAccessed(shortCode)
		if err != nil {
			http.Error(w, "Ошибка при обновлении времени последнего доступа", http.StatusInternalServerError)
			return
		}

		urlCache.AddEntry(originalURL, shortCode)
		urlCache.IncrementCount(shortCode)

		http.Redirect(w, r, originalURL, http.StatusFound)
	}
}

func WebInterfaceHandler(db storage.Storage, c *cache.URLCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			CreateShortURLHandler(db, c)(w, r)
			return
		}
		err := tmpl.Execute(w, PageData{})
		if err != nil {
			log.Printf("Ошибка при выполнении шаблона: %v", err)
		}
	}
}

func WriteJSON(w http.ResponseWriter, r *http.Request, status int, a interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(a)
}
