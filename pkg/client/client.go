package client

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/qiushiyan/gcache/pkg/peer"
	"github.com/qiushiyan/gcache/pkg/store"
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

func (c *Client) Get(group string, key store.Key) ([]byte, error) {
	u := fmt.Sprintf(
		"%v%v/%v",
		c.baseURL,
		url.QueryEscape(group),
		url.QueryEscape(string(key)),
	)

	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned: %v", resp.Status)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}

	return bytes, nil
}

var _ peer.PeerClient = (*Client)(nil)
