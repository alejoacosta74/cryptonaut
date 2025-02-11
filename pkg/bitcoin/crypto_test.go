package bitcoin

import (
	"encoding/hex"
	"testing"

	"github.com/alejoacosta74/cryptonaut/pkg/crypto/ecdsa/secp256k1"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/stretchr/testify/require"
)

func TestBitcoinKeys(t *testing.T) {
	testCases := []struct {
		name          string
		privateKeyHex string
		privateKeyWIF string
		address       string
		isTestnet     bool
		isCompressed  bool
	}{
		{
			name:          "mainnet compressed key",
			privateKeyHex: "9df5a907ff17ed6a4e02c00c2c119049a045f52a4e817b06b2ec54eb68f70079",
			privateKeyWIF: "L2WmFR8WMr5GSprjt7UTA7WQ23WDEZPVRimrZv1dmz7e4JzxqSNq",
			address:       "1EoxGLjv4ZADtRBjTVeXY35czVyDdp7rU4",
			isTestnet:     false,
			isCompressed:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test ConvertKey from hex to WIF
			wifStr, err := ConvertKey(tc.privateKeyHex)
			if err != nil {
				t.Fatalf("ConvertKey(hex) failed: %v", err)
			}
			if wifStr != tc.privateKeyWIF {
				t.Errorf("ConvertKey(hex) = %v, want %v", wifStr, tc.privateKeyWIF)
			}

			// Test ConvertKey from WIF to hex
			hexStr, err := ConvertKey(tc.privateKeyWIF)
			if err != nil {
				t.Fatalf("ConvertKey(WIF) failed: %v", err)
			}
			if hexStr != tc.privateKeyHex {
				t.Errorf("ConvertKey(WIF) = %v, want %v", hexStr, tc.privateKeyHex)
			}

			// Test GenerateAddress with hex input
			addr, err := GenerateAddressFromPrivateKey(tc.privateKeyHex, tc.isTestnet)
			if err != nil {
				t.Fatalf("GenerateAddress(hex) failed: %v", err)
			}
			if addr != tc.address {
				t.Errorf("GenerateAddress(hex) = %v, want %v", addr, tc.address)
			}

			// Test GenerateAddress with WIF input
			addr, err = GenerateAddressFromPrivateKey(tc.privateKeyWIF, tc.isTestnet)
			if err != nil {
				t.Fatalf("GenerateAddress(WIF) failed: %v", err)
			}
			if addr != tc.address {
				t.Errorf("GenerateAddress(WIF) = %v, want %v", addr, tc.address)
			}

			// Test DerivePublicKey
			wif, _ := btcutil.DecodeWIF(tc.privateKeyWIF)
			pubKey, err := GeneratePublicKeyFromWIF(wif)
			if err != nil {
				t.Fatalf("DerivePublicKey failed: %v", err)
			}
			expectedPrivKey, _ := hex.DecodeString(tc.privateKeyHex)
			expectedKey, _ := btcec.PrivKeyFromBytes(expectedPrivKey)
			expectedPubKey := expectedKey.PubKey().SerializeCompressed()
			if string(pubKey) != string(expectedPubKey) {
				t.Errorf("DerivePublicKey = %x, want %x", pubKey, expectedPubKey)
			}
		})
	}
}

// TestGeneratePrivateKey tests the key generation functions
func TestGeneratePrivateKey(t *testing.T) {
	t.Run("generate hex private key", func(t *testing.T) {
		privKey, err := secp256k1.GeneratePrivateKey()
		if err != nil {
			t.Fatalf("GeneratePrivateKeyHex failed: %v", err)
		}
		if privKey == nil {
			t.Error("GeneratePrivateKeyHex returned nil")
		}
		if len(privKey.Serialize()) != 32 {
			t.Error("Generated invalid private key")
		}
	})

	t.Run("generate WIF private key", func(t *testing.T) {
		testCases := []struct {
			name       string
			testnet    bool
			compressed bool
		}{
			{"mainnet compressed", false, true},
			{"mainnet uncompressed", false, false},
			{"testnet compressed", true, true},
			{"testnet uncompressed", true, false},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				privKey, err := secp256k1.GeneratePrivateKey()
				require.NoError(t, err)
				wif, err := ConvertPrivateKeyToWIF(privKey, tc.testnet, tc.compressed)
				if err != nil {
					t.Fatalf("GeneratePrivateKeyWIF failed: %v", err)
				}
				if wif == nil {
					t.Error("GeneratePrivateKeyWIF returned nil")
				}
				if wif.IsForNet(getChainParams(!tc.testnet)) {
					t.Error("Generated WIF for wrong network")
				}
				if wif.CompressPubKey != tc.compressed {
					t.Error("Generated WIF with wrong compression setting")
				}
			})
		}
	})
}
