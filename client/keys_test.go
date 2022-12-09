package client

import (
	"testing"
)

var testKeys = `{
	"serverPublicKey": "serverPubKey",
	"subnetPublicKey": "subnetPubKey",
	"networkPublicKey": "networkPubKey",
	"networkPublicKeySet": "networkPubKeySet"
}`

func TestServerKeysKeys(t *testing.T) {
	sk := ServerKeys{
		"serverPubKey",
		"subnetPubKey",
		"networkPubKey",
		"networkPubKeySet",
	}

	if spk, ok := sk.Key("ServerPubKey"); !ok || spk != "serverPubKey" {
		t.Errorf("Unexpected key for `ServerPubKey`: %s", spk)
	}

	if sbk, ok := sk.Key("SubnetPubKey"); !ok || sbk != "subnetPubKey" {
		t.Errorf("Unexpected key for `SubnetPubKey`: %s", sbk)
	}

	if npk, ok := sk.Key("NetworkPubKey"); !ok || npk != "networkPubKey" {
		t.Errorf("Unexpected key for `NetworkPubKey`: %s", npk)
	}

	if npks, ok := sk.Key("NetworkPubKeySet"); !ok || npks != "networkPubKeySet" {
		t.Errorf("Unexpected key for `NetworkPubKeySet`: %s", npks)
	}
}

func TestMostCommonKey(t *testing.T) {
	client := &Client{
		ServerKeysForNode: map[string]ServerKeys{
			"http://localhost:7470": ServerKeys{
				ServerPubKey: "common",
			},
			"http://localhost:7471": ServerKeys{
				ServerPubKey: "common",
			},
			"http://localhost:7472": ServerKeys{
				ServerPubKey: "uncommon",
			},
		},
	}

	key := client.MostCommonKey("ServerPubKey")
	if key != "common" {
		t.Errorf("Unexpected result from MostCommonKey: expected `common` got %s", key)
	}
}
