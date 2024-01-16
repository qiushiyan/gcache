package client

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	pb "github.com/qiushiyan/gcache/pkg/gcachepb"
	"github.com/qiushiyan/gcache/pkg/peer"
	"google.golang.org/protobuf/proto"
)

// A client is associated with a peer and request the server given group name and key
// implements peer.PeerClient
type Client struct {
	baseURL string
}

func New(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
	}
}

func (c *Client) Get(in *pb.Request, out *pb.Response) error {
	u := fmt.Sprintf(
		"%v/_gcache/%v/%v",
		c.baseURL,
		url.QueryEscape(in.GetGroup()),
		url.QueryEscape(in.GetKey()),
	)
	resp, err := http.Get(u)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned: %v", resp.Status)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %v", err)
	}

	if err = proto.Unmarshal(bytes, out); err != nil {
		return fmt.Errorf("decoding response body: %v", err)
	}

	return nil
}

var _ peer.PeerClient = (*Client)(nil)
