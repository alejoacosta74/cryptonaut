package ethereum

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
)

// GeneratePrivateKey generates a new ECDSA private key using the secp256k1 curve.
// secp256k1 is a Koblitz curve used by Bitcoin and Ethereum, defined by the equation y² = x³ + 7.
// Unlike other curves, secp256k1 was not generated randomly but chosen for having certain
// desirable properties. Its parameters were selected in a way that allows for especially
// efficient computation, making it ideal for blockchain applications where performance is critical.
// The curve provides 128-bit security level and is considered cryptographically secure for
// digital signatures and key exchange protocols.
func GeneratePrivateKey() (*ecdsa.PrivateKey, error) {
	return crypto.GenerateKey()
}

func DerivePublicKey(privateKey *ecdsa.PrivateKey) (*ecdsa.PublicKey, error) {
	return privateKey.Public().(*ecdsa.PublicKey), nil
}

func GenerateAddress(privateKey *ecdsa.PrivateKey) (string, error) {
	pubKey, err := DerivePublicKey(privateKey)
	if err != nil {
		return "", err
	}
	return crypto.PubkeyToAddress(*pubKey).Hex(), nil
}

// signMessage signs the given message using the provided ECDSA private key.
// It returns a 65-byte signature where:
//   - bytes [0:32] contain R,
//   - bytes [32:64] contain S, and
//   - byte [64] contains V (the recovery id plus 27).
func SignMessage(privateKey *ecdsa.PrivateKey, message []byte) ([]byte, error) {
	// In Ethereum, the signature is computed on the Keccak256 hash of the message.
	// (Often you’ll prefix the message to prevent signing arbitrary data,
	//  but here we assume the message is already properly formatted.)
	hash := crypto.Keccak256(message)

	// crypto.Sign returns the 65-byte [R || S || V] signature.
	signature, err := crypto.Sign(hash, privateKey)
	if err != nil {
		return nil, err
	}
	return signature, nil
}

// recoverPublicKey recovers the ECDSA public key from the signature and message.
// The function returns the recovered public key if the signature is valid.
func RecoverPublicKey(message, signature []byte) (*ecdsa.PublicKey, error) {
	// Compute the message hash.
	hash := crypto.Keccak256(message)

	// crypto.SigToPub recovers the public key from the signature and the hash.
	pubKey, err := crypto.SigToPub(hash, signature)
	if err != nil {
		return nil, err
	}
	return pubKey, nil
}
