package client

import (
	"testing"

	"github.com/crcls/lit-go-sdk/config"
)

var testConfig = &config.Config{
	MinimumNodeCount: 1,
	Network:          "localhost",
	RequestTimeout:   config.REQUEST_TIMEOUT,
	Version:          config.VERSION,
}

func TestNewWithDefaultConfig(t *testing.T) {
	httpClient = &MockHttpClient{testKeys}
	c, err := New(nil)

	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	if c.Config.Network != "jalapeno" {
		t.Errorf("Unexpected network value: %s", c.Config.Network)
	}
}

func TestNewWithConfig(t *testing.T) {
	httpClient = &MockHttpClient{testKeys}
	c, err := New(config.New("localhost"))

	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	if c.Config.Network != "localhost" {
		t.Errorf("Unexpected network value: %s", c.Config.Network)
	}
}

func TestNewFailConnect(t *testing.T) {
	httpClient = &MockHttpClient{"500"}
	_, err := New(testConfig)

	if err == nil {
		t.Errorf("Expected an error when the client connects")
	}
}
