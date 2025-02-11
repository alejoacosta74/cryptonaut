/*
Package secp256k1 implements ECDSA signatures using secp256k1 curve implementation
*/
package secp256k1

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
)

func GeneratePrivateKey() (*btcec.PrivateKey, error) {
	privKey, err := btcec.NewPrivateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to create private key: %v", err)
	}
	return privKey, nil
}

func GeneratePrivateKeyHex() (string, error) {
	privKey, err := GeneratePrivateKey()
	if err != nil {
		return "", err
	}
	privKeyHex := hex.EncodeToString(privKey.Serialize())
	return privKeyHex, nil
}

func ParsePrivateKeyFromString(privKeyStr string) (*btcec.PrivateKey, error) {
	privKeyBytes, err := hex.DecodeString(privKeyStr)
	if err != nil {
		return nil, err
	}
	privKey, _ := btcec.PrivKeyFromBytes(privKeyBytes)
	return privKey, nil
}

// DerivePublicKey derives the public key from a private key
func DerivePublicKey(privateKey *btcec.PrivateKey) *btcec.PublicKey {
	if privateKey == nil {
		return nil
	}
	return privateKey.PubKey()
}

func ParsePublicKeyFromString(pubKeyStr string) (*btcec.PublicKey, error) {
	pubKeyBytes, err := hex.DecodeString(pubKeyStr)
	if err != nil {
		return nil, err
	}
	pubKey, err := btcec.ParsePubKey(pubKeyBytes)
	if err != nil {
		return nil, err
	}
	return pubKey, nil
}

func SerializePublicKeyToString(pubKey *btcec.PublicKey, isCompressed bool) string {
	var pubKeyBytes []byte
	if isCompressed {
		pubKeyBytes = pubKey.SerializeCompressed()
	} else {
		pubKeyBytes = pubKey.SerializeUncompressed()
	}

	pubkeyStr := hex.EncodeToString(pubKeyBytes)

	return pubkeyStr
}

// SignMessage signs a message using the private key
// The message is first hashed using SHA256
func SignMessage(privateKey *btcec.PrivateKey, message []byte) (*ecdsa.Signature, error) {
	if privateKey == nil {
		return nil, fmt.Errorf("private key cannot be nil")
	}
	if len(message) == 0 {
		return nil, fmt.Errorf("message cannot be empty")
	}

	// Hash the message using SHA256
	messageHash := sha256.Sum256(message)

	// Sign the hash
	signature := ecdsa.Sign(privateKey, messageHash[:])

	return signature, nil
}

// VerifySignature verifies a signature against a message and public key
// The message is first hashed using SHA256
func VerifySignature(publicKey *btcec.PublicKey, message []byte, signature *ecdsa.Signature) (bool, error) {
	if publicKey == nil {
		return false, fmt.Errorf("public key cannot be nil")
	}
	if signature == nil {
		return false, fmt.Errorf("signature cannot be nil")
	}
	if len(message) == 0 {
		return false, fmt.Errorf("message cannot be empty")
	}

	// Hash the message using SHA256
	messageHash := sha256.Sum256(message)

	// Verify the signature
	return signature.Verify(messageHash[:], publicKey), nil
}
