package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"

	"github.com/qiushiyan/gcache/pkg/client"
	"github.com/qiushiyan/gcache/pkg/consistenthash"
	pb "github.com/qiushiyan/gcache/pkg/gcachepb"
	"github.com/qiushiyan/gcache/pkg/group"
	"github.com/qiushiyan/gcache/pkg/peer"
	"github.com/qiushiyan/gcache/pkg/store"
	"google.golang.org/protobuf/proto"
)

const (
	defaultBasePath = "/_gcache/"
	defaultReplicas = 10
)

type Pool struct {
	host    string
	path    string
	mu      sync.RWMutex
	peers   *consistenthash.Ring
	clients map[string]peer.PeerClient // keyed by peer host, e.g. "http://10.0.0.2:8008"
}

type Response struct {
	Data  string  `json:"data"` // base 64 encoded
	Error *string `json:"error"`
}

func NewPool(host string) *Pool {
	return &Pool{
		host:    host,
		path:    defaultBasePath,
		peers:   consistenthash.New(defaultReplicas, nil),
		clients: make(map[string]peer.PeerClient),
	}
}

func (p *Pool) Host() string {
	return p.host
}

func (p *Pool) log(text string, args ...any) {
	slog.Info(fmt.Sprintf("[Server %s] %s", p.host, text), args...)
}

func (p *Pool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if parts, ok := p.validateRequest(r); ok {
		p.log("", "method", r.Method, "path", r.URL.Path)
		groupName, key := parts[0], parts[1]

		g := group.GetGroup(groupName)
		if g == nil {
			http.Error(w, "no such group: "+groupName, http.StatusNotFound)
			http.Error(w, "available groups: "+strings.Join(group.List(), ", "), http.StatusNotFound)
			return
		}

		v, err := g.Get(store.Key(key))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		view := v.(store.ByteView)
		body, err := proto.Marshal(&pb.Response{Value: view.ByteSlice()})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(body)

	} else {
		http.Error(w, "Pool serving unexpected path: "+r.URL.Path, http.StatusBadRequest)
		return
	}

}

func (p *Pool) AddPeer(peerAddrs ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, addr := range peerAddrs {
		p.peers.Add(addr)
		p.clients[addr] = client.New(addr)
	}

}

// implements peer.PeerPicker
func (p *Pool) PickPeer(key store.Key) (peer.PeerClient, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	peer := p.peers.Get(string(key))
	if peer == "" {
		p.log("peer is empty")
	} else if peer == p.host {
		p.log("peer is self")
	} else {
		p.log(fmt.Sprintf("%s pick peer %s", p.host, peer))
		return p.clients[peer], true
	}

	return nil, false

}

// check if path is in the format /<basepath>/<groupname>/<key>
// parts should be [groupname, key]
func (p *Pool) validateRequest(r *http.Request) (parts []string, ok bool) {
	if !strings.HasPrefix(r.URL.Path, p.path) {
		return nil, false
	}
	parts = strings.SplitN(r.URL.Path[len(p.path):], "/", 2)
	if len(parts) != 2 {
		return nil, false
	}

	return parts, true
}

var _ peer.PeerPicker = (*Pool)(nil)
