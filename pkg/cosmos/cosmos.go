package cosmos

import (
	"encoding/hex"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/types"
)

// AddressConfig holds the configuration for address generation
type AddressConfig struct {
	AccountAddressPrefix string
	AccountPubKeyPrefix  string
}

// SetupPrefixes configures the global SDK configuration with the given prefix
func SetupPrefixes(config AddressConfig) {
	sdkConfig := types.GetConfig()
	sdkConfig.SetBech32PrefixForAccount(config.AccountAddressPrefix, config.AccountPubKeyPrefix)
	sdkConfig.Seal()
}

func GeneratePrivateKey() *ed25519.PrivKey {
	privKey := ed25519.GenPrivKey()
	return privKey
}

func GeneratePublicKey(privKey *ed25519.PrivKey) cryptotypes.PubKey {
	return privKey.PubKey()
}

func GeneratePublicKeyFromPrivateKeyHex(privKeyHex string) cryptotypes.PubKey {
	privKeyBytes, err := hex.DecodeString(privKeyHex)
	if err != nil {
		return nil
	}
	privKey := ed25519.PrivKey{Key: privKeyBytes}
	return GeneratePublicKey(&privKey)
}

func GenerateBech32Address(pubKey cryptotypes.PubKey, config AddressConfig) string {
	SetupPrefixes(config)
	address := types.AccAddress(pubKey.Address())
	return address.String()
}

func GenerateBech32AddressFromPrivateKeyHex(privKeyHex string, config AddressConfig) string {
	privKeyBytes, err := hex.DecodeString(privKeyHex)
	if err != nil {
		return ""
	}
	privKey := ed25519.PrivKey{Key: privKeyBytes}
	return GenerateBech32Address(GeneratePublicKey(&privKey), config)
}
