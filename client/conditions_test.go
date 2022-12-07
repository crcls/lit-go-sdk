package client

import (
	"testing"

	"github.com/crcls/lit-go-sdk/auth"
	"github.com/crcls/lit-go-sdk/config"
)

func init() {
	testResponse = `{
		"result": "success",
		"error": ""
	}`
}

func TestStoreEncryptionConditionWithNode(t *testing.T) {
	c, _ := New(config.New("localhost"))
	ch := make(chan SaveCondMsg, 1)

	params := SaveCondParams{
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

	c.StoreEncryptionConditionWithNode("/test", params, ch)

	select {
	case msg := <-ch:
		if msg.Err != nil {
			t.Errorf("Unexpected error: %s", msg.Err.Error())
		}
	}
}
