package group

import (
	"sync"

	"github.com/qiushiyan/gcache/pkg/cache"
	"github.com/qiushiyan/gcache/pkg/peer"
	"github.com/qiushiyan/gcache/pkg/store"
)

var once = &sync.Once{}

// A Group is a cache namespace and associated data loaded spread over
type Group struct {
	name       string
	getter     store.Getter // callback to get data in case of cache miss
	mainCache  *cache.Cache
	peerPicker peer.PeerPicker
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

func (g *Group) RegisterPeerPicker(picker peer.PeerPicker) {
	once.Do(func() {
		g.peerPicker = picker
	})
}

func (g *Group) Set(key store.Key, value store.Value) error {
	if key.Empty() {
		return store.ErrKeyEmpty
	}
	g.mainCache.Set(key, value)
	return nil
}

func (g *Group) Get(key store.Key) (store.Value, error) {
	if key.Empty() {
		return nil, store.ErrKeyEmpty
	}
	if v, ok := g.getLocal(key); ok {
		return v, nil
	}

	v, err := g.load(key)
	// if fetched from peer or getter successfully, update to main cache
	if err == nil {
		g.mainCache.Set(key, v)
	}

	return v, err
}

func (g *Group) load(key store.Key) (store.Value, error) {
	// fetch from peer nodes
	if g.peerPicker != nil {
		if client, ok := g.peerPicker.PickPeer(key); ok {
			v, err := g.getFromPeer(client, key)
			if err != nil {
				return nil, err
			}

			return v, nil
		}
	}

	if g.getter == nil {
		return nil, nil
	}
	// fetch from getter func
	value, err := g.getter.Get(key)
	if err != nil {
		return nil, store.ErrGetter.With(err)
	}

	g.mainCache.Set(key, value)
	return value, nil
}

func (g *Group) getFromPeer(client peer.PeerClient, key store.Key) (store.Value, error) {
	if bytes, err := client.Get(g.name, key); err != nil {
		return store.NewByteView(nil), nil
	} else {
		return store.NewByteView(bytes), nil
	}
}

func (g *Group) getLocal(key store.Key) (store.Value, bool) {
	return g.mainCache.Get(key)
}
