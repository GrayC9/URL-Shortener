package main

import (
	"log"

	"github.com/GrayC9/URL-Shortener/internal/handlers"
	"github.com/GrayC9/URL-Shortener/internal/server"
)

func main() {
	router := handlers.NewRouter()
	go log.Fatalln(server.New().Run(router))
}
