/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/alejoacosta74/cryptonaut/pkg/bitcoin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var bitcoinAddressCmd = &cobra.Command{
	Use:   "address",
	Short: "Generate Bitcoin addresses from private keys",
	Long: `Generate Bitcoin addresses from private keys.
    
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("keys called")
	},
}

var createAddressCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Bitcoin address from a private key",
	Run:   runCreateAddressCmd,
}

func init() {
	bitcoinAddressCmd.AddCommand(createAddressCmd)

	createAddressCmd.Flags().String("from-key", "", "Private key to generate address from")
	viper.BindPFlag("from-key", createAddressCmd.Flags().Lookup("from-key"))

	bitcoinCmd.AddCommand(bitcoinAddressCmd)

}

func runCreateAddressCmd(cmd *cobra.Command, args []string) {
	testnet := viper.GetBool("testnet")
	privKey := viper.GetString("from-key")
	address := bitcoin.GenerateAddress(privKey, testnet)
	cmd.Printf("Address: %s\n", address)

}
