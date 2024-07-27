package tests

import (
	"bytes"
	"embed"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
	"url_shortener/internal/cache"
	"url_shortener/internal/handlers"
	"url_shortener/internal/storage"
)

//go:embed values.txt
var valuesFile embed.FS

func loadURLs(filePath string) ([]string, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	lines := bytes.Split(file, []byte("\n"))
	urls := make([]string, 0, len(lines))
	for _, line := range lines {
		if len(line) > 0 {
			urls = append(urls, string(line))
		}
	}
	return urls, nil
}

func TestShorteningPerformance(t *testing.T) {
	urls, err := loadURLs(valuesFile)
	if err != nil {
		t.Fatalf("Failed load URLs from file: %v", err)
	}

	db, err := storage.NewMariaDBStorage("секретка")
	if err != nil {
		t.Fatalf("Failed to create st: %v", err)
	}
	urlCache := cache.NewURLCache()

	r := mux.NewRouter()
	r.HandleFunc("/", handlers.WebInterfaceHandler(db, urlCache)).Methods("GET", "POST")
	server := httptest.NewServer(r)
	defer server.Close()

	for _, originalURL := range urls {
		start := time.Now()
		resp, err := http.PostForm(server.URL+"/", url.Values{"original_url": {originalURL}})
		if err != nil {
			t.Fatalf("Failed shorten URL: %v", err)
		}
		resp.Body.Close()
		firstDuration := time.Since(start)

		start = time.Now()
		resp, err = http.PostForm(server.URL+"/", url.Values{"original_url": {originalURL}})
		if err != nil {
			t.Fatalf("Failed to shorten URL: %v", err)
		}
		resp.Body.Close()
		secondDuration := time.Since(start)

		t.Logf("First shortening took %v for URL: %s", firstDuration, originalURL)
		t.Logf("Second shortening took %v for URL: %s", secondDuration, originalURL)

		if secondDuration >= firstDuration {
			t.Errorf("Second URL shortening did not improve for URL: %s (first: %v, second: %v)", originalURL, firstDuration, secondDuration)
		}
	}
}
