package crypto

import (
	"bytes"
	"context"
	"encoding/hex"
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

func TestThresholdDecrypt(t *testing.T) {
	newWasmInstance = mockNewWasmInstance

	t.Run("Success", func(t *testing.T) {
		share := DecryptionShare{1, hex.EncodeToString([]byte("test"))}

		result, err := ThresholdDecrypt(
			context.Background(),
			[]DecryptionShare{share},
			hex.EncodeToString([]byte("ciphertext")),
			hex.EncodeToString([]byte("netPubKeySet")),
		)

		if err != nil {
			t.Errorf("Unexpected error: %s", err.Error())
		}

		// MockWasm Call returns uint64(1) when a return is expected
		if result[0] != 1 {
			t.Errorf("Unexpected result: %+v", result)
		}
	})

	t.Run("Fail if share is not hex", func(t *testing.T) {
		share := DecryptionShare{1, "test"}

		_, err := ThresholdDecrypt(
			context.Background(),
			[]DecryptionShare{share},
			hex.EncodeToString([]byte("ciphertext")),
			hex.EncodeToString([]byte("netPubKeySet")),
		)

		if err == nil {
			t.Errorf("Expected error")
		}
	})

	t.Run("Fail if ciphertext is not hex", func(t *testing.T) {
		share := DecryptionShare{1, hex.EncodeToString([]byte("test"))}

		_, err := ThresholdDecrypt(
			context.Background(),
			[]DecryptionShare{share},
			"ciphertext",
			hex.EncodeToString([]byte("netPubKeySet")),
		)

		if err == nil {
			t.Errorf("Expected error")
		}
	})

	t.Run("Fail if netPubKeySet is not hex", func(t *testing.T) {
		share := DecryptionShare{1, hex.EncodeToString([]byte("test"))}

		_, err := ThresholdDecrypt(
			context.Background(),
			[]DecryptionShare{share},
			hex.EncodeToString([]byte("ciphertext")),
			"netPubKeySet",
		)

		if err == nil {
			t.Errorf("Expected error")
		}
	})
}
