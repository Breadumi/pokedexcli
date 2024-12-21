package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cachemap map[string]cacheEntry
	mu       sync.RWMutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache() Cache {
	return Cache{}
}
