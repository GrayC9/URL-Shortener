package shortener

import (
	"math/rand"
	"time"
)

type URL struct {
	OriginalURL string `json:"original_url"`
	ShortCode   string `json:"short_code"`
}

func GenerateShortCode() string {
	rand.Seed(time.Now().UnixNano())
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, 6)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
