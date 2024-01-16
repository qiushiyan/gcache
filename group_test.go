package gcache

import (
	"fmt"
	"testing"

	"github.com/qiushiyan/gcache/pkg/cache"
	"github.com/qiushiyan/gcache/pkg/store"
)

var db1 = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestCreateGroup(t *testing.T) {
	name := "scores"
	NewGroup(name, 10, cache.LRU, nil)
	if g := GetGroup(name); g == nil {
		t.Fatal("create group failed")
	}
}

func TestGetGroup(t *testing.T) {
	name := "scores"
	NewGroup(name, 10, cache.LRU, nil)

	if g := GetGroup(name); g.name != name {
		t.Fatalf("get group failed, expect %s, get %s", name, g.name)
	}

}

func TestGroupGet(t *testing.T) {
	loadCounts := make(map[string]int, len(db1))
	g := NewGroup("scores", 2<<10, cache.LRU, store.GetterFunc(
		func(key store.Key) (store.Value, error) {
			if v, ok := db1[string(key)]; ok {
				if _, ok := loadCounts[string(key)]; !ok {
					loadCounts[string(key)] = 0
				}
				loadCounts[string(key)] += 1
				return store.NewByteViewFromStr(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	for k, v := range db1 {
		value, err := g.Get(store.Key(k))
		if err != nil {
			t.Fatalf("get key %s failed: %s", k, err)
		}
		view := value.(store.ByteView)
		if view.String() != v {
			t.Fatalf("get value failed, expect %s, get %s", v, view.String())
		}

		if loadCounts[k] > 1 {
			t.Fatalf("cache %s miss %d times", k, loadCounts[k])
		}

	}

	if view, err := g.Get("unknown"); err == nil {
		t.Fatalf("the value of unknown should be empty, but %s got", view)
	}
}
