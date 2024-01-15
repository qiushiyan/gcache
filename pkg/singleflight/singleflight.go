package singleflight

import (
	"sync"

	"github.com/qiushiyan/gcache/pkg/store"
)

// a Call represents multiple pending requests for the same key
// the first request invokes the wait group, which halts any other requests
// val and err will be set upon completion of the first request, and all other requests can return the same value
type Call struct {
	wg  sync.WaitGroup
	val store.Value
	err error
}

// a CallGroup coordinates calls with the same key using Call
type CallGroup struct {
	mu sync.Mutex

	calls map[store.Key]*Call
}

func (g *CallGroup) Do(key store.Key, fn func() (store.Value, error)) (store.Value, error) {
	g.mu.Lock()
	g.lazyInitCalls()

	if c, ok := g.calls[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}

	c := new(Call)
	g.calls[key] = c
	g.mu.Unlock()

	c.wg.Add(1)

	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.calls, key)
	g.mu.Unlock()

	return c.val, c.err
}

func (g *CallGroup) lazyInitCalls() {
	if g.calls == nil {
		g.calls = make(map[store.Key]*Call)
	}
}
