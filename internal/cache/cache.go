package cache

import (
	"sort"
	"sync"
)

type URLCache struct {
	mu    sync.RWMutex
	cache map[string]*CacheEntry
}

type CacheEntry struct {
	OriginalURL string
	ShortURL    string
	Count       int
}

func NewURLCache() *URLCache {
	return &URLCache{
		cache: make(map[string]*CacheEntry),
	}
}

func (c *URLCache) AddEntry(originalURL, shortURL string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if entry, exists := c.cache[shortURL]; exists {
		entry.Count++
	} else {
		c.cache[shortURL] = &CacheEntry{
			OriginalURL: originalURL,
			ShortURL:    shortURL,
			Count:       1,
		}
	}
}

func (c *URLCache) GetEntry(shortURL string) (*CacheEntry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.cache[shortURL]
	return entry, exists
}

func (c *URLCache) DeleteEntry(shortURL string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.cache, shortURL)
}

func (c *URLCache) GetMostPopular(limit int) []*CacheEntry {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entries := make([]*CacheEntry, 0, len(c.cache))
	for _, entry := range c.cache {
		entries = append(entries, entry)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Count > entries[j].Count
	})

	if limit > len(entries) {
		limit = len(entries)
	}
	return entries[:limit]
}
