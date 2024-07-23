package handlers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"url_shortener/internal/cache"
	"url_shortener/internal/shortener"
	"url_shortener/internal/storage"

	"github.com/gorilla/mux"
)

var tmpl = template.Must(template.ParseFiles("internal/web/index.html"))

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
				ShortURL:    cacheEntry.ShortURL,
			}
			err := tmpl.Execute(w, data)
			if err != nil {
				log.Printf("Error executing template: %v", err)
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
				http.Error(w, "Error saving URL", http.StatusInternalServerError)
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
			log.Printf("Error executing template: %v", err)
		}
	}
}

func RedirectHandler(db storage.Storage, urlCache *cache.URLCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		shortCode := vars["shortCode"]

		cacheEntry, exists := urlCache.GetEntry(shortCode)
		if exists {
			cacheEntry.Count++
			urlCache.AddEntry(cacheEntry.OriginalURL, cacheEntry.ShortURL)

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
			http.Error(w, "Error update click count", http.StatusInternalServerError)
			return
		}

		err = db.UpdateLastAccessed(shortCode)
		if err != nil {
			http.Error(w, "Error updating last accesseded time", http.StatusInternalServerError)
			return
		}

		urlCache.AddEntry(originalURL, shortCode)

		http.Redirect(w, r, originalURL, http.StatusFound)
	}
}

func WebInterfaceHandler(db storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			CreateShortURLHandler(db, nil)(w, r)
			return
		}
		err := tmpl.Execute(w, PageData{})
		if err != nil {
			log.Printf("Error executing template: %v", err)
		}
	}
}

func WriteJSON(w http.ResponseWriter, r *http.Request, status int, a interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(a)
}
