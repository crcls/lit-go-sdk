package client

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/crcls/lit-go-sdk/config"
)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var httpClient HttpClient

func init() {
	httpClient = &http.Client{}
}

func (c *Client) Connect() error {
	nodes := config.NETWORKS[c.Config.Network]
	ch := make(chan HnskMsg, len(nodes))

	for _, url := range nodes {
		go c.Handshake(url, ch)
	}

	var count uint8
	for msg := range ch {
		if msg.Connected {
			c.ConnectedNodes = append(c.ConnectedNodes, msg.Url)
			keys := *msg.Keys
			c.ServerKeysForNode[msg.Url] = keys

			if count >= c.Config.MinimumNodeCount {
				var err error
				c.ServerPubKey, err = c.MostCommonKey("ServerPubKey")
				if err != nil {
					fmt.Printf("%v\n", err)
					return err
				}
				c.SubnetPubKey, err = c.MostCommonKey("SubnetPubKey")
				if err != nil {
					fmt.Printf("%v\n", err)
					return err
				}
				c.NetworkPubKey, err = c.MostCommonKey("NetworkPubKey")
				if err != nil {
					fmt.Printf("%v\n", err)
					return err
				}
				c.NetworkPubKeySet, err = c.MostCommonKey("NetworkPubKeySet")
				if err != nil {
					fmt.Printf("%v\n", err)
					return err
				}
			}
		} else {
			log.Printf("Failed to connect to Lit Node: %s\n", msg.Url)
			log.Printf("\tReason: %s\n", msg.Error.Error())
		}

		count++
		if count == uint8(len(nodes)) {
			break
		}
	}

	if uint8(len(c.ConnectedNodes)) >= c.Config.MinimumNodeCount {
		c.Ready = true
		return nil
	}

	return fmt.Errorf("Failed to connect to enough nodes")
}

func (c *Client) NodeRequest(ctx context.Context, url string, body []byte) (*http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("LitClient: Failed to create the request for %s.\n", url)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("lit-js-sdk-version", c.Config.Version)

	return httpClient.Do(request)
}
