package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/qiushiyan/gcache/pkg/group"
	"github.com/qiushiyan/gcache/pkg/store"
)

var defaultPath = "/_gcache/"

type Pool struct {
	host string
	path string
}

func NewPool(host string) *Pool {
	return &Pool{
		host: host,
		path: defaultPath,
	}
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

		if v, err := g.Get(store.Key(key)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			if v == nil {
				w.Header().Set("Content-Type", "text/plain")
				w.Write([]byte("nil"))
			} else {
				w.Header().Set("Content-Type", "application/octet-stream")
				view := v.(store.ByteView)
				w.Write(view.ByteSlice())

			}
		}

	} else {
		http.Error(w, "Pool serving unexpected path: "+r.URL.Path, http.StatusBadRequest)
		return
	}

}

// check if path is in the format /<basepath>/<groupname>/<key>
// returns [groupname, key] and true if path is valid
func (p *Pool) validateRequest(r *http.Request) ([]string, bool) {
	if !strings.HasPrefix(r.URL.Path, p.path) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	parts := strings.SplitN(r.URL.Path[len(p.path):], "/", 2)
	if len(parts) != 2 {
		return nil, false
	}

	return parts, true
}
