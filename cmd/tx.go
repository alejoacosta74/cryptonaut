/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/alejoacosta74/cryptonaut/tx"
	"github.com/spf13/cobra"
)

// txCmd represents the tx command
var txCmd = &cobra.Command{
	Use:   "tx",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("tx called")
	},
}

var decodeRawTxCmd = &cobra.Command{
	Use:   "decode",
	Short: "Decode a raw transaction",
	Long:  `Decode a raw transaction`,
}

var decodeEthereumRawTxCmd = &cobra.Command{
	Use:   "ethereum",
	Short: "Decode an Ethereum raw transaction",
	Long:  `Decode an Ethereum raw transaction`,
	Args:  cobra.ExactArgs(1),
	Run:   decodeEthereumRawTx,
}

var decodeBitcoinRawTxCmd = &cobra.Command{
	Use:   "bitcoin",
	Short: "Decode a Bitcoin raw transaction",
	Long:  `Decode a Bitcoin raw transaction`,
	Args:  cobra.ExactArgs(1),
	Run:   decodeBitcoinRawTx,
}

func init() {
	txCmd.AddCommand(decodeRawTxCmd)
	decodeRawTxCmd.AddCommand(decodeEthereumRawTxCmd)
	decodeRawTxCmd.AddCommand(decodeBitcoinRawTxCmd)

	rootCmd.AddCommand(txCmd)
}

func decodeEthereumRawTx(cmd *cobra.Command, args []string) {
	rawTx := args[0]
	tx, err := tx.DecodeEthereumRawTx(rawTx)
	if err != nil {
		fmt.Println("Error decoding Ethereum raw transaction:", err)
		return
	}

	txInfo := ethereumTxInfo{
		Hash:     tx.Hash().String(),
		Nonce:    tx.Nonce(),
		GasPrice: tx.GasPrice().String(),
		Gas:      tx.Gas(),
		To:       tx.To().String(),
		Value:    tx.Value().String(),
		Data:     string(tx.Data()),
		ChainID:  tx.ChainId().String(),
		Type:     tx.Type(),
	}

	jsonData, err := json.MarshalIndent(txInfo, "", "    ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	fmt.Println(string(jsonData))
}

type ethereumTxInfo struct {
	Hash                 string `json:"hash"`
	Nonce                uint64 `json:"nonce"`
	GasPrice             string `json:"gasPrice"`
	Gas                  uint64 `json:"gas"`
	To                   string `json:"to"`
	Value                string `json:"value"`
	Data                 string `json:"data"`
	ChainID              string `json:"chainId"`
	Type                 uint8  `json:"type"`
	AccessList           string `json:"accessList,omitempty"`
	MaxFeePerGas         string `json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas string `json:"maxPriorityFeePerGas,omitempty"`
	Signature            string `json:"signature,omitempty"`
}

func decodeBitcoinRawTx(cmd *cobra.Command, args []string) {
	rawTx := args[0]
	tx, err := tx.DecodeBitcoinRawTx(rawTx)
	if err != nil {
		fmt.Println("Error decoding Bitcoin raw transaction:", err)
		return
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
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	fmt.Println(string(jsonData))
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
