package crypto

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

func PKCS7Padding(plaintext []byte) []byte {
	padding := aes.BlockSize - (len(plaintext) % aes.BlockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(plaintext, padtext...)
}

func Prng(length uint64) []byte {
	values := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, values); err != nil {
		panic(err)
	}

	return values
}

func AesEncrypt(key []byte, plaintext []byte) (ciphertext []byte) {
	padded := PKCS7Padding(plaintext)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	ciphertext = make([]byte, aes.BlockSize+len(padded))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], padded)

	return
}

func ThresholdEncrypt(ctx context.Context, subPubKey []byte, message []byte) ([]byte, error) {
	wasm, err := newWasmInstance(ctx)
	if err != nil {
		return nil, err
	}

	rngSize, err := wasm.Call("get_rng_values_size")
	if err != nil {
		return nil, err
	}

	rngValues := Prng(rngSize.(uint64))

	for i, value := range rngValues {
		if _, err := wasm.Call("set_rng_value", uint64(i), uint64(value)); err != nil {
			return nil, err
		}
	}

	for i, b := range subPubKey {
		if _, err := wasm.Call("set_pk_byte", uint64(i), uint64(b)); err != nil {
			return nil, err
		}
	}

	for i, b := range message {
		if _, err := wasm.Call("set_msg_byte", uint64(i), uint64(b)); err != nil {
			return nil, err
		}
	}

	ctSize, err := wasm.Call("encrypt", uint64(len(message)))
	ciphertext := make([]byte, 0, ctSize.(uint64))
	for i := uint64(0); i < ctSize.(uint64); i++ {
		b, err := wasm.Call("get_ct_byte", i)
		if err != nil {
			return nil, err
		}
		ciphertext = append(ciphertext, byte(b.(uint64)))
	}

	return ciphertext, nil
}
