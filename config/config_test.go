package config

import (
	"testing"
)

func TestNewDefault(t *testing.T) {
	c := New("localhost")

	if c.Network != "localhost" {
		t.Errorf("Unexpected network value: %s", c.Network)
	}

	if c.MinimumNodeCount != MINIMUM_NODE_COUNT {
		t.Errorf("Unexpected minimum node count value: %d", c.MinimumNodeCount)
	}

	if c.RequestTimeout != REQUEST_TIMEOUT {
		t.Errorf("Unexpected request timeout value: %d", c.RequestTimeout)
	}

	if c.Version != VERSION {
		t.Errorf("Unexpected version value: %s", c.Version)
	}
}
