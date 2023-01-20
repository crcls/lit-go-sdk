package client

import (
	"github.com/crcls/lit-go-sdk/auth"
)

type SaveCondParams struct {
	Key       string       `json:"key"`
	Val       string       `json:"val"`
	AuthSig   auth.AuthSig `json:"authSig"`
	Chain     string       `json:"chain"`
	Permanent int          `json:"permanant"` // Purposely misspelled to match API
}

type SaveCondResponse struct {
	Result string `json:"result"`
	Error  string `json:"error"`
}

type SaveCondMsg struct {
	Response *SaveCondResponse
	Err      error
}
