package client

import (
	"github.com/crcls/lit-go-sdk/config"
)

type Client struct {
	Config            *config.Config
	ConnectedNodes    []string
	Ready             bool
	ServerKeysForNode map[string]ServerKeys
	ServerPubKey      string
	SubnetPubKey      string
	NetworkPubKey     string
	NetworkPubKeySet  string
}

func New(c *config.Config) (*Client, error) {
	if c == nil {
		c = config.New(config.DEFAULT_NETWORK)
	}

	client := &Client{
		Config:            c,
		Ready:             false,
		ServerKeysForNode: make(map[string]ServerKeys),
	}

	if err := client.Connect(); err != nil {
		return nil, err
	}

	return client, nil
}
