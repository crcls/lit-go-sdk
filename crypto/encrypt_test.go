package crypto

import (
	"bytes"
	"context"
	"reflect"
	"testing"
)

func TestThresholdEncrypt(t *testing.T) {
	newWasmInstance = mockNewWasmInstance

	result, err := ThresholdEncrypt(context.Background(), []byte("subPubKey"), []byte("message"))

	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	expected := bytes.Repeat([]byte{byte(1)}, 16)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Unexpected result: %v", result)
	}
}
