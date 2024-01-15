package peer

import (
	"github.com/qiushiyan/gcache/pkg/store"
)

type PeerClient interface {
	Get(group string, key store.Key) ([]byte, error)
}

type PeerPicker interface {
	Host() string
	PickPeer(key store.Key) (client PeerClient, ok bool)
}
