package group

import (
	"sync"

	"github.com/qiushiyan/gcache/pkg/cache"
	"github.com/qiushiyan/gcache/pkg/store"
)

type Group struct {
	name      string
	getter    store.Getter // callback to get data in case of cache miss
	mainCache *cache.Cache
}

var (
	mu     = &sync.RWMutex{}
	groups = make(map[string]*Group)
)

func New(name string, cacheCap int64, getter store.Getter) *Group {
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache.New(cacheCap, cache.LRU),
	}

	mu.Lock()
	defer mu.Unlock()

	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()

	g := groups[name]
	return g
}

func (g *Group) Get(key store.Key) (store.Value, error) {
	if key.Empty() {
		return nil, store.ErrKeyEmpty
	}
	if v, ok := g.mainCache.Get(key); ok {
		return v, nil
	}

	return g.load(key)
}

func (g *Group) load(key store.Key) (store.Value, error) {
	value, err := g.getter.Get(key)
	if err != nil {
		return nil, store.ErrGetter.With(err)
	}

	g.mainCache.Set(key, value)
	return value, nil
}
