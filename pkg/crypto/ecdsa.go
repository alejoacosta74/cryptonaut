package crypto

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
)

// GeneratePrivateKey generates a new private key using the btcsuite library.
func GeneratePrivateKey() (*btcec.PrivateKey, error) {
	// Generate a private key using btcec
	privKey, err := btcec.NewPrivateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}
	return privKey, nil
}

// GeneratePublicKey generates a public key from a private key.
func GeneratePublicKey(privateKey *btcec.PrivateKey) (*btcec.PublicKey, error) {
	return privateKey.PubKey(), nil
}
