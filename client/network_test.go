package client

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"testing"
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
		Config:            testConfig,
		Ready:             false,
		ServerKeysForNode: make(map[string]ServerKeys),
	}

	if err := client.Connect(); err != nil {
		t.Errorf("%+v", err)
	}
}

func TestConnectFail(t *testing.T) {
	httpClient = &MockHttpClient{Response: `{
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
