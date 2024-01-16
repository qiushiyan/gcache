package peer

import (
	pb "github.com/qiushiyan/gcache/pkg/gcachepb"
	"github.com/qiushiyan/gcache/pkg/store"
)

type PeerClient interface {
	Get(in *pb.Request, out *pb.Response) error
}

type PeerPicker interface {
	Host() string
	PickPeer(key store.Key) (client PeerClient, ok bool)
}
