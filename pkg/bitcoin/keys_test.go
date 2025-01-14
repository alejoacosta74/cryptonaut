package bitcoin

import (
	"encoding/hex"
	"testing"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"
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
				t.Errorf("ConvertKey(hex) failed: %v", err)
			}
			if wifStr != tc.privateKeyWIF {
				t.Errorf("ConvertKey(hex) = %v, want %v", wifStr, tc.privateKeyWIF)
			}

			// Test ConvertKey from WIF to hex
			hexStr, err := ConvertKey(tc.privateKeyWIF)
			if err != nil {
				t.Errorf("ConvertKey(WIF) failed: %v", err)
			}
			if hexStr != tc.privateKeyHex {
				t.Errorf("ConvertKey(WIF) = %v, want %v", hexStr, tc.privateKeyHex)
			}

			// Test GenerateAddress with hex input
			addr := GenerateAddress(tc.privateKeyHex, tc.isTestnet)
			if addr != tc.address {
				t.Errorf("GenerateAddress(hex) = %v, want %v", addr, tc.address)
			}

			// Test GenerateAddress with WIF input
			addr = GenerateAddress(tc.privateKeyWIF, tc.isTestnet)
			if addr != tc.address {
				t.Errorf("GenerateAddress(WIF) = %v, want %v", addr, tc.address)
			}

			// Test DerivePublicKey
			wif, _ := btcutil.DecodeWIF(tc.privateKeyWIF)
			pubKey := DerivePublicKey(wif)
			expectedPrivKey, _ := hex.DecodeString(tc.privateKeyHex)
			expectedKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), expectedPrivKey)
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
		privKey := GeneratePrivateKeyHex()
		if privKey == nil {
			t.Error("GeneratePrivateKeyHex returned nil")
		}
		// Verify the key is valid by checking if serialization works
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
				wif := GeneratePrivateKeyWIF(tc.testnet, tc.compressed)
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
