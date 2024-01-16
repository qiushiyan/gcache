package group

import (
	"log"
	"sync"

	"github.com/qiushiyan/gcache/pkg/cache"
	pb "github.com/qiushiyan/gcache/pkg/gcachepb"
	"github.com/qiushiyan/gcache/pkg/peer"
	"github.com/qiushiyan/gcache/pkg/singleflight"
	"github.com/qiushiyan/gcache/pkg/store"
)

var once = &sync.Once{}

// A Group is a cache namespace and associated data loaded spread over
type Group struct {
	name       string
	getter     store.Getter // callback to get data in case of cache miss
	mainCache  *cache.Cache
	peerPicker peer.PeerPicker
	loader     *singleflight.CallGroup
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
		loader:    &singleflight.CallGroup{},
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
	if v, ok := g.mainCache.Get(key); ok {
		g.log("cache hit for local cache")
		return v, nil
	}

	g.log("cache miss for local cache")

	v, err := g.load(key)

	if err != nil {
		g.log("returning error:", err)
	} else {
		g.log("returning value:", v)
	}
	return v, err
}

func (g *Group) load(key store.Key) (store.Value, error) {
	value, err := g.loader.Do(key, func() (store.Value, error) {
		if g.peerPicker != nil {
			if client, ok := g.peerPicker.PickPeer(key); ok {
				if v, err := g.getFromPeer(client, key); err == nil {
					// if fetched from peer or getter successfully, update to main cache
					if err == nil && v != nil {
						g.log("sync value to local cache")
						g.mainCache.Set(key, v)
					}
					return v, nil
				}
			}
		}

		// in any case peer is not available, fetch from getter func
		v, err := g.getFromGetter(key)
		// if fetched from peer or getter successfully, update to main cache
		if err == nil && v != nil {
			g.log("sync value to local cache")
			g.mainCache.Set(key, v)
		}

		return v, err

	})
	// fetch from peer nodes

	if err != nil {
		return nil, store.ErrGetter.With(err)
	}

	// update remote value to main cache
	g.mainCache.Set(key, value)
	return value, nil
}

func (g *Group) getFromPeer(client peer.PeerClient, key store.Key) (store.Value, error) {
	req := &pb.Request{
		Group: g.name,
		Key:   string(key),
	}
	res := &pb.Response{}
	if err := client.Get(req, res); err != nil {
		g.log("Failed to get from peer", err)
		return nil, err
	} else {
		return store.NewByteView(res.Value), nil
	}
}

func (g *Group) getFromGetter(key store.Key) (store.Value, error) {
	if g.getter == nil {
		g.log("getter is nil, returning")
		return nil, nil
	}

	g.log("fetch from getter func")
	value, err := g.getter.Get(key)
	return value, err
}

func (g *Group) log(text string, args ...any) {
	log.Println(g.peerPicker.Host(), text, args)
}
