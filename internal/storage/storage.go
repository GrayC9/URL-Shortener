package storage

import (
	"errors"
	"sync"
)

var ErrNotFound = errors.New("not found")

type MemoryStorage struct {
	mu     sync.RWMutex
	urls   map[string]string
	clicks map[string]int
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		urls:   make(map[string]string),
		clicks: make(map[string]int),
	}
}

func (m *MemoryStorage) GetShortCode(originalURL string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for k, v := range m.urls {
		if v == originalURL {
			return k, nil
		}
	}
	return "", ErrNotFound
}

func (m *MemoryStorage) SaveURL(shortCode, originalURL string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.urls[shortCode] = originalURL
	return nil
}

func (m *MemoryStorage) GetURL(shortCode string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	url, ok := m.urls[shortCode]
	if !ok {
		return "", ErrNotFound
	}
	return url, nil
}

func (m *MemoryStorage) IncrementClickCount(shortCode string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.clicks[shortCode]++
	return nil
}

func (m *MemoryStorage) UpdateLastAccessed(shortCode string) error {
	return nil
}
