package client

import (
	"context"
	"testing"

	"github.com/crcls/lit-go-sdk/config"
)

var testConfig = &config.Config{
	MinimumNodeCount: 1,
	Network:          "localhost",
	RequestTimeout:   config.REQUEST_TIMEOUT,
	Version:          config.VERSION,
}

var testctx = context.Background()

func TestNewWithDefaultConfig(t *testing.T) {
	httpClient = &MockHttpClient{Response: testKeys}
	c, err := New(testctx, nil)

	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	if c.Config.Network != "jalapeno" {
		t.Errorf("Unexpected network value: %s", c.Config.Network)
	}
}

func TestNewWithConfig(t *testing.T) {
	httpClient = &MockHttpClient{Response: testKeys}
	c, err := New(testctx, config.New("localhost"))

	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	if c.Config.Network != "localhost" {
		t.Errorf("Unexpected network value: %s", c.Config.Network)
	}
}

func TestNewFailConnect(t *testing.T) {
	httpClient = &MockHttpClient{StatusCode: 500}
	_, err := New(testctx, testConfig)

	if err == nil {
		t.Errorf("Expected an error when the client connects")
	}
}
