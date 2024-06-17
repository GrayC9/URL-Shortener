package main

import (
	"log"

	"github.com/GrayC9/URL-Shortener/internal/server"
)

func main() {
	s := server.New()
	log.Fatal(s.Run())
}
