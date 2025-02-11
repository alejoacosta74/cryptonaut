package schnorr

import (
	"crypto/sha256"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
)

// SignMessage signs an arbitrary message using Schnorr signature.
func SignMessage(privateKey *btcec.PrivateKey, message []byte) (*schnorr.Signature, error) {
	// Hash the message first to get a 32-byte input to comply with BIP340 signature format
	digest := sha256.Sum256(message)
	// Compute the Schnorr signature for the message
	signature, err := schnorr.Sign(privateKey, digest[:])
	if err != nil {
		return nil, fmt.Errorf("failed to sign message: %w", err)
	}
	return signature, nil
}

// VerifyMessage verifies a Schnorr signature for a given message and public key.
func VerifyMessage(publicKey []byte, message []byte, signature []byte) (bool, error) {
	sig, err := schnorr.ParseSignature(signature)
	if err != nil {
		return false, fmt.Errorf("failed to parse signature: %w", err)
	}

	pubKey, err := btcec.ParsePubKey(publicKey)
	if err != nil {
		return false, fmt.Errorf("failed to parse public key: %w", err)
	}
	// Hash the message first to get a 32-byte input to comply with BIP340 signature format
	digest := sha256.Sum256(message)
	return sig.Verify(digest[:], pubKey), nil
}
