package client

import (
	"context"
	"encoding/json"
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

func (c *Client) StoreEncryptionConditionWithNode(
	url string,
	params SaveCondParams,
	ch chan SaveCondMsg,
) {
	reqBody, err := json.Marshal(params)
	if err != nil {
		ch <- SaveCondMsg{nil, err}
		close(ch)
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), c.Config.RequestTimeout)
	resp, err := c.NodeRequest(ctx, url+"/web/encryption/store", reqBody)
	if err != nil {
		ch <- SaveCondMsg{nil, err}
		close(ch)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ch <- SaveCondMsg{nil, err}
		close(ch)
		return
	}

	r := &SaveCondResponse{}
	if err := json.Unmarshal(body, r); err != nil {
		ch <- SaveCondMsg{nil, err}
		close(ch)
		return
	}

	ch <- SaveCondMsg{r, nil}
}
