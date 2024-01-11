package group

import (
	"sync"

	"github.com/qiushiyan/gcache/pkg/cache"
	"github.com/qiushiyan/gcache/pkg/store"
)

type Group struct {
	name      string
	getter    store.Getter // callback to get data in case of cache miss
	mainCache cache.Cache
}

var (
	mu     = &sync.RWMutex{}
	groups = make(map[string]*Group)
)

func New(name string, cacheCap int64, getter store.Getter) *Group {
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: *cache.New(cacheCap),
	}

	mu.Lock()
	defer mu.Unlock()

	groups[name] = g
	return g
}

func Get(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()

	g := groups[name]
	return g
}
