package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/crcls/lit-go-sdk/auth"
)

type MockHttpClient struct {
	Response   string
	StatusCode int
}

func (mhc *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	body := ioutil.NopCloser(strings.NewReader(mhc.Response))
	status := mhc.StatusCode
	if status == 0 {
		status = 200
	}

	return &http.Response{
		Body:       body,
		Request:    req,
		Status:     strconv.Itoa(status),
		StatusCode: status,
	}, nil
}

func TestConnect(t *testing.T) {
	httpClient = &MockHttpClient{Response: testKeys}
	client := &Client{
		Config: testConfig,
	}

	if err := client.Connect(testctx); err != nil {
		t.Errorf("%+v", err)
	}
}

func TestConnectFail(t *testing.T) {
	httpClient = &MockHttpClient{Response: `{
		"result": "fail"
	}`}
	client := &Client{
		Config: testConfig,
	}

	if err := client.Connect(testctx); err == nil {
		t.Errorf("Expected an error from Connect")
	}
}

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

	// Simulate success response
	httpClient = &MockHttpClient{Response: `{
		"result": "success",
		"error": ""
	}`}

	ch := make(chan SaveCondMsg, 1)

	reqBody, _ := json.Marshal(&params)

	StoreEncryptionConditionWithNode(testctx, "http://localhost", "version", reqBody, ch)

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

func TestStoreEncryptionConditionWithNodeFailedRequest(t *testing.T) {
	httpClient = &MockHttpClient{Response: testKeys} // Needed for client.Connect

	// Simulate success response
	httpClient = &MockHttpClient{StatusCode: 500}

	ch := make(chan SaveCondMsg, 1)
	reqBody, _ := json.Marshal(&params)

	StoreEncryptionConditionWithNode(testctx, "http://localhost", "version", reqBody, ch)

	select {
	case msg := <-ch:
		if msg.Err == nil {
			t.Errorf("Expected Request failed error")
		}
	}
}

func TestStoreEncryptionConditionWithNodeUnexpectedResponse(t *testing.T) {
	httpClient = &MockHttpClient{Response: testKeys} // Needed for client.Connect

	// Simulate success response
	httpClient = &MockHttpClient{Response: ""}

	ch := make(chan SaveCondMsg, 1)
	reqBody, _ := json.Marshal(&params)

	StoreEncryptionConditionWithNode(testctx, "http://localhost", "version", reqBody, ch)

	select {
	case msg := <-ch:
		if msg.Err == nil {
			t.Errorf("Expected JSON unmarshal error")
		}
	}
}