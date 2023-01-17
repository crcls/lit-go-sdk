package client

import (
	"context"

	"github.com/crcls/lit-go-sdk/config"
)

type Client struct {
	Config           *config.Config
	ConnectedNodes   []string
	Ready            bool
	ServerKeys       []ServerKeys
	ServerPubKey     string
	SubnetPubKey     string
	NetworkPubKey    string
	NetworkPubKeySet string
}

func New(ctx context.Context, c *config.Config) (*Client, error) {
	if c == nil {
		c = config.New(config.DEFAULT_NETWORK)
	}

	client := &Client{
		Config: c,
	}

	if err := client.Connect(ctx); err != nil {
		return nil, err
	}

	return client, nil
}
