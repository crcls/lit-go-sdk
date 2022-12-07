package client

import (
	"testing"

	"github.com/crcls/lit-go-sdk/config"
)

var keys = `{
	"serverPublicKey": "ServerPubKey",
	"subnetPublicKey": "SubnetPubKey",
	"networkPublicKey": "NetworkPubKey",
	"networkPublicKeySet": "NetworkPubKeySet"
}`

func init() {
	testResponse = keys
}

func TestHandshake(t *testing.T) {
	c, _ := New(config.New("localhost"))
	ch := make(chan HnskMsg, 1)

	c.Handshake("/web/handshake", ch)

	select {
	case msg := <-ch:
		if !msg.Connected {
			t.Errorf("Handshake returned false connection")
		} else if msg.Keys == nil {
			t.Errorf("Handshake response keys are nil")
		} else if msg.Keys.NetworkPubKeySet != "NetworkPubKeySet" {
			t.Errorf("Unexpected NetworkPubKeySet key %s", msg.Keys.NetworkPubKeySet)
		}
	}
}
