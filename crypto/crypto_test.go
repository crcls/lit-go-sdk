package crypto

import (
	"context"

	"github.com/crcls/lit-go-sdk/wasm"
)

type MockWasm struct {
	Context context.Context
}

func (mw MockWasm) Call(name string, args ...uint64) (interface{}, error) {
	switch name {
	case "combine_decryption_shares":
		return uint64(1), nil
	case "get_msg_byte":
		return uint64(1), nil
	case "get_rng_values_size":
		return uint64(16), nil
	case "encrypt":
		return uint64(16), nil
	case "get_ct_byte":
		return uint64(1), nil
	default:
		return nil, nil
	}
}

func (mw MockWasm) Close() {}

func mockNewWasmInstance(ctx context.Context) (wasm.Wasm, error) {
	return MockWasm{ctx}, nil
}
