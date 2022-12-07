package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/crcls/lit-go-sdk/auth"
	"github.com/crcls/lit-go-sdk/conditions"
	"github.com/crcls/lit-go-sdk/config"
)

type MockClient struct {
	client *ClientFactory
}

func (mc *MockClient) Connect() (bool, error) {
	mc.client.Ready = true
	return true, nil
}

func (mc *MockClient) GetEncryptionKey(params EncryptedKeyParams) ([]byte, error) {
	return mc.client.GetEncryptionKey(params)
}

func (mc *MockClient) NodeRequest(url string, body []byte) (*http.Response, error) {
	fmt.Printf("%s\n", url)
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
	} else if strings.HasSuffix(url, "/web/encryption/retrieve") {
		resp := &DecryptionShareResponse{
			DecryptionShare: "shareone",
			Result:          "success",
			ShareIndex:      1,
		}
		respBody, _ := json.Marshal(resp)
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

func NewMockClient() *MockClient {
	conf := config.New("localhost")
	conf.MinimumNodeCount = 1

	return &MockClient{
		client: &ClientFactory{
			Config: conf,
			ConnectedNodes: map[string]bool{
				"test://": true,
			},
			Ready:             false,
			ServerKeysForNode: make(map[string]ServerKeys),
		},
	}
}
