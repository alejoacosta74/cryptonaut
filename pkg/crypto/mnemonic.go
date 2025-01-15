package crypto

import (
	"github.com/tyler-smith/go-bip39"
)

// GenerateMnemonic creates a new random BIP39 mnemonic
func GenerateMnemonic() (string, error) {
	// Generate 256 bits of entropy (24 words)
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return "", err
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}

	return mnemonic, nil
}
