package lru

import (
	"container/list"
	"log"

	"github.com/qiushiyan/gcache/pkg/store"
	"github.com/qiushiyan/gcache/pkg/strategy"
)

type Entry struct {
	key   store.Key
	value store.Value
}

type Cache struct {
	cap       int64 // max number of bytes
	size      int64 // current size
	ll        *list.List
	cache     map[store.Key]*list.Element
	OnEvicted store.EvictedCallback
}

func New(cap int64, cb store.EvictedCallback) *Cache {
	return &Cache{
		cap:       cap,
		ll:        list.New(),
		cache:     make(map[store.Key]*list.Element),
		OnEvicted: cb,
	}
}

func (c *Cache) Create(cap int64, cb store.EvictedCallback) *Cache {
	return New(cap, cb)
}

func (c *Cache) Get(key store.Key) (store.Value, bool) {
	el, ok := c.cache[key]
	if !ok {
		return nil, false
	}

	c.ll.MoveToFront(el)

	entry := toEntry(el)
	return entry.value, true
}

func (c *Cache) Evict(strata strategy.EvictType) {
	switch strata {
	case strategy.LRU:
		el := c.ll.Back()
		if el != nil {
			entry := toEntry(el)
			c.ll.Remove(el)
			delete(c.cache, entry.key)

			c.size -= int64(entry.value.Len() + len(entry.key))

			if c.OnEvicted != nil {
				c.OnEvicted(entry.key, entry.value)
			}
		}

	default:
		log.Fatalf("unknown eviction strategy %v", strata)
	}
}

func toEntry(el *list.Element) *Entry {
	entry, ok := el.Value.(*Entry)
	if !ok {
		log.Fatal("cache entry is not of type *Entry")
	}
	return entry
}

func (c *Cache) Set(key store.Key, v store.Value) {
	if el, ok := c.cache[key]; ok {
		c.ll.MoveToFront(el)
		entry := toEntry(el)
		// compare size between new and old value
		c.size += int64(v.Len() - entry.value.Len())
		entry.value = v
	} else {
		c.cache[key] = c.ll.PushFront(&Entry{key, v})
		// add size of the new value
		c.size += int64(len(key) + v.Len())
	}

	for c.cap != 0 && c.size > c.cap {
		c.Evict(strategy.LRU)
	}
}

// Removes the key from cache
// Returns true if the key was removed
func (c *Cache) Delete(key store.Key) bool {
	if el, ok := c.cache[key]; ok {
		c.ll.Remove(el)
		entry := toEntry(el)
		delete(c.cache, entry.key)
		c.size -= int64(len(entry.key) + entry.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(entry.key, entry.value)
		}

		return true
	}

	return false
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
