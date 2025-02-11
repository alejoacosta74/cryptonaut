package ethereum

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestDecodeEthereumRawTx(t *testing.T) {
	// Create a private key for signing
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed to generate private key: %v", err)
	}

	tests := []struct {
		name    string
		rawTx   string
		wantErr bool
		setup   func() (*types.Transaction, string)
	}{
		{
			name:    "valid legacy transaction",
			wantErr: false,
			setup: func() (*types.Transaction, string) {
				tx := types.NewTransaction(
					1, // nonce
					common.HexToAddress("0x1234567890123456789012345678901234567890"), // to
					big.NewInt(1000000000000000000),                                   // value (1 ETH)
					21000,                                                             // gas limit
					big.NewInt(1000000000),                                            // gas price
					[]byte{},                                                          // data
				)

				signer := types.NewEIP155Signer(big.NewInt(1)) // mainnet
				signedTx, err := types.SignTx(tx, signer, privateKey)
				if err != nil {
					t.Fatalf("failed to sign tx: %v", err)
				}

				buf := new(bytes.Buffer)
				if err := signedTx.EncodeRLP(buf); err != nil {
					t.Fatalf("failed to encode tx: %v", err)
				}

				return signedTx, hex.EncodeToString(buf.Bytes())
			},
		},
		{
			name:    "invalid hex string",
			rawTx:   "not a hex string",
			wantErr: true,
			setup: func() (*types.Transaction, string) {
				return nil, ""
			},
		},
		{
			name:    "invalid RLP encoding",
			rawTx:   "0123456789",
			wantErr: true,
			setup: func() (*types.Transaction, string) {
				return nil, ""
			},
		},
		{
			name:    "empty transaction",
			rawTx:   "",
			wantErr: true,
			setup: func() (*types.Transaction, string) {
				return nil, ""
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var expectedTx *types.Transaction
			var rawTx string

			if tt.setup != nil {
				expectedTx, rawTx = tt.setup()
			} else {
				rawTx = tt.rawTx
			}

			gotTx, err := DecodeEthereumRawTx(rawTx)

			// Check error cases
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeEthereumRawTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If we expect an error, we don't need to check the transaction
			if tt.wantErr {
				return
			}

			// Compare the decoded transaction with the expected one
			if expectedTx != nil {
				// Compare transaction fields
				if gotTx.Nonce() != expectedTx.Nonce() {
					t.Errorf("Nonce mismatch: got %v, want %v", gotTx.Nonce(), expectedTx.Nonce())
				}
				if gotTx.To().Hex() != expectedTx.To().Hex() {
					t.Errorf("To address mismatch: got %v, want %v", gotTx.To().Hex(), expectedTx.To().Hex())
				}
				if gotTx.Value().Cmp(expectedTx.Value()) != 0 {
					t.Errorf("Value mismatch: got %v, want %v", gotTx.Value(), expectedTx.Value())
				}
				if gotTx.Gas() != expectedTx.Gas() {
					t.Errorf("Gas mismatch: got %v, want %v", gotTx.Gas(), expectedTx.Gas())
				}
				if gotTx.GasPrice().Cmp(expectedTx.GasPrice()) != 0 {
					t.Errorf("GasPrice mismatch: got %v, want %v", gotTx.GasPrice(), expectedTx.GasPrice())
				}
			}
		})
	}
}
