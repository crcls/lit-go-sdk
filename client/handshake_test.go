package client

import (
	"testing"
)

func TestHandshake(t *testing.T) {
	c := &MockClient{}
	ch := make(chan HnskMsg, 1)

	Handshake("/web/handshake", ch, c.NodeRequest)

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
