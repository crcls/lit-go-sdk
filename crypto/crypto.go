package crypto

import (
	"context"

	"github.com/crcls/lit-go-sdk/wasm"
)

var newWasmInstance func(ctx context.Context) (wasm.Wasm, error)

func init() {
	newWasmInstance = wasm.NewWasmInstance
}
