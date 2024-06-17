package main

import (
	"github.com/GrayC9/URL-Shortener/internal/server"
)

func main() {
	s := server.New()
	s.Run()
}
