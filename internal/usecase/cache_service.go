package usecase

import (
	"sync"
	"time"
)

type CacheItem struct {
	Value      []byte
	Expiration int64
}

type CacheService struct {
	mu    sync.RWMutex
	items map[string]CacheItem
}

func NewCacheService() *CacheService {
	return &CacheService{
		items: make(map[string]CacheItem),
	}
}

func (c *CacheService) Set(key string, value []byte, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = CacheItem{
		Value:      value,
		Expiration: time.Now().Add(ttl).UnixNano(),
	}
}

func (c *CacheService) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	item, found := c.items[key]
	c.mu.RUnlock()

	if !found {
		return nil, false
	}

	if time.Now().UnixNano() > item.Expiration {
		c.mu.Lock()
		delete(c.items, key)
		c.mu.Unlock()
		return nil, false
	}

	return item.Value, true
}

func (c *CacheService) Invalidate(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

func (c *CacheService) ClearAll() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]CacheItem)
}
