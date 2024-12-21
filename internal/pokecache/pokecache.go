package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	CacheMap map[string]cacheEntry
	mu       *sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		CacheMap: make(map[string]cacheEntry),
		mu:       &sync.Mutex{},
	}
	go c.reapLoop(interval)
	return c
}

func (c *Cache) Add(key string, v []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.CacheMap[key] = cacheEntry{
		val: v,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, found := c.CacheMap[key]
	return v.val, found
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.mu.Lock()
		t := time.Now()
		for key, val := range c.CacheMap {
			if val.createdAt.Before(t.Add(-interval)) {
				delete(c.CacheMap, key)
			}
		}
		c.mu.Unlock()
	}
}
