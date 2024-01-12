package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/qiushiyan/gcache/pkg/cache"
	"github.com/qiushiyan/gcache/pkg/group"
	"github.com/qiushiyan/gcache/pkg/store"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestServer(t *testing.T) {
	g := group.New("scores", 2<<10, cache.LRU, nil)

	addr := "http://localhost:8080"
	p := NewPool(addr)
	for k, v := range db {
		g.Set(store.Key(k), store.NewByteViewFromStr(v))
		request, _ := http.NewRequest(http.MethodGet, addr+"/_gcache/scores/"+k, nil)
		response := httptest.NewRecorder()

		p.ServeHTTP(response, request)
		v := response.Body.String()
		if v != db[k] {
			t.Errorf("get value failed, expect %s, get %s", db[k], v)
		}
	}

}
