package client

import (
	"fmt"
	"net/http"

	"github.com/crcls/lit-go-sdk/auth"
	"github.com/crcls/lit-go-sdk/conditions"
	"github.com/crcls/lit-go-sdk/config"
)

type Client interface {
	Connect() (bool, error)
	GetEncryptionKey(params EncryptedKeyParams) ([]byte, error)
	NodeRequest(url string, body []byte) (*http.Response, error)
	SaveEncryptionKey(
		symmetricKey []byte,
		authSig auth.AuthSig,
		authConditions []conditions.EvmContractCondition,
		chain string,
		permanent bool,
	) (string, error)
}

type ClientFactory struct {
	Config            *config.Config
	ConnectedNodes    map[string]bool
	Ready             bool
	ServerKeysForNode map[string]ServerKeys
	ServerPubKey      string
	SubnetPubKey      string
	NetworkPubKey     string
	NetworkPubKeySet  string
}

func New(c *config.Config) (*ClientFactory, error) {
	if c == nil {
		c = config.New(config.DEFAULT_NETWORK)
	}

	client := &ClientFactory{
		Config:            c,
		ConnectedNodes:    make(map[string]bool),
		Ready:             false,
		ServerKeysForNode: make(map[string]ServerKeys),
	}

	if ok, err := client.Connect(); !ok || err != nil {
		e := err

		if err == nil {
			e = fmt.Errorf("Client connect failed")
		}

		return nil, e
	}

	return client, nil
}
