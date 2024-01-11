package cache

import (
	"sync"

	"github.com/qiushiyan/gcache/pkg/store"
)

type Cache struct {
	mu    sync.RWMutex
	store store.Store
	cap   int64
}

func New(cap int64) *Cache {
	return &Cache{
		cap: cap,
	}
}

func (c *Cache) initStore() {
	c.store = c.store.Create(c.cap, nil)
}

func (c *Cache) lazyInitStore() {
	if c.store == nil {
		c.initStore()
	}
}

func (c *Cache) get(key store.Key) (store.Value, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	c.lazyInitStore()
	return c.store.Get(key)
}

func (c *Cache) set(key store.Key, value store.Value) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.lazyInitStore()
	c.store.Set(key, value)
}
