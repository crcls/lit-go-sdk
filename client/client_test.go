package client

import "testing"

func init() {
	testResponse = testKeys
}

func TestNewDefaultConfig(t *testing.T) {
	c, err := New(nil)

	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	if c.Config.Network != "jalapeno" {
		t.Errorf("Unexpected network value: %s", c.Config.Network)
	}
}
