package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

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

func StoreEncryptionConditionWithNode(
	ctx context.Context,
	url,
	version string,
	params SaveCondParams,
	ch chan SaveCondMsg,
) {
	reqBody, err := json.Marshal(params)
	if err != nil {
		ch <- SaveCondMsg{nil, err}
		return
	}

	resp, err := NodeRequest(ctx, url+"/web/encryption/store", version, reqBody)
	if err != nil {
		ch <- SaveCondMsg{nil, err}
		return
	}

	if resp.StatusCode == 500 {
		ch <- SaveCondMsg{nil, fmt.Errorf("Request failed: %s", resp.Status)}
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ch <- SaveCondMsg{nil, err}
		return
	}

	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		ch <- SaveCondMsg{nil, fmt.Errorf("Request failed: %s", string(body))}
		return
	}

	r := &SaveCondResponse{}
	if err := json.Unmarshal(body, r); err != nil {
		ch <- SaveCondMsg{nil, err}
		return
	}

	ch <- SaveCondMsg{r, nil}
}
