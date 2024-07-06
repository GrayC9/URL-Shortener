package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"url_shortener/internal/shortener"
	"url_shortener/internal/storage"
)

func CreateShortURLHandler(db storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var url shortener.URL
		err := json.NewDecoder(r.Body).Decode(&url)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		url.ShortCode = shortener.GenerateShortCode()
		db.SaveURL(url.ShortCode, url.OriginalURL)

		json.NewEncoder(w).Encode(url)
	}
}

func RedirectHandler(db storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		shortCode := vars["shortCode"]

		originalURL, err := db.GetURL(shortCode)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		http.Redirect(w, r, originalURL, http.StatusFound)
	}
}
