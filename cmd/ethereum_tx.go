package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"

	"github.com/alejoacosta74/cryptonaut/internal/config"
	"github.com/alejoacosta74/cryptonaut/pkg/ethereum"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ethereumTxCmd = &cobra.Command{
	Use:   "tx",
	Short: "Ethereum transaction commands",
	Long:  "Ethereum transaction commands",
}

var ethereumDecodeRawTxCmd = &cobra.Command{
	Use:   "decode",
	Short: "Decode a raw Ethereum transaction",
	Long:  "Decode a raw Ethereum transaction",
	Args:  cobra.ExactArgs(1),
	Run:   runDecodeEthereumRawTx,
}

var ethereumMempoolSubscribeCmd = &cobra.Command{
	Use:   "mempool",
	Short: "Subscribe to Ethereum mempool transactions",
	Long:  "Subscribe to Ethereum mempool transactions",
	RunE:  runSubscribeEthereumMempool,
}

func init() {
	ethereumTxCmd.AddCommand(ethereumDecodeRawTxCmd)
	ethereumTxCmd.AddCommand(ethereumMempoolSubscribeCmd)

	ethereumMempoolSubscribeCmd.Flags().StringP("to-address", "t", "", "Filter transactions by to address")
	viper.BindPFlag("to-address", ethereumMempoolSubscribeCmd.Flags().Lookup("to-address"))
	ethereumMempoolSubscribeCmd.Flags().StringP(config.FlagWsUrl, "w", "", "Websocket URL")
	ethereumMempoolSubscribeCmd.MarkFlagRequired(config.FlagWsUrl)
	viper.BindPFlag(config.FlagWsUrl, ethereumMempoolSubscribeCmd.Flags().Lookup(config.FlagWsUrl))

	ethereumCmd.AddCommand(ethereumTxCmd)
}

func runDecodeEthereumRawTx(cmd *cobra.Command, args []string) {
	rawTx := args[0]
	tx, err := ethereum.DecodeEthereumRawTx(rawTx)
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

func runSubscribeEthereumMempool(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithCancel(cmd.Context())
	defer cancel()

	wsURL := viper.GetString("ws-url")
	toAddress := viper.GetString("to-address")

	// Create the client
	client, err := ethereum.NewEthereumClient(wsURL)
	if err != nil {
		return err
	}
	defer client.Close()

	// Create and start the subscription
	sub, err := ethereum.NewMempoolSubscription(client, toAddress)
	if err != nil {
		return err
	}
	defer sub.Stop()

	if err := sub.Start(ctx); err != nil {
		return err
	}

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	return nil
}
