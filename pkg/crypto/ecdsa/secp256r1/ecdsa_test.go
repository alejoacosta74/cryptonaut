package secp256r1

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"math/big"
	"strings"
	"testing"
)

func TestGenerateAndDeriveKeys(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "successful key generation",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test private key generation
			privKey, err := GeneratePrivateKey()
			if (err != nil) != tt.wantErr {
				t.Errorf("GeneratePrivateKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if privKey == nil {
					t.Error("GeneratePrivateKey() returned nil key")
					return
				}
				if privKey.Curve != elliptic.P256() {
					t.Error("GeneratePrivateKey() used wrong curve")
				}
			}

			// Test public key derivation
			pubKey, err := DerivePublicKey(privKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("DerivePublicKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if pubKey == nil {
					t.Error("DerivePublicKey() returned nil key")
					return
				}
				if pubKey.Curve != elliptic.P256() {
					t.Error("DerivePublicKey() used wrong curve")
				}
			}
		})
	}
}

func TestSignAndVerify(t *testing.T) {
	tests := []struct {
		name        string
		message     []byte
		modifyR     bool // for testing invalid signatures
		modifyS     bool
		wantErr     bool
		errContains string
	}{
		{
			name:    "simple message",
			message: []byte("Hello, World!"),
			wantErr: false,
		},
		{
			name:    "empty message",
			message: []byte(""),
			wantErr: false,
		},
		{
			name:    "long message",
			message: bytes.Repeat([]byte("long message"), 100),
			wantErr: false,
		},
		{
			name:    "binary message",
			message: []byte{0x00, 0x01, 0x02, 0x03},
			wantErr: false,
		},
		{
			name:    "modified R component",
			message: []byte("Hello, World!"),
			modifyR: true,
			wantErr: false, // should not error but verification should fail
		},
		{
			name:    "modified S component",
			message: []byte("Hello, World!"),
			modifyS: true,
			wantErr: false, // should not error but verification should fail
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Generate key pair for each test
			privKey, err := GeneratePrivateKey()
			if err != nil {
				t.Fatalf("failed to generate private key: %v", err)
			}
			pubKey, err := DerivePublicKey(privKey)
			if err != nil {
				t.Fatalf("failed to derive public key: %v", err)
			}

			// Sign message
			r, s, err := SignMessage(privKey, tt.message)
			if (err != nil) != tt.wantErr {
				t.Fatalf("SignMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			// Modify signature components if needed (for negative testing)
			if tt.modifyR {
				r.Add(r, big.NewInt(1))
			}
			if tt.modifyS {
				s.Add(s, big.NewInt(1))
			}

			// Verify signature
			valid, err := VerifySignature(pubKey, tt.message, r, s)
			if err != nil {
				t.Fatalf("VerifySignature() error = %v", err)
			}

			// Check if verification result matches expectations
			expectedValid := !tt.modifyR && !tt.modifyS
			if valid != expectedValid {
				t.Errorf("VerifySignature() = %v, want %v", valid, expectedValid)
			}
		})
	}
}

func TestVerifySignature_Invalid(t *testing.T) {
	tests := []struct {
		name        string
		setupKey    func() *ecdsa.PublicKey
		message     []byte
		r, s        *big.Int
		wantValid   bool
		wantErr     bool
		errContains string
	}{
		{
			name: "zero R value",
			setupKey: func() *ecdsa.PublicKey {
				key, _ := GeneratePrivateKey()
				pub, _ := DerivePublicKey(key)
				return pub
			},
			message:   []byte("test"),
			r:         big.NewInt(0),
			s:         big.NewInt(1),
			wantValid: false,
			wantErr:   false,
		},
		{
			name: "zero S value",
			setupKey: func() *ecdsa.PublicKey {
				key, _ := GeneratePrivateKey()
				pub, _ := DerivePublicKey(key)
				return pub
			},
			message:   []byte("test"),
			r:         big.NewInt(1),
			s:         big.NewInt(0),
			wantValid: false,
			wantErr:   false,
		},
		{
			name:        "nil public key",
			message:     []byte("test"),
			r:           big.NewInt(1),
			s:           big.NewInt(1),
			wantValid:   false,
			wantErr:     true,
			errContains: "public key cannot be nil",
		},
		{
			name:        "nil signature components",
			message:     []byte("test"),
			wantValid:   false,
			wantErr:     true,
			errContains: "signature components (r,s) cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var pubKey *ecdsa.PublicKey
			if tt.setupKey != nil {
				pubKey = tt.setupKey()
			}

			valid, err := VerifySignature(pubKey, tt.message, tt.r, tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifySignature() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if valid != tt.wantValid {
				t.Errorf("VerifySignature() valid = %v, want %v", valid, tt.wantValid)
			}
		})
	}
}

// Benchmarks
func BenchmarkGeneratePrivateKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := GeneratePrivateKey()
		if err != nil {
			b.Fatalf("GeneratePrivateKey() error = %v", err)
		}
	}
}

func BenchmarkDerivePublicKey(b *testing.B) {
	privKey, err := GeneratePrivateKey()
	if err != nil {
		b.Fatalf("failed to generate private key: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := DerivePublicKey(privKey)
		if err != nil {
			b.Fatalf("DerivePublicKey() error = %v", err)
		}
	}
}

func BenchmarkSignMessage(b *testing.B) {
	privKey, err := GeneratePrivateKey()
	if err != nil {
		b.Fatalf("failed to generate private key: %v", err)
	}
	message := []byte("Hello, World!")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := SignMessage(privKey, message)
		if err != nil {
			b.Fatalf("SignMessage() error = %v", err)
		}
	}
}

func BenchmarkVerifySignature(b *testing.B) {
	privKey, err := GeneratePrivateKey()
	if err != nil {
		b.Fatalf("failed to generate private key: %v", err)
	}
	pubKey, err := DerivePublicKey(privKey)
	if err != nil {
		b.Fatalf("failed to derive public key: %v", err)
	}
	message := []byte("Hello, World!")
	r, s, err := SignMessage(privKey, message)
	if err != nil {
		b.Fatalf("SignMessage() error = %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		valid, err := VerifySignature(pubKey, message, r, s)
		if err != nil {
			b.Fatalf("VerifySignature() error = %v", err)
		}
		if !valid {
			b.Fatal("signature verification failed")
		}
	}
}

func TestPrivateKeySerialization(t *testing.T) {
	tests := []struct {
		name        string
		setupKey    func() (*ecdsa.PrivateKey, error)
		hexKey      string // for parse-only tests
		wantErr     bool
		errContains string
	}{
		{
			name: "roundtrip serialization",
			setupKey: func() (*ecdsa.PrivateKey, error) {
				return GeneratePrivateKey()
			},
			wantErr: false,
		},
		{
			name:        "parse invalid hex string",
			hexKey:      "invalid hex",
			wantErr:     true,
			errContains: "failed to decode hex string",
		},
		{
			name:    "parse with 0x prefix",
			hexKey:  "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			wantErr: false,
		},
		{
			name:        "parse empty string",
			hexKey:      "",
			wantErr:     true,
			errContains: "failed to decode hex string",
		},
		{
			name:        "parse odd length hex",
			hexKey:      "123",
			wantErr:     true,
			errContains: "failed to decode hex string",
		},
		{
			name:        "parse too short key",
			hexKey:      "1234",
			wantErr:     true,
			errContains: "invalid key length",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var originalKey *ecdsa.PrivateKey
			var serialized string
			var err error

			// Test serialization if we have a key setup function
			if tt.setupKey != nil {
				originalKey, err = tt.setupKey()
				if err != nil {
					t.Fatalf("failed to setup test key: %v", err)
				}

				// Test SerializePrivateKey
				serialized, err = SerializePrivateKey(originalKey)
				if (err != nil) != tt.wantErr {
					t.Errorf("SerializePrivateKey() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if tt.wantErr {
					if tt.errContains != "" && err != nil {
						if !strings.Contains(err.Error(), tt.errContains) {
							t.Errorf("error message '%v' does not contain '%v'", err, tt.errContains)
						}
					}
					return
				}

				// Verify serialized string format
				if !isValidHexString(serialized) {
					t.Errorf("SerializePrivateKey() produced invalid hex string: %v", serialized)
				}
			}

			// Test ParseECDSAPrivateKeyFromHex
			hexToTest := tt.hexKey
			if tt.setupKey != nil {
				hexToTest = serialized
			}

			parsedKey, err := ParseECDSAPrivateKeyFromHex(hexToTest)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseECDSAPrivateKeyFromHex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if tt.errContains != "" && err != nil {
					if !strings.Contains(err.Error(), tt.errContains) {
						t.Errorf("error message '%v' does not contain '%v'", err, tt.errContains)
					}
				}
				return
			}

			// For roundtrip tests, verify the parsed key matches the original
			if tt.setupKey != nil {
				if !privateKeysEqual(originalKey, parsedKey) {
					t.Error("parsed key does not match original key")
				}
			}
		})
	}
}

// Helper function to check if a string is valid hex
func isValidHexString(s string) bool {
	_, err := hex.DecodeString(strings.TrimPrefix(s, "0x"))
	return err == nil
}

// Helper function to compare private keys
func privateKeysEqual(a, b *ecdsa.PrivateKey) bool {
	if a == nil || b == nil {
		return a == b
	}
	return a.D.Cmp(b.D) == 0 &&
		a.PublicKey.X.Cmp(b.PublicKey.X) == 0 &&
		a.PublicKey.Y.Cmp(b.PublicKey.Y) == 0 &&
		a.PublicKey.Curve == b.PublicKey.Curve
}

// Benchmark serialization operations
func BenchmarkSerializePrivateKey(b *testing.B) {
	privKey, err := GeneratePrivateKey()
	if err != nil {
		b.Fatalf("failed to generate private key: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := SerializePrivateKey(privKey)
		if err != nil {
			b.Fatalf("SerializePrivateKey() error = %v", err)
		}
	}
}

func BenchmarkParseECDSAPrivateKeyFromHex(b *testing.B) {
	privKey, err := GeneratePrivateKey()
	if err != nil {
		b.Fatalf("failed to generate private key: %v", err)
	}
	serialized, err := SerializePrivateKey(privKey)
	if err != nil {
		b.Fatalf("failed to serialize key: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ParseECDSAPrivateKeyFromHex(serialized)
		if err != nil {
			b.Fatalf("ParseECDSAPrivateKeyFromHex() error = %v", err)
		}
	}
}
