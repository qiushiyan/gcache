package consistenthash

import (
	"hash/crc32"
	"slices"
	"sort"
	"strconv"
	"sync"
)

// Hash maps bytes to uint32
type HashFunc func(data []byte) uint32

// A Ring is a hash ring of virtual nodes
type Ring struct {
	hashFunc HashFunc
	replicas int            // virtual node amplification factor
	nodes    []int          // sorted virtual nodes
	nodeMap  map[int]string // map virtual node to real node name
	mu       sync.RWMutex
}

func New(replicas int, fn HashFunc) *Ring {
	r := &Ring{
		replicas: replicas,
		hashFunc: fn,
		nodeMap:  make(map[int]string),
	}
	if r.hashFunc == nil {
		// default hash func
		r.hashFunc = crc32.ChecksumIEEE
	}
	return r
}

func (r *Ring) Add(keys ...string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i := range keys {
		key := keys[i]
		for j := 0; j < r.replicas; j++ {
			hash := r.hashWithReplica(key, j)
			r.nodes = append(r.nodes, hash)
			r.nodeMap[hash] = key
		}
	}

	slices.Sort(r.nodes)
}

func (r *Ring) Get(key string) string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if len(r.nodes) == 0 {
		return ""
	}

	hash := int(r.hashFunc([]byte(key)))

	idx := sort.SearchInts(r.nodes, hash)

	vnode := r.nodes[idx%len(r.nodes)]
	return r.nodeMap[vnode]
}

func (r *Ring) Remove(keys ...string) {
	for i := range keys {
		k := keys[i]
		for j := 0; j < r.replicas; j++ {
			hash := r.hashWithReplica(k, j)
			r.deleteVNode(hash)
		}
	}
}
func (r *Ring) hashWithReplica(key string, replica int) int {
	return int(r.hashFunc([]byte(key + strconv.Itoa(replica))))
}

func (r *Ring) deleteVNode(hash int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	idx := sort.SearchInts(r.nodes, hash)
	if idx < len(r.nodes) && r.nodes[idx] == hash {
		copy(r.nodes[idx:], r.nodes[idx+1:])
		r.nodes[len(r.nodes)-1] = 0
		r.nodes = r.nodes[:len(r.nodes)-1]
		delete(r.nodeMap, hash)
	}
}
