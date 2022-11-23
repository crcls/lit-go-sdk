package auth

type AuthSig struct {
	Sig           string `json:"sig"`
	DerivedVia    string `json:"derivedVia"`
	SignedMessage string `json:"signedMessage"`
	Address       string `json:"address"`
}
