package client

import (
	"testing"
)

func TestHandshake(t *testing.T) {
	httpClient = &MockHttpClient{Response: testKeys}
	ch := make(chan HnskMsg, 1)

	Handshake(testctx, "http://localhost", "version", ch)

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

func TestHandshakeResponseMissingKeys(t *testing.T) {
	httpClient = &MockHttpClient{Response: `{"result": "fail"}`}
	ch := make(chan HnskMsg, 1)

	Handshake(testctx, "http://localhost", "version", ch)

	select {
	case msg := <-ch:
		if msg.Error == nil {
			t.Errorf("Expected to fail with Missing Key error")
		}
	}
}

func TestHandshakeFailedResponse(t *testing.T) {
	httpClient = &MockHttpClient{StatusCode: 500}
	ch := make(chan HnskMsg, 1)

	Handshake(testctx, "http://localhost", "version", ch)

	select {
	case msg := <-ch:
		if msg.Error == nil {
			t.Errorf("Expected to fail with Missing Key error")
		}
	}
}
