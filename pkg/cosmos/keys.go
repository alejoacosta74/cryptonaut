package cosmos

import (
	"encoding/hex"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/types"
)

var config *types.Config

const (
	Bech32PrefixAccAddr = "cosmos"
	Bech32PrefixAccPub  = "cosmospub"
)

func init() {
	// Set the Bech32 prefix for Cosmos Hub
	config = types.GetConfig()
	config.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	config.Seal()
}

func GeneratePrivateKey() *ed25519.PrivKey {
	privKey := ed25519.GenPrivKey()
	return privKey
}

func GeneratePublicKey(privKey *ed25519.PrivKey) cryptotypes.PubKey {
	return privKey.PubKey()
}

func GenerateBech32Address(pubKey cryptotypes.PubKey) string {
	address := types.AccAddress(pubKey.Address())
	return address.String()
}

func GenerateBech32AddressFromPrivateKeyHex(privKeyHex string) string {
	privKeyBytes, err := hex.DecodeString(privKeyHex)
	if err != nil {
		return ""
	}
	privKey := ed25519.PrivKey{Key: privKeyBytes}
	return GenerateBech32Address(GeneratePublicKey(&privKey))
}
