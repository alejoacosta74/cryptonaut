package bitcoin

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

func TestDecodeBitcoinRawTx(t *testing.T) {
	// Create a private key for signing
	privateKey, err := btcec.NewPrivateKey()
	if err != nil {
		t.Fatalf("failed to generate private key: %v", err)
	}
	pubKey := privateKey.PubKey()

	tests := []struct {
		name    string
		rawTx   string
		wantErr bool
		setup   func() (*wire.MsgTx, string)
	}{
		{
			name:    "valid P2PKH transaction",
			wantErr: false,
			setup: func() (*wire.MsgTx, string) {
				// Create a new transaction
				tx := wire.NewMsgTx(wire.TxVersion)

				// Add an input
				prevHash, _ := chainhash.NewHashFromStr("1234567890123456789012345678901234567890123456789012345678901234")
				outpoint := wire.NewOutPoint(prevHash, 0)
				txIn := wire.NewTxIn(outpoint, nil, nil)
				tx.AddTxIn(txIn)

				// Create a P2PKH script
				pubKeyBytes := pubKey.SerializeCompressed()
				pubKeyHash := btcutil.Hash160(pubKeyBytes)
				pkScript, err := txscript.NewScriptBuilder().
					AddOp(txscript.OP_DUP).
					AddOp(txscript.OP_HASH160).
					AddData(pubKeyHash).
					AddOp(txscript.OP_EQUALVERIFY).
					AddOp(txscript.OP_CHECKSIG).
					Script()
				if err != nil {
					t.Fatalf("failed to create pkScript: %v", err)
				}

				// Add an output
				txOut := wire.NewTxOut(100000000, pkScript) // 1 BTC
				tx.AddTxOut(txOut)

				// Serialize the transaction
				buf := new(bytes.Buffer)
				if err := tx.Serialize(buf); err != nil {
					t.Fatalf("failed to serialize tx: %v", err)
				}

				return tx, hex.EncodeToString(buf.Bytes())
			},
		},
		{
			name:    "invalid hex string",
			rawTx:   "not a hex string",
			wantErr: true,
			setup: func() (*wire.MsgTx, string) {
				return nil, ""
			},
		},
		{
			name:    "invalid transaction format",
			rawTx:   "0123456789",
			wantErr: true,
			setup: func() (*wire.MsgTx, string) {
				return nil, ""
			},
		},
		{
			name:    "empty transaction",
			rawTx:   "",
			wantErr: true,
			setup: func() (*wire.MsgTx, string) {
				return nil, ""
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var expectedTx *wire.MsgTx
			var rawTx string

			if tt.setup != nil {
				expectedTx, rawTx = tt.setup()
			} else {
				rawTx = tt.rawTx
			}
			fmt.Println("rawTx", rawTx)
			gotTx, err := DecodeBitcoinRawTx(rawTx)

			// Check error cases
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeBitcoinRawTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// If we expect an error, we don't need to check the transaction
			if tt.wantErr {
				return
			}

			// Compare the decoded transaction with the expected one
			if expectedTx != nil {
				// Compare version
				if gotTx.Version != expectedTx.Version {
					t.Errorf("Version mismatch: got %v, want %v", gotTx.Version, expectedTx.Version)
				}

				// Compare number of inputs
				if len(gotTx.TxIn) != len(expectedTx.TxIn) {
					t.Errorf("Input count mismatch: got %v, want %v", len(gotTx.TxIn), len(expectedTx.TxIn))
				}

				// Compare number of outputs
				if len(gotTx.TxOut) != len(expectedTx.TxOut) {
					t.Errorf("Output count mismatch: got %v, want %v", len(gotTx.TxOut), len(expectedTx.TxOut))
				}

				// Compare first input
				if len(gotTx.TxIn) > 0 && len(expectedTx.TxIn) > 0 {
					if gotTx.TxIn[0].PreviousOutPoint.Hash.String() != expectedTx.TxIn[0].PreviousOutPoint.Hash.String() {
						t.Errorf("Input hash mismatch: got %v, want %v",
							gotTx.TxIn[0].PreviousOutPoint.Hash.String(),
							expectedTx.TxIn[0].PreviousOutPoint.Hash.String())
					}
				}

				// Compare first output
				if len(gotTx.TxOut) > 0 && len(expectedTx.TxOut) > 0 {
					if gotTx.TxOut[0].Value != expectedTx.TxOut[0].Value {
						t.Errorf("Output value mismatch: got %v, want %v",
							gotTx.TxOut[0].Value, expectedTx.TxOut[0].Value)
					}
				}

				// Compare locktime
				if gotTx.LockTime != expectedTx.LockTime {
					t.Errorf("LockTime mismatch: got %v, want %v", gotTx.LockTime, expectedTx.LockTime)
				}
			}
		})
	}
}
