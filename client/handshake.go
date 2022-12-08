package client

import (
	"context"
	"encoding/json"
	"io/ioutil"
)

type HnskMsg struct {
	Url       string
	Connected bool
	Keys      *ServerKeys
	Error     error
}

func (c *Client) Handshake(url string, ch chan HnskMsg) {
	// TODO: make this configurable once supported by the network
	reqBody, err := json.Marshal(map[string]string{
		"clientPublicKey": "test",
	})
	if err != nil {
		ch <- HnskMsg{url, false, nil, err}
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), c.Config.RequestTimeout)
	resp, err := c.NodeRequest(ctx, url+"/web/handshake", reqBody)
	if err != nil {
		ch <- HnskMsg{url, false, nil, err}
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ch <- HnskMsg{url, false, nil, err}
		return
	}

	keys := ServerKeys{}
	if err := json.Unmarshal(body, &keys); err != nil {
		ch <- HnskMsg{url, false, nil, err}
		return
	}

	ch <- HnskMsg{url, true, &keys, nil}
}
