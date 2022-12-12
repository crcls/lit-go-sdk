package client

import (
	"testing"

	"github.com/crcls/lit-go-sdk/auth"
	"github.com/crcls/lit-go-sdk/config"
)

var params = SaveCondParams{
	Key: "key",
	Val: "val",
	AuthSig: auth.AuthSig{
		Sig:           "sig",
		DerivedVia:    "derivedVia",
		SignedMessage: "signedMessage",
		Address:       "address",
	},
	Chain:     "chain",
	Permanent: 0,
}

func TestStoreEncryptionConditionWithNode(t *testing.T) {
	httpClient = &MockHttpClient{Response: testKeys} // Needed for client.Connect
	c, _ := New(config.New("localhost"))

	// Simulate success response
	httpClient = &MockHttpClient{Response: `{
		"result": "success",
		"error": ""
	}`}

	ch := make(chan SaveCondMsg, 1)

	c.StoreEncryptionConditionWithNode("/test", params, ch)

	select {
	case msg := <-ch:
		if msg.Err != nil {
			t.Errorf("Unexpected error: %s", msg.Err.Error())
		} else if msg.Response.Result != "success" {
			t.Errorf("Unexpected response result: %s", msg.Response.Result)
		} else if msg.Response.Error != "" {
			t.Errorf("Unexpected response error: %s", msg.Response.Error)
		}
	}
}
