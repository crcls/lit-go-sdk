package auth

type AuthSig struct {
	Sig           string `json:"sig", yaml:"sig"`
	DerivedVia    string `json:"derivedVia", yaml:"derivedBy"`
	SignedMessage string `json:"signedMessage", yaml:"signedMessage"`
	Address       string `json:"address", yaml:"address"`
}
