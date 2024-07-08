package handlers

import (
	"encoding/json"
	"net/http"
	"url_shortener/internal/shortener"
	"url_shortener/internal/storage"

	"github.com/gorilla/mux"
)

func CreateShortURLHandler(db storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var url shortener.URL
		err := json.NewDecoder(r.Body).Decode(&url)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		shortCode, err := db.GetShortCode(url.OriginalURL)
		if err == nil {
			url.ShortCode = shortCode
		} else {
			url.ShortCode = shortener.GenerateShortCode()
			err = db.SaveURL(url.ShortCode, url.OriginalURL)
			if err != nil {
				http.Error(w, "ERorr to save url", http.StatusInternalServerError)
				return
			}
		}

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
