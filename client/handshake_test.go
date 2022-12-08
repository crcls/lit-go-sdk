package client

import (
	"testing"

	"github.com/crcls/lit-go-sdk/config"
)

func TestHandshake(t *testing.T) {
	httpClient = &MockHttpClient{testKeys}
	c, _ := New(config.New("localhost"))
	ch := make(chan HnskMsg, 1)

	c.Handshake("/web/handshake", ch)

	select {
	case msg := <-ch:
		if !msg.Connected {
			t.Errorf("Handshake returned false connection")
		} else if msg.Keys == nil {
			t.Errorf("Handshake response keys are nil")
		} else if msg.Keys.NetworkPubKeySet != "networkPubKeySet" {
			t.Errorf("Unexpected NetworkPubKeySet key %s", msg.Keys.NetworkPubKeySet)
		}
	}
}
