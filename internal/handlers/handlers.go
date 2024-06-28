package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/GrayC9/URL-Shortener/internal/shortener"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Router struct {
	R    *mux.Router
	logg *logrus.Logger
}

func NewRouter(log *logrus.Logger) *Router {
	return &Router{
		R:    mux.NewRouter(),
		logg: log,
	}
}

type URL struct {
	OriginalURL string `json:"original"`
	ShortURL    string `json:"short"`
}

func (r *Router) InitRoutes() {
	r.R.HandleFunc("/original", r.getOriginal).Methods("GET")
	r.R.HandleFunc("/short", r.shortened).Methods("POST")
}

func (r *Router) getOriginal(w http.ResponseWriter, req *http.Request) {
	url := &URL{}
	if err := json.NewDecoder(req.Body).Decode(&url); err != nil {
		r.logg.Errorln(err)
		return
	}
	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"original": url.OriginalURL,
	})
}

func (r *Router) shortened(w http.ResponseWriter, req *http.Request) {
	url := &URL{}
	if err := json.NewDecoder(req.Body).Decode(&url); err != nil {
		r.logg.Errorln(err)
		return
	}
	url.ShortURL = shortener.MakeShort()
	if url.ShortURL == "" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		r.logg.Errorln("shortened URL is empty")
		return
	}
	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"original": url.OriginalURL,
		"short":    url.ShortURL,
	})
	//http.Redirect(w, req, newURL.OriginalURL, http.StatusFound)
}

func WriteJSON(w http.ResponseWriter, status int, a interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(a)
}
