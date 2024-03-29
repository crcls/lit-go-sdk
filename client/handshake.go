package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type HnskMsg struct {
	Url       string
	Connected bool
	Keys      *ServerKeys
	Error     error
}

func Handshake(ctx context.Context, url, version string, ch chan HnskMsg) {
	// TODO: make this configurable once supported by the network
	reqBody, err := json.Marshal(map[string]string{
		"clientPublicKey": "test",
	})
	if err != nil {
		ch <- HnskMsg{url, false, nil, err}
		return
	}

	resp, err := NodeRequest(ctx, url+"/web/handshake", version, reqBody)
	if err != nil {
		ch <- HnskMsg{url, false, nil, err}
		return
	}

	if resp.StatusCode == 500 {
		ch <- HnskMsg{url, false, nil, fmt.Errorf("Request failed: %s", resp.Status)}
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ch <- HnskMsg{url, false, nil, err}
		return
	}

	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		ch <- HnskMsg{url, false, nil, fmt.Errorf("Request failed: %s", string(body))}
		return
	}

	keys := ServerKeys{}
	if err := json.Unmarshal(body, &keys); err != nil {
		ch <- HnskMsg{url, false, nil, err}
		return
	}

	keyNames := [4]string{
		"ServerPubKey",
		"SubnetPubKey",
		"NetworkPubKey",
		"NetworkPubKeySet",
	}

	for _, keyName := range keyNames {
		key, _ := keys.Key(keyName)
		if key == "" {
			ch <- HnskMsg{url, false, nil, fmt.Errorf("Missing Key in handshake response.")}
			return
		}
	}

	ch <- HnskMsg{url, true, &keys, nil}
}
