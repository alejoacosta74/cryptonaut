package secp256r1

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
)

// GeneratePrivateKey generates a new ECDSA private key using the P-256 (secp256r1) curve.
// P-256 is a NIST-standardized elliptic curve widely used in TLS/SSL and many government standards.
// It is a prime field curve with equation y² = x³ - 3x + b, where b is a specific constant.
// The curve provides 128-bit security level and is considered cryptographically secure for
// digital signatures and key exchange protocols.
func GeneratePrivateKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}

func DerivePublicKey(privateKey *ecdsa.PrivateKey) (*ecdsa.PublicKey, error) {
	return privateKey.Public().(*ecdsa.PublicKey), nil
}

func SignMessage(privateKey *ecdsa.PrivateKey, message []byte) (r, s *big.Int, err error) {
	if privateKey == nil {
		return nil, nil, fmt.Errorf("private key cannot be nil")
	}
	digest := sha256.Sum256(message)
	r, s, err = ecdsa.Sign(rand.Reader, privateKey, digest[:])
	if err != nil {
		return nil, nil, fmt.Errorf("failed to sign message: %w", err)
	}
	return r, s, nil
}

func VerifySignature(publicKey *ecdsa.PublicKey, message []byte, r, s *big.Int) (bool, error) {
	if publicKey == nil {
		return false, fmt.Errorf("public key cannot be nil")
	}
	if r == nil || s == nil {
		return false, fmt.Errorf("signature components (r,s) cannot be nil")
	}
	digest := sha256.Sum256(message)
	return ecdsa.Verify(publicKey, digest[:], r, s), nil
}

// ParseECDSAPrivateKeyFromHex converts a hex string into an *ecdsa.PrivateKey.
// It assumes the key is for the P-256 (secp256r1) curve.
func ParseECDSAPrivateKeyFromHex(hexKey string) (*ecdsa.PrivateKey, error) {
	if hexKey == "" {
		return nil, fmt.Errorf("failed to decode hex string: empty string")
	}

	// Remove optional "0x" prefix if present.
	hexKey = strings.TrimPrefix(hexKey, "0x")

	// Decode the hex string into bytes.
	keyBytes, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex string: %w", err)
	}

	if len(keyBytes) < 32 {
		return nil, fmt.Errorf("invalid key length: expected at least 32 bytes, got %d", len(keyBytes))
	}

	// Try to parse as DER encoded private key
	privKey, err := x509.ParseECPrivateKey(keyBytes)
	if err == nil {
		return privKey, nil
	}

	// If DER parsing fails, try as raw private key
	d := new(big.Int).SetBytes(keyBytes)

	// Create a new private key structure for the P-256 curve.
	curve := elliptic.P256()
	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = curve
	priv.D = d

	// Compute the public key: (X, Y) = d * G, where G is the base point.
	priv.PublicKey.X, priv.PublicKey.Y = curve.ScalarBaseMult(d.Bytes())

	// check that (X, Y) is on the curve.
	if !curve.IsOnCurve(priv.PublicKey.X, priv.PublicKey.Y) {
		return nil, fmt.Errorf("the computed public key is not on the curve")
	}

	return priv, nil
}

func ParseECDSAPublicKeyFromHex(hexKey string) (*ecdsa.PublicKey, error) {
	keyBytes, err := validateHexKeyString(hexKey)
	if err != nil {
		return nil, fmt.Errorf("failed to validate hex key string: %w", err)
	}

	// Try to parse as DER encoded public key first
	genericPubKey, err := x509.ParsePKIXPublicKey(keyBytes)
	if err == nil {
		if pubKey, ok := genericPubKey.(*ecdsa.PublicKey); ok {
			return pubKey, nil
		}
		return nil, fmt.Errorf("parsed key is not an ECDSA public key")
	}

	// If DER parsing fails, try as raw public key
	curve := elliptic.P256()
	x, y := elliptic.Unmarshal(curve, keyBytes)
	if x == nil {
		return nil, fmt.Errorf("failed to parse public key coordinates")
	}

	// Verify the point is on the curve
	if !curve.IsOnCurve(x, y) {
		return nil, fmt.Errorf("public key point is not on the curve")
	}

	return &ecdsa.PublicKey{
		Curve: curve,
		X:     x,
		Y:     y,
	}, nil
}

func SerializePrivateKey(privKey *ecdsa.PrivateKey) (string, error) {
	// Serialize (marshal) the private key to ASN.1 DER encoded form.
	derBytes, err := x509.MarshalECPrivateKey(privKey)
	if err != nil {
		return "", fmt.Errorf("failed to marshal ECDSA private key: %w", err)
	}

	return hex.EncodeToString(derBytes), nil
}

func SerializePublicKey(pubKey *ecdsa.PublicKey) string {
	pubKeyBytes := elliptic.Marshal(pubKey.Curve, pubKey.X, pubKey.Y)
	return hex.EncodeToString(pubKeyBytes)
}

func validateHexKeyString(hexKey string) ([]byte, error) {
	if hexKey == "" {
		return nil, fmt.Errorf("failed to decode hex string: empty string")
	}

	// Remove optional "0x" prefix if present
	hexKey = strings.TrimPrefix(hexKey, "0x")

	// Decode the hex string into bytes
	keyBytes, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex string: %w", err)
	}

	// Check minimum length (33 bytes for compressed, 65 for uncompressed)
	if len(keyBytes) < 33 {
		return nil, fmt.Errorf("invalid key length: must be at least 33 bytes but got %d", len(keyBytes))
	}
	return keyBytes, nil
}
