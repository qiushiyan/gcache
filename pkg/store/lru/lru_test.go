package lru

import (
	"testing"

	"github.com/qiushiyan/gcache/pkg/store"
)

type String string

func (d String) Len() int {
	return len(d)
}

func TestSet(t *testing.T) {
	c := New(20, nil)
	c.Set("hello", String("world"))
	c.Set("hello", String("squirrel"))

	if c.Len() != 1 {
		t.Fatalf("expected size to be 1, got %d", c.Len())
	}

}

func TestGet(t *testing.T) {
	c := New(int64(0), nil)
	c.Set("key1", String("1234"))
	if v, ok := c.Get("key1"); !ok || string(v.(String)) != "1234" {
		t.Fatalf("cache hit key1=1234 failed")
	}
	if _, ok := c.Get("key2"); ok {
		t.Fatalf("cache miss key2 failed")
	}
}

func TestAutoEvict(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "k3"
	v1, v2, v3 := "value1", "value2", "v3"
	cap := len(k1 + k2 + v1 + v2)
	c := New(int64(cap), nil)

	c.Set(store.Key(k1), String(v1))
	c.Set(store.Key(k2), String(v2))
	c.Set(store.Key(k3), String(v3))

	if _, ok := c.Get("key1"); ok || c.Len() != 2 {
		t.Fatalf("RemoveAutoEvict key1 failed")
	}
}

func TestDelete(t *testing.T) {
	c := New(int64(20), nil)
	c.Set("key1", String("1234"))
	c.Delete("key1")

	if _, ok := c.Get("key1"); ok {
		t.Fatalf("delete key1 failed")
	}
}
