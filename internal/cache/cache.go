package cache

import "C"
import (
	"log"
	"sort"
	"sync"
	"url_shortener/internal/storage"
)

const popularURLLimit = 1000

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
	var cacheNew = &CacheEntry{
		OriginalURL: originalURL,
		ShortURL:    shortURL,
		Count:       0,
	}

	c.mu.RLock()
	_, ok := c.cache[originalURL]
	c.mu.RUnlock()

	if !ok {
		c.mu.Lock()
		c.cache[originalURL] = cacheNew
		c.mu.Unlock()
	}
}

func (c *URLCache) GetEntry(originalURL string) (*CacheEntry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.cache[originalURL]
	return entry, ok
}

func (c *URLCache) IncrementCount(originalURL string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if entry, ok := c.cache[originalURL]; ok {
		entry.Count++
		return true
	}
	return false
}

func (c *URLCache) DeleteEntry(originalURL string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.cache, originalURL)
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

func PreloadCache(db storage.Storage, urlCache *URLCache) {
	popularURLs, err := db.GetPopularURLs(popularURLLimit)
	if err != nil {
		log.Printf("Ошибка при получении популярных URL: %v", err)
		return
	}

	for _, url := range popularURLs {
		urlCache.AddEntry(url.OriginalURL, url.ShortCode)
		log.Printf("URL добавлен в кеш: %s -> %s", url.ShortCode, url.OriginalURL)
	}
}
