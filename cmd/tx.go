/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
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

	txInfo := txInfo{
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

type txInfo struct {
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
