package cmd

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/alejoacosta74/cryptonaut/pkg/bitcoin"
	"github.com/spf13/cobra"
)

var bitcoinTxCmd = &cobra.Command{
	Use:   "tx",
	Short: "Bitcoin transaction commands",
	Long:  "Bitcoin transaction commands",
}

var bitcoinTxDecodeCmd = &cobra.Command{
	Use:   "decode",
	Short: "Decode a Bitcoin transaction",
	Long:  "Decode a Bitcoin transaction",
	RunE:  runBitcoinTxDecodeCmd,
}

func init() {
	bitcoinTxCmd.AddCommand(bitcoinTxDecodeCmd)

	// Add the bitcoinTxCmd to the bitcoinCmd root command
	bitcoinCmd.AddCommand(bitcoinTxCmd)
}

func runBitcoinTxDecodeCmd(cmd *cobra.Command, args []string) error {
	rawTx := args[0]
	tx, err := bitcoin.DecodeBitcoinRawTx(rawTx)
	if err != nil {
		return fmt.Errorf("Error decoding Bitcoin raw transaction: %v", err)
	}

	// Convert to our display format
	txInfo := bitcoinTxInfo{
		Hash:     tx.TxHash().String(),
		Version:  tx.Version,
		Locktime: tx.LockTime,
		Size:     tx.SerializeSize(),
		Inputs:   make([]bitcoinInputInfo, len(tx.TxIn)),
		Outputs:  make([]bitcoinOutputInfo, len(tx.TxOut)),
	}

	// Convert inputs
	for i, input := range tx.TxIn {
		txInfo.Inputs[i] = bitcoinInputInfo{
			TxID:      input.PreviousOutPoint.Hash.String(),
			Vout:      input.PreviousOutPoint.Index,
			ScriptSig: hex.EncodeToString(input.SignatureScript),
			Sequence:  input.Sequence,
		}
		// Add witness data if present
		if len(input.Witness) > 0 {
			witnessData := make([]string, len(input.Witness))
			for j, w := range input.Witness {
				witnessData[j] = hex.EncodeToString(w)
			}
			witnessJSON, _ := json.Marshal(witnessData)
			txInfo.Inputs[i].WitnessTypes = string(witnessJSON)
		}
	}

	// Convert outputs
	for i, output := range tx.TxOut {
		txInfo.Outputs[i] = bitcoinOutputInfo{
			Value:        output.Value,
			ScriptPubKey: hex.EncodeToString(output.PkScript),
		}
	}

	jsonData, err := json.MarshalIndent(txInfo, "", "    ")
	if err != nil {
		return fmt.Errorf("Error marshaling JSON: %v", err)
	}

	fmt.Println(string(jsonData))
	return nil
}

type bitcoinTxInfo struct {
	Hash     string              `json:"hash"`
	Version  int32               `json:"version"`
	Locktime uint32              `json:"locktime"`
	Size     int                 `json:"size"`
	Inputs   []bitcoinInputInfo  `json:"inputs"`
	Outputs  []bitcoinOutputInfo `json:"outputs"`
}

type bitcoinInputInfo struct {
	TxID         string `json:"txid"`
	Vout         uint32 `json:"vout"`
	ScriptSig    string `json:"scriptSig"`
	Sequence     uint32 `json:"sequence"`
	WitnessTypes string `json:"witness,omitempty"`
}

type bitcoinOutputInfo struct {
	Value        int64  `json:"value"`
	ScriptPubKey string `json:"scriptPubKey"`
}
