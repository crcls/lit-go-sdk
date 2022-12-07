package client

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/crcls/lit-go-sdk/config"
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

var testResponse string

func init() {
	httpClient = &MockHttpClient{testResponse}
}

func TestConnect(t *testing.T) {
	_, err := New(config.New("localhost"))

	if err != nil {
		t.Errorf("%+v", err)
	}
}
