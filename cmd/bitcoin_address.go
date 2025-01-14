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
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("keys called")
		return nil
	},
}

var createAddressCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Bitcoin address from a private key",
	RunE:  runCreateAddressCmd,
}

func init() {
	bitcoinAddressCmd.AddCommand(createAddressCmd)

	createAddressCmd.Flags().String("from-key", "", "Private key to generate address from")
	viper.BindPFlag("from-key", createAddressCmd.Flags().Lookup("from-key"))

	bitcoinCmd.AddCommand(bitcoinAddressCmd)

}

func runCreateAddressCmd(cmd *cobra.Command, args []string) error {
	testnet := viper.GetBool("testnet")
	privKey := viper.GetString("from-key")
	address, err := bitcoin.GenerateAddress(privKey, testnet)
	if err != nil {
		cmd.PrintErrln(err)
		return err
	}
	cmd.Printf("Address: %s\n", address)
	return nil
}
