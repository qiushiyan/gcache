package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/qiushiyan/gcache/pkg/cache"
	"github.com/qiushiyan/gcache/pkg/group"
	"github.com/qiushiyan/gcache/pkg/server"
	"github.com/qiushiyan/gcache/pkg/store"
)

type Response struct {
	Data  string  `json:"data"` // base 64 encoded
	Error *string `json:"error"`
}

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func createGroup() *group.Group {
	return group.New("scores", 2<<10, cache.LRU, store.GetterFunc(
		func(key store.Key) (store.Value, error) {
			if v, ok := db[string(key)]; ok {
				return store.NewByteView([]byte(v)), nil
			}
			return nil, fmt.Errorf("cannot find value for %s", key)
		}))
}

func startCacheServer(host string, peerAddrs []string, g *group.Group) {
	pool := server.NewPool(host)
	pool.AddPeer(peerAddrs...)
	g.RegisterPeerPicker(pool)
	log.Fatal(http.ListenAndServe(host[7:], pool))
}

func startAPIServer(apiAddr string, g *group.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			v, err := g.Get(store.Key(key))

			var response Response
			if err != nil {
				s := err.Error()
				response.Error = &s
			} else {
				if v != nil {
					view := v.(store.ByteView)
					response.Data = base64.StdEncoding.EncodeToString(view.ByteSlice())
				}
			}
			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(response)
			if err != nil {
				http.Error(w, fmt.Sprintf("error building the response, %v", err), http.StatusInternalServerError)
				return
			}
		}))
	log.Println("frontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))
}

func main() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "gcache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}
	host := addrMap[port]

	var peers = make([]string, len(addrMap))
	for _, v := range addrMap {
		peers = append(peers, v)
	}

	g := createGroup()
	if api {
		go startAPIServer(apiAddr, g)
	}
	startCacheServer(host, peers, g)
}
