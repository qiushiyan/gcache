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

func List() []string {
	mu.RLock()
	defer mu.RUnlock()

	var names []string
	for k := range groups {
		names = append(names, k)
	}
	return names
}

func New(name string, cacheCap int64, cacheType cache.CacheType, getter store.Getter) *Group {
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache.New(cacheCap, cacheType),
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

func (g *Group) Set(key store.Key, value store.Value) error {
	if key.Empty() {
		return store.ErrKeyEmpty
	}
	g.mainCache.Set(key, value)
	return nil
}

func (g *Group) load(key store.Key) (store.Value, error) {
	if g.getter == nil {
		return nil, nil
	}
	value, err := g.getter.Get(key)
	if err != nil {
		return nil, store.ErrGetter.With(err)
	}

	g.mainCache.Set(key, value)
	return value, nil
}
