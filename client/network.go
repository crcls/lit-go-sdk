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

func (c *Client) Connect(ctx context.Context) error {
	nodes := config.NETWORKS[c.Config.Network]
	ch := make(chan HnskMsg, len(nodes))

	ctx, cancel := context.WithTimeout(ctx, c.Config.RequestTimeout)
	defer cancel()

	for _, url := range nodes {
		go Handshake(ctx, url, c.Config.Version, ch)
	}

	var count uint8
	for msg := range ch {
		if msg.Connected {
			c.ConnectedNodes = append(c.ConnectedNodes, msg.Url)
			keys := *msg.Keys
			c.ServerKeys = append(c.ServerKeys, keys)

			if count >= c.Config.MinimumNodeCount {
				c.ServerPubKey = MostCommonKey(c.ServerKeys, "ServerPubKey")
				c.SubnetPubKey = MostCommonKey(c.ServerKeys, "SubnetPubKey")
				c.NetworkPubKey = MostCommonKey(c.ServerKeys, "NetworkPubKey")
				c.NetworkPubKeySet = MostCommonKey(c.ServerKeys, "NetworkPubKeySet")
			}
		} else if c.Config.Debug {
			log.Printf("Failed to connect to Lit Node: %s\n", msg.Url)
			log.Printf("\tReason: %s\n", msg.Error.Error())
		}

		count++
		if count == uint8(len(nodes)) {
			break
		}
	}

	if uint8(len(c.ConnectedNodes)) >= c.Config.MinimumNodeCount {
		if c.Config.Debug {
			log.Println("Connected to the Lit Network.")
		}

		c.Ready = true
		return nil
	}

	return fmt.Errorf("Failed to connect to enough nodes")
}

func NodeRequest(ctx context.Context, url, version string, body []byte) (*http.Response, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("lit-js-sdk-version", version)

	return httpClient.Do(request)
}
