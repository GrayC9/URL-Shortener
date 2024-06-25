package storage

import (
	"errors"
	"sync"
)

type Storage interface {
	SaveURL(shortCode, originalURL string) error
	GetURL(shortCode string) (string, error)
}

type InMemoryStorage struct {
	urls map[string]string
	mu   sync.RWMutex
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		urls: make(map[string]string),
	}
}

func (s *InMemoryStorage) SaveURL(shortCode, originalURL string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.urls[shortCode] = originalURL
	return nil
}

func (s *InMemoryStorage) GetURL(shortCode string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if url, exists := s.urls[shortCode]; exists {
		return url, nil
	}
	return "", errors.New("url not found")
}
