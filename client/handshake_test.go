package client

import (
	"testing"
)

func TestHandshake(t *testing.T) {
	httpClient = &MockHttpClient{Response: testKeys}
	c := &Client{
		Config:            testConfig,
		Ready:             false,
		ServerKeysForNode: make(map[string]ServerKeys),
	}
	ch := make(chan HnskMsg, 1)

	c.Handshake(testctx, "http://localhost", ch)

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
	c := &Client{
		Config:            testConfig,
		Ready:             false,
		ServerKeysForNode: make(map[string]ServerKeys),
	}
	ch := make(chan HnskMsg, 1)

	c.Handshake(testctx, "http://localhost", ch)

	select {
	case msg := <-ch:
		if msg.Error == nil {
			t.Errorf("Expected to fail with Missing Key error")
		}
	}
}

func TestHandshakeFailedResponse(t *testing.T) {
	httpClient = &MockHttpClient{StatusCode: 500}
	c := &Client{
		Config:            testConfig,
		Ready:             false,
		ServerKeysForNode: make(map[string]ServerKeys),
	}
	ch := make(chan HnskMsg, 1)

	c.Handshake(testctx, "http://localhost", ch)

	select {
	case msg := <-ch:
		if msg.Error == nil {
			t.Errorf("Expected to fail with Missing Key error")
		}
	}
}
