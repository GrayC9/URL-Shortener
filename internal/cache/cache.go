package cache

import (
	"container/list"
	"log"
	"sort"
	"sync"
	"url_shortener/internal/storage"
)

type URLCache struct {
	mu       sync.Mutex
	cache    map[string]*list.Element
	lru      *list.List
	capacity int
}

type CacheEntry struct {
	OriginalURL string
	ShortURL    string
	Count       int
}

type entry struct {
	key   string
	value *CacheEntry
}

func NewURLCache(capacity int) *URLCache {
	return &URLCache{
		cache:    make(map[string]*list.Element),
		lru:      list.New(),
		capacity: capacity,
	}
}

func (c *URLCache) AddEntry(originalURL, shortURL string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.cache[originalURL]; ok {
		c.lru.MoveToFront(elem)
		return
	}

	newEntry := &CacheEntry{
		OriginalURL: originalURL,
		ShortURL:    shortURL,
		Count:       0,
	}

	elem := c.lru.PushFront(&entry{key: originalURL, value: newEntry})
	c.cache[originalURL] = elem

	if c.lru.Len() > c.capacity {
		c.removeOldest()
	}
}

func (c *URLCache) GetEntry(originalURL string) (*CacheEntry, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.cache[originalURL]; ok {
		c.lru.MoveToFront(elem)
		return elem.Value.(*entry).value, true
	}
	return nil, false
}

func (c *URLCache) IncrementCount(originalURL string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.cache[originalURL]; ok {
		elem.Value.(*entry).value.Count++
		c.lru.MoveToFront(elem)
		return true
	}
	return false
}

func (c *URLCache) DeleteEntry(originalURL string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.cache[originalURL]; ok {
		c.removeElement(elem)
	}
}

func (c *URLCache) GetMostPopular(limit int) []*CacheEntry {
	c.mu.Lock()
	defer c.mu.Unlock()

	entries := make([]*CacheEntry, 0, len(c.cache))
	for _, elem := range c.cache {
		entries = append(entries, elem.Value.(*entry).value)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Count > entries[j].Count
	})

	if limit > len(entries) {
		limit = len(entries)
	}
	return entries[:limit]
}

func (c *URLCache) removeOldest() {
	if elem := c.lru.Back(); elem != nil {
		c.removeElement(elem)
	}
}

func (c *URLCache) removeElement(elem *list.Element) {
	c.lru.Remove(elem)
	delete(c.cache, elem.Value.(*entry).key)
}

func PreloadCache(db storage.Storage, urlCache *URLCache) {
	popularURLs, err := db.GetPopularURLs(urlCache.capacity)
	if err != nil {
		log.Printf("Ошибка при получении популярных URL: %v", err)
		return
	}

	for _, url := range popularURLs {
		urlCache.AddEntry(url.OriginalURL, url.ShortCode)
		log.Printf("URL добавлен в кеш: %s -> %s", url.ShortCode, url.OriginalURL)
	}
}
