package schnorr

import (
	"encoding/hex"
	"testing"

	"github.com/alejoacosta74/cryptonaut/pkg/crypto/ecdsa/secp256k1"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSchnorrSignAndVerify(t *testing.T) {
	tests := []struct {
		name          string
		message       string
		tamperMessage string // for negative testing
		wantErr       bool
		wantValid     bool
	}{
		{
			name:      "valid simple message",
			message:   "Hello, Cryptonaut!",
			wantErr:   false,
			wantValid: true,
		},
		{
			name:          "invalid signature - tampered message",
			message:       "Original message",
			tamperMessage: "Tampered message",
			wantErr:       false,
			wantValid:     false,
		},
		{
			name:      "valid empty message",
			message:   "",
			wantErr:   false,
			wantValid: true,
		},
		{
			name:      "valid long message",
			message:   "This is a very long message that we'll use to test the Schnorr signature scheme with more than just a few bytes of data to ensure it works correctly with larger inputs as well.",
			wantErr:   false,
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Generate a new private key for each test
			privKey, err := secp256k1.GeneratePrivateKey()
			require.NoError(t, err)
			require.NotNil(t, privKey)

			// Get the public key
			pubKey := privKey.PubKey()
			pubKeyBytes := pubKey.SerializeCompressed()

			// Sign the message
			signature, err := SignMessage(privKey, []byte(tt.message))
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, signature)

			// Verify the signature
			messageToVerify := tt.message
			if tt.tamperMessage != "" {
				messageToVerify = tt.tamperMessage
			}

			valid, err := VerifyMessage(pubKeyBytes, []byte(messageToVerify), signature.Serialize())
			require.NoError(t, err)
			assert.Equal(t, tt.wantValid, valid)
		})
	}
}

func TestSchnorrVerify_InvalidInputs(t *testing.T) {
	tests := []struct {
		name      string
		pubKey    string
		signature string
		message   string
		wantErr   bool
	}{
		{
			name:      "invalid public key",
			pubKey:    "invalid",
			signature: "00" + hex.EncodeToString(make([]byte, 64)), // 65-byte invalid signature
			message:   "test message",
			wantErr:   true,
		},
		{
			name:      "invalid signature",
			pubKey:    hex.EncodeToString(make([]byte, btcec.PubKeyBytesLenCompressed)),
			signature: "invalid",
			message:   "test message",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pubKeyBytes, _ := hex.DecodeString(tt.pubKey)
			sigBytes, _ := hex.DecodeString(tt.signature)

			_, err := VerifyMessage(pubKeyBytes, []byte(tt.message), sigBytes)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
