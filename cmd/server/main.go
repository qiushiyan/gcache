package main

import (
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

func main() {
	g := group.New("scores", 2<<10, cache.LRU, nil)
	for k, v := range db {
		val := store.NewByteViewFromStr(v)
		g.Set(store.Key(k), val)
		// value, err := g.Get(store.Key(k))
		// if err != nil {
		// 	fmt.Println(fmt.Sprintf("get key %s failed: %s", k, err))
		// 	return
		// }

		// view := value.(store.ByteView)
		// if view.String() != v {
		// 	fmt.Println(fmt.Sprintf("get value failed, expect %s, get %s", v, view.String()))
		// 	return
		// }
	}

	addr := "localhost:8080"
	p := server.NewPool(addr)
	http.ListenAndServe(addr, p)
}
