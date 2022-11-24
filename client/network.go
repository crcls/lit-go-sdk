package client

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/crcls/lit-go-sdk/config"
)

func (c *Client) Connect() (bool, error) {
	nodes := config.NETWORKS[c.Config.Network]
	ch := make(chan HnskMsg, len(nodes))

	for _, url := range nodes {
		go Handshake(url, ch, c.NodeRequest)
	}

	var count uint8
	for msg := range ch {
		if msg.Connected {
			c.ConnectedNodes[msg.Url] = msg.Connected
			keys := *msg.Keys
			c.ServerKeysForNode[msg.Url] = keys
			// fmt.Printf("Connected to Lit Node at: %s\n", msg.Url)

			if count >= c.Config.MinimumNodeCount {
				var err error
				c.ServerPubKey, err = c.MostCommonKey("ServerPubKey")
				if err != nil {
					fmt.Printf("%v\n", err)
					return false, err
				}
				c.SubnetPubKey, err = c.MostCommonKey("SubnetPubKey")
				if err != nil {
					fmt.Printf("%v\n", err)
					return false, err
				}
				c.NetworkPubKey, err = c.MostCommonKey("NetworkPubKey")
				if err != nil {
					fmt.Printf("%v\n", err)
					return false, err
				}
				c.NetworkPubKeySet, err = c.MostCommonKey("NetworkPubKeySet")
				if err != nil {
					fmt.Printf("%v\n", err)
					return false, err
				}
			}
		}

		count++
		if count == uint8(len(nodes)) {
			break
		}
	}

	if uint8(len(c.ConnectedNodes)) >= c.Config.MinimumNodeCount {
		c.Ready = true
		return true, nil
	}

	return false, fmt.Errorf("Failed to connect to enough nodes")
}

func (c *Client) NodeRequest(url string, body []byte) (*http.Response, error) {
	client := http.Client{
		Timeout: c.Config.RequestTimeout,
	}

	// fmt.Printf("Body: %s\n", string(body))

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("LitClient: Failed to create the request for %s.\n", url)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("lit-js-sdk-version", c.Config.Version)

	return client.Do(request)
}
