package client

import (
	"testing"
)

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
