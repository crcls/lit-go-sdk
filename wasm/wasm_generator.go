//go:generate go-bindata -pkg wasm -o wasm_base64.go .

package wasm

import (
	"bytes"
	"io/ioutil"
)

func WasmBin() ([]byte, error) {
	data, err := Asset("./threshold_crypto_wasm_bridge_bg.wasm")
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(data)
	return ioutil.ReadAll(buf)
}
