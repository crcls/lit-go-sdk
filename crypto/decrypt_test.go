package crypto

import (
	"bytes"
	"testing"
)

func TestPKCS7UnPadding(t *testing.T) {
	plaintext := []byte("testtesttest")
	padding := bytes.Repeat([]byte{byte(4)}, 4)
	padded := append(plaintext, padding...)

	unpadded := PKCS7UnPadding(padded)

	if string(unpadded) != "testtesttest" {
		t.Errorf("Unexpested result: %s", string(unpadded))
	}
}

func TestAesEncryptDecrypt(t *testing.T) {
	key := Prng(16)
	ciphertext := AesEncrypt(key, []byte("secretsecret"))
	plaintext := AesDecrypt(key, ciphertext)

	if string(plaintext) != "secretsecret" {
		t.Errorf("Unexpected plaintext: %s", plaintext)
	}
}
