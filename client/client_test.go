package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/crcls/lit-go-sdk/auth"
	"github.com/crcls/lit-go-sdk/conditions"
)

type MockClient struct{}

func (mc *MockClient) Connect() (bool, error) {
	return true, nil
}

func (mc *MockClient) GetEncryptionKey(params EncryptedKeyParams) ([]byte, error) {
	return make([]byte, 0), nil
}

func (mc *MockClient) NodeRequest(url string, body []byte) (*http.Response, error) {
	response := &http.Response{}

	if strings.HasSuffix(url, "/web/handshake") {
		keys := ServerKeys{
			ServerPubKey:     "ServerPubKey",
			SubnetPubKey:     "SubnetPubKey",
			NetworkPubKey:    "NetworkPubKey",
			NetworkPubKeySet: "NetworkPubKeySet",
		}
		respBody, _ := json.Marshal(keys)
		response.Body = io.NopCloser(bytes.NewBuffer(respBody))
	}

	return response, nil
}

func (mc *MockClient) SaveEncryptionKey(
	symmetricKey []byte,
	authSig auth.AuthSig,
	authConditions []conditions.EvmContractCondition,
	chain string,
	permanent bool,
) (string, error) {
	return "", nil
}
