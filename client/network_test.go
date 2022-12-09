package client

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

type MockHttpClient struct {
	Response string
}

func (mhc *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	body := ioutil.NopCloser(strings.NewReader(mhc.Response))

	return &http.Response{
		Body:    body,
		Request: req,
	}, nil
}

func TestConnect(t *testing.T) {
	client := &Client{
		Config:            testConfig,
		Ready:             false,
		ServerKeysForNode: make(map[string]ServerKeys),
	}

	if err := client.Connect(); err != nil {
		t.Errorf("%+v", err)
	}
}

func TestConnectFail(t *testing.T) {
	httpClient = &MockHttpClient{`{
		"result": "fail"
	}`}
	client := &Client{
		Config:            testConfig,
		Ready:             false,
		ServerKeysForNode: make(map[string]ServerKeys),
	}

	if err := client.Connect(); err == nil {
		t.Errorf("Expected an error from Connect")
	}
}
