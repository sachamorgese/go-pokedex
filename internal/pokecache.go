package internal

import (
	"sync"
	"time"
)

// cacheEntry should be a generic type
type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	cache map[string]cacheEntry
	mx    sync.RWMutex
}

func (c *Cache) Add(key string, val []byte) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.cache[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	if entry, ok := c.cache[key]; ok {
		return entry.val, true
	}

	return nil, false
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)

	for range ticker.C {
		for key, value := range c.cache {
			if value.createdAt.Add(interval).Before(time.Now()) {
				delete(c.cache, key)
			}
		}
	}
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		cache: make(map[string]cacheEntry),
	}

	go c.reapLoop(interval)

	return c
}
