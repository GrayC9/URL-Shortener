package main

import (
	"log"
	"net/http"
	"url-shortener/internal/storage/memory"
)

func main() {
	store := memory.NewMemoryStore()
	//http.HandleFunc("/shorten", handlers.ShortenURL(store))
	//http.HandleFunc("/expand", handlers.(store))

	log.Println("Server port 10000")
	err := http.ListenAndServe(":10000", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
