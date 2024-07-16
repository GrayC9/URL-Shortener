package shortener

import (
	"testing"
	"time"
	"unicode/utf8"
)

func TestGenerateShortCode(t *testing.T) {
	code := GenerateShortCode()
	if utf8.RuneCountInString(code) != 6 {
		t.Errorf("expected length of 6  %d", utf8.RuneCountInString(code))
	}

	for _, r := range code {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') {
			t.Errorf("invalid character in code: %c", r)
		}
	}

	codeSet := make(map[string]bool)
	for i := 0; i < 1000; i++ {
		code := GenerateShortCode()
		if _, exists := codeSet[code]; exists {
			t.Errorf("duplicat code found: %s", code)
		}
		codeSet[code] = true
	}
}

func TestURLStruct(t *testing.T) {
	url := URL{
		OriginalURL:    "http://example.com",
		ShortCode:      GenerateShortCode(),
		LastAccessedAt: time.Now(),
	}

	if url.OriginalURL != "http://example.com" {
		t.Errorf("expected originalURL to 'http://example.com' %s", url.OriginalURL)
	}

	if utf8.RuneCountInString(url.ShortCode) != 6 {
		t.Errorf("expected shortCode length of 6 %d", utf8.RuneCountInString(url.ShortCode))
	}

	if time.Since(url.LastAccessedAt) > time.Second {
		t.Errorf(": %s", url.LastAccessedAt)
	}
}
