package client

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sort"

	"github.com/crcls/lit-go-sdk/auth"
	"github.com/crcls/lit-go-sdk/conditions"
	"github.com/crcls/lit-go-sdk/crypto"
)

var thresholdDecrypt func(ctx context.Context, shares []crypto.DecryptionShare, ciphertext, netPubKeySet string) ([]byte, error)
var thresholdEncrypt func(ctx context.Context, subPubKey []byte, message []byte) ([]byte, error)

func init() {
	thresholdDecrypt = crypto.ThresholdDecrypt
	thresholdEncrypt = crypto.ThresholdEncrypt
}

type ServerKeys struct {
	ServerPubKey     string `json:"serverPublicKey"`
	SubnetPubKey     string `json:"subnetPublicKey"`
	NetworkPubKey    string `json:"networkPublicKey"`
	NetworkPubKeySet string `json:"networkPublicKeySet"`
}

func (s ServerKeys) Key(name string) (string, bool) {
	switch name {
	case "ServerPubKey":
		return s.ServerPubKey, true
	case "SubnetPubKey":
		return s.SubnetPubKey, true
	case "NetworkPubKey":
		return s.NetworkPubKey, true
	case "NetworkPubKeySet":
		return s.NetworkPubKeySet, true
	default:
		return "", false
	}
}

type DecryptionShareResponse struct {
	DecryptionShare string `json:"decryptionShare"`
	ErrorCode       string `json:"errorCode"`
	Message         string `json:"message"`
	Result          string `json:"result"`
	ShareIndex      uint8  `json:"shareIndex"`
	Status          string `json:"status"`
}

type DecryptResMsg struct {
	Share *DecryptionShareResponse
	Err   error
}

func (c *Client) GetDecryptionShare(ctx context.Context, url string, params EncryptedKeyParams, ch chan DecryptResMsg) {
	reqBody, err := json.Marshal(params)
	if err != nil {
		ch <- DecryptResMsg{nil, err}
		return
	}

	ctx, cancel := context.WithTimeout(ctx, c.Config.RequestTimeout)
	defer cancel()

	resp, err := c.NodeRequest(ctx, url+"/web/encryption/retrieve", reqBody)
	if err != nil {
		ch <- DecryptResMsg{nil, err}
		return
	}

	if resp.StatusCode == 500 {
		ch <- DecryptResMsg{nil, fmt.Errorf("Request failed: %s", resp.Status)}
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ch <- DecryptResMsg{nil, err}
		return
	}

	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		ch <- DecryptResMsg{nil, fmt.Errorf("Request failed: %s", string(body))}
		return
	}

	share := &DecryptionShareResponse{}
	if err := json.Unmarshal(body, share); err != nil {
		ch <- DecryptResMsg{nil, err}
		return
	}

	ch <- DecryptResMsg{share, nil}
}

type EncryptedKeyParams struct {
	AuthSig               *auth.AuthSig                      `json:"authSig"`
	Chain                 string                             `json:"chain"`
	EvmContractConditions []*conditions.EvmContractCondition `json:"evmContractConditions"`
	ToDecrypt             string                             `json:"toDecrypt"`
}

func (c *Client) GetEncryptionKey(
	ctx context.Context,
	params EncryptedKeyParams,
) ([]byte, error) {
	if !c.Ready {
		return nil, fmt.Errorf("LitClient: not ready")
	}

	ch := make(chan DecryptResMsg)

	for _, url := range c.ConnectedNodes {
		go c.GetDecryptionShare(ctx, url, params, ch)
	}

	shares := make([]crypto.DecryptionShare, 0)
	count := 0
	for resp := range ch {
		if resp.Err != nil || resp.Share.ErrorCode != "" {
			if c.Config.Debug {
				if resp.Err != nil {
					log.Print(resp.Err)
				} else if resp.Share.Message != "" {
					log.Print(resp.Share.Message)
				}
			}
		} else if resp.Share.Status == "fulfilled" || resp.Share.Result == "success" {
			shares = append(shares, crypto.DecryptionShare{
				Index: resp.Share.ShareIndex,
				Share: resp.Share.DecryptionShare,
			})
		}
		count++

		if count >= len(c.ConnectedNodes) {
			break
		}
	}

	if len(shares) < int(c.Config.MinimumNodeCount) {
		return nil, fmt.Errorf("LitClient: failed to retrieve enough shares")
	}

	sort.SliceStable(shares, func(i, j int) bool {
		return shares[i].Index < shares[j].Index
	})

	return thresholdDecrypt(ctx, shares, params.ToDecrypt, c.NetworkPubKeySet)
}

func (c *Client) SaveEncryptionKey(
	ctx context.Context,
	symmetricKey []byte,
	authSig auth.AuthSig,
	authConditions []conditions.EvmContractCondition,
	chain string,
	permanent bool,
) (string, error) {
	subPubKey, err := hex.DecodeString(c.SubnetPubKey)
	if err != nil {
		return "", err
	}

	key, err := thresholdEncrypt(ctx, subPubKey, symmetricKey)
	if err != nil {
		return "", err
	}

	hash := sha256.New()
	hash.Write(key)
	hashStr := hex.EncodeToString(hash.Sum(nil))

	condJson, err := json.Marshal(authConditions)
	if err != nil {
		return "", err
	}

	cHash := sha256.New()
	cHash.Write(condJson)
	cHashStr := hex.EncodeToString(cHash.Sum(nil))

	ch := make(chan SaveCondMsg)

	scp := SaveCondParams{
		Key:     hashStr,
		Val:     cHashStr,
		AuthSig: authSig,
		Chain:   chain,
	}

	if permanent {
		scp.Permanent = 1
	} else {
		scp.Permanent = 0
	}

	for _, url := range c.ConnectedNodes {
		go c.StoreEncryptionConditionWithNode(
			ctx,
			url,
			scp,
			ch,
		)
	}

	count := 0
	var e error
	for msg := range ch {
		if msg.Err != nil || msg.Response == nil {
			if c.Config.Debug {
				log.Print(msg.Err)
			}

			e = msg.Err
		}
		count++

		if count >= len(c.ConnectedNodes) {
			break
		}
	}

	if e != nil {
		return "", e
	}

	return hex.EncodeToString(key), nil
}

func (c *Client) MostCommonKey(name string) string {
	keyList := make(map[string]int)
	for _, keys := range c.ServerKeysForNode {
		k, ok := keys.Key(name)
		if !ok {
			if c.Config.Debug {
				log.Printf("MostCommonKey: Key not found: %s\n", name)
			}

			continue
		}

		if _, ok := keyList[k]; ok {
			keyList[k] += 1
		} else {
			keyList[k] = 1
		}
	}

	if len(keyList) == 0 {
		return ""
	}

	keys := make([]string, 0, len(keyList))
	for key := range keyList {
		keys = append(keys, key)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return keyList[keys[i]] > keyList[keys[j]]
	})

	return keys[0]
}
