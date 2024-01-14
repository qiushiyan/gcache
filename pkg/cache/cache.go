package cache

import (
	"sync"

	"github.com/qiushiyan/gcache/pkg/store"
	"github.com/qiushiyan/gcache/pkg/store/lru"
)

type Cache struct {
	mu        *sync.RWMutex
	store     store.Store
	cap       int64
	cacheType CacheType
}

func New(cap int64, cacheType CacheType) *Cache {
	return &Cache{
		cap:       cap,
		cacheType: cacheType,
	}
}

func (c *Cache) initStore() {
	switch c.cacheType {
	case LRU:
		c.store = lru.New(c.cap, nil)
	}
}

func (c *Cache) lazyInitStore() {
	if c.store == nil {
		c.initStore()
	}
}

func (c *Cache) Get(key store.Key) (store.Value, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	c.lazyInitStore()
	return c.store.Get(key)
}

func (c *Cache) Set(key store.Key, value store.Value) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.lazyInitStore()
	c.store.Set(key, value)
}
