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
		}
	}
}
