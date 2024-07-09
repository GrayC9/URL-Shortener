package handlers

import (
	"html/template"
	"net/http"
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

func CreateShortURLHandler(db storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		originalURL := r.FormValue("original_url")

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

		data := PageData{
			OriginalURL: url.OriginalURL,
			ShortURL:    r.Host + "/" + url.ShortCode,
		}

		tmpl.Execute(w, data)
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

func WebInterfaceHandler(db storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			CreateShortURLHandler(db)(w, r)
			return
		}
		tmpl.Execute(w, PageData{})
	}
}
