package client

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/crcls/lit-go-sdk/auth"
	"github.com/crcls/lit-go-sdk/conditions"
	"github.com/crcls/lit-go-sdk/config"
	"github.com/crcls/lit-go-sdk/crypto"
)

var testKeys = `{
	"serverPublicKey": "serverPubKey",
	"subnetPublicKey": "subnetPubKey",
	"networkPublicKey": "networkPubKey",
	"networkPublicKeySet": "networkPubKeySet"
}`

func TestServerKeysKeys(t *testing.T) {
	sk := ServerKeys{
		"serverPubKey",
		"subnetPubKey",
		"networkPubKey",
		"networkPubKeySet",
	}

	if spk, ok := sk.Key("ServerPubKey"); !ok || spk != "serverPubKey" {
		t.Errorf("Unexpected key for `ServerPubKey`: %s", spk)
	}

	if sbk, ok := sk.Key("SubnetPubKey"); !ok || sbk != "subnetPubKey" {
		t.Errorf("Unexpected key for `SubnetPubKey`: %s", sbk)
	}

	if npk, ok := sk.Key("NetworkPubKey"); !ok || npk != "networkPubKey" {
		t.Errorf("Unexpected key for `NetworkPubKey`: %s", npk)
	}

	if npks, ok := sk.Key("NetworkPubKeySet"); !ok || npks != "networkPubKeySet" {
		t.Errorf("Unexpected key for `NetworkPubKeySet`: %s", npks)
	}

	if _, ok := sk.Key("test"); ok {
		t.Errorf("Expected Key to return 'not ok' with unsupported key")
	}
}

func TestMostCommonKey(t *testing.T) {
	client := &Client{
		ServerKeysForNode: map[string]ServerKeys{
			"http://localhost:7470": ServerKeys{
				ServerPubKey: "common",
			},
			"http://localhost:7471": ServerKeys{
				ServerPubKey: "common",
			},
			"http://localhost:7472": ServerKeys{
				ServerPubKey: "uncommon",
			},
		},
	}

	key := client.MostCommonKey("ServerPubKey")
	if key != "common" {
		t.Errorf("Unexpected result from MostCommonKey: expected `common` got %s", key)
	}
}

var testAuthSig = auth.AuthSig{
	Sig:           "sig",
	DerivedVia:    "derivedVia",
	SignedMessage: "signedMessage",
	Address:       "address",
}

var testCondition = conditions.EvmContractCondition{
	ContractAddress: "contractAddress",
	FunctionName:    "functionName",
	FunctionParams:  []string{"functionParams"},
	FunctionAbi: conditions.AbiMember{
		Name:            "name",
		Inputs:          []conditions.AbiIO{},
		Outputs:         []conditions.AbiIO{},
		Constant:        false,
		StateMutability: "stateMutability",
	},
}

var testParams = EncryptedKeyParams{
	AuthSig: &testAuthSig,
	Chain:   "ethereum",
	EvmContractConditions: []*conditions.EvmContractCondition{
		&testCondition,
	},
	ToDecrypt: "toDecrypt",
}

func mockThresholdDecrypt(shares []crypto.DecryptionShare, ciphertext, netPubKeySet string) ([]byte, error) {
	return []byte("test"), nil
}

func TestGetEncryptionKey(t *testing.T) {
	httpClient = &MockHttpClient{Response: testKeys}
	c, _ := New(config.New("localhost"))

	// Mock the Threshold Decrypt function
	thresholdDecrypt = mockThresholdDecrypt

	// Mock the network response
	httpClient = &MockHttpClient{Response: `{
		"decryptionShare": "share",
		"result": "success"
	}`}

	key, err := c.GetEncryptionKey(testParams)
	if err != nil {
		t.Errorf(err.Error())
	}

	if string(key) != "test" {
		t.Errorf("Unexpected value for key: %s", string(key))
	}
}

func TestGetEncryptionKeyClientNotReady(t *testing.T) {
	c := &Client{
		Config:            testConfig,
		Ready:             false,
		ServerKeysForNode: make(map[string]ServerKeys),
	}

	_, err := c.GetEncryptionKey(testParams)

	if err == nil {
		t.Errorf("Expected a client not ready error")
	}
}

func TestGetDecryptionShare(t *testing.T) {
	httpClient = &MockHttpClient{Response: testKeys}
	c, _ := New(config.New("localhost"))

	// Mock the network response
	httpClient = &MockHttpClient{StatusCode: 500}
	ch := make(chan DecryptResMsg, 1)

	c.GetDecryptionShare("http://localhost", testParams, ch)

	select {
	case msg := <-ch:
		if msg.Err == nil {
			t.Errorf("Expected a request failed error")
		}
	}
}

func mockThresholdEncrypt(subPubKey []byte, message []byte) ([]byte, error) {
	if string(message) == "privateKey" {
		return []byte("encryptedKey"), nil
	} else {
		return []byte{}, fmt.Errorf("failed")
	}
}

func TestSaveEncryptionKey(t *testing.T) {
	h := hex.EncodeToString([]byte("subnetPubKey"))
	keys := fmt.Sprintf(`{
		"serverPublicKey": "serverPubKey",
		"subnetPublicKey": "%s",
		"networkPublicKey": "networkPubKey",
		"networkPublicKeySet": "networkPubKeySet"
	}`, h)
	httpClient = &MockHttpClient{Response: keys}
	c, _ := New(config.New("localhost"))

	thresholdEncrypt = mockThresholdEncrypt

	encryptedKey, err := c.SaveEncryptionKey(
		[]byte("privateKey"),
		testAuthSig,
		[]conditions.EvmContractCondition{testCondition},
		"ethereum",
		false,
	)

	if err != nil {
		t.Errorf(err.Error())
	}

	if encryptedKey == "" {
		t.Errorf("Expected a value for encryptedKey")
	}
}
