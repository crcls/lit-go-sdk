package client

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"sort"

	"github.com/crcls/lit-go-sdk/auth"
	"github.com/crcls/lit-go-sdk/conditions"
	"github.com/crcls/lit-go-sdk/crypto"
	"github.com/crcls/lit-go-sdk/jsonUtils"
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

func MostCommonKey(serverKeys []ServerKeys, name string) string {
	counts := make(map[string]int)
	for _, keys := range serverKeys {
		if k, ok := keys.Key(name); ok {
			counts[k]++
		}
	}

	mode := ""
	maxCount := 0

	for key, count := range counts {
		if count > maxCount {
			maxCount = count
			mode = key
		}
	}

	return mode
}

type EncryptedKeyParams struct {
	AuthSig                 auth.AuthSig                        `json:"authSig"`
	Chain                   string                              `json:"chain"`
	AccessControlConditions []conditions.AccessControlCondition `json:"accessControlCondition"`
	EvmContractConditions   []conditions.EvmContractCondition   `json:"evmContractConditions"`
	SolRpcConditions        []conditions.SolRpcCondition        `json:"solRpcConditions"`
	ToDecrypt               string                              `json:"toDecrypt"`
}

func GetEncryptionKey[AC conditions.AuthCondition](
	c *Client,
	ctx context.Context,
	authSig auth.AuthSig,
	authConditions []AC,
	chain string,
	toDecrypt string,
) ([]byte, error) {
	if !c.Ready {
		return nil, fmt.Errorf("LitClient: not ready")
	}

	ch := make(chan DecryptResMsg)

	ctx, cancel := context.WithTimeout(ctx, c.Config.RequestTimeout)
	defer cancel()

	params := &EncryptedKeyParams{
		AuthSig:   authSig,
		Chain:     chain,
		ToDecrypt: toDecrypt,
	}

	ci := (interface{})(authConditions)

	switch ci.(type) {
	case []conditions.AccessControlCondition:
		params.AccessControlConditions = ci.([]conditions.AccessControlCondition)
	case []conditions.EvmContractCondition:
		params.EvmContractConditions = ci.([]conditions.EvmContractCondition)
		// case []conditions.SolRpcCondition:
	}

	reqBody, err := jsonUtils.JSONMarshal(params)
	if err != nil {
		ch <- DecryptResMsg{nil, err}
		return nil, err
	}

	for _, url := range c.ConnectedNodes {
		go GetEncryptionShare(ctx, url, c.Config.Version, reqBody, ch)
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

func SaveEncryptionKey[AC conditions.AuthCondition](
	c *Client,
	ctx context.Context,
	symmetricKey []byte,
	authSig auth.AuthSig,
	conditions []AC,
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

	condJson, err := jsonUtils.JSONMarshal(conditions)
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

	ctx, cancel := context.WithTimeout(ctx, c.Config.RequestTimeout)
	defer cancel()

	reqBody, err := json.Marshal(&scp)
	if err != nil {
		return "", err
	}

	if c.Config.Debug {
		log.Printf("Lit: Request Body JSON: %s\n", reqBody)
	}

	for _, url := range c.ConnectedNodes {
		go StoreEncryptionConditionWithNode(
			ctx,
			url,
			c.Config.Version,
			reqBody,
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
