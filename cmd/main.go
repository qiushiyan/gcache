package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/qiushiyan/gcache/pkg/cache"
	"github.com/qiushiyan/gcache/pkg/group"
	"github.com/qiushiyan/gcache/pkg/server"
	"github.com/qiushiyan/gcache/pkg/store"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func createGroup() *group.Group {
	return group.New("scores", 2<<10, cache.LRU, store.GetterFunc(
		func(key store.Key) (store.Value, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[string(key)]; ok {
				return store.NewByteView([]byte(v)), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
}

func startCacheServer(host string, peerAddrs []string, g *group.Group) {
	pool := server.NewPool(host)
	pool.AddPeer(peerAddrs...)
	g.RegisterPeerPicker(pool)
	http.ListenAndServe(host, pool)
}

func startAPIServer(apiAddr string, g *group.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			v, err := g.Get(store.Key(key))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			view := v.(*store.ByteView)
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Write(view.ByteSlice())

		}))
	log.Println("frontend server is running at", apiAddr)
	http.ListenAndServe(apiAddr, nil)
}

func main() {
	port := 8001

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	g := createGroup()
	go startAPIServer(apiAddr, g)
	startCacheServer(addrMap[port], []string(addrs), g)
}
