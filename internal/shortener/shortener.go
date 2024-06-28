package shortener

import (
	"math/rand"
	"time"
)

type URL struct {
	OriginalURL string `json:"original"`
	ShortURL    string `json:"short"`
}

func MakeShort() string {
	rand.Seed(time.Now().UnixNano())
	alphabet := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()")
	short := make([]rune, 6)
	for rn := range short {
		short[rn] = alphabet[rand.Intn(len(alphabet))]
	}
	return string(short)
}
