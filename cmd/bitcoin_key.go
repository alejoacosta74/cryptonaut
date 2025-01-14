/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/hex"
	"fmt"

	"github.com/alejoacosta74/cryptonaut/pkg/bitcoin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// keysCmd represents the keys command
var bitcoinKeyCmd = &cobra.Command{
	Use:   "key",
	Short: "Create and manage Bitcoin private keys in various formats",
	Long: `Create and manage Bitcoin private keys in various formats.
    
Supported formats:
- WIF (Wallet Import Format)
- Hex (64-character hexadecimal)

Examples:
  # Create a new WIF private key for mainnet
  cryptonaut bitcoin key create --format wif

  # Create a new hex private key for testnet
  cryptonaut bitcoin key create --format hex --testnet`,
}

var createKeyCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Bitcoin private key",
	RunE:  runCreateKey,
}

var convertKeyCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert a Bitcoin private key from/to WIF and Hex",
	Args:  cobra.ExactArgs(1),
	RunE:  runConvertKey,
}

func init() {
	bitcoinKeyCmd.AddCommand(createKeyCmd)
	bitcoinKeyCmd.AddCommand(convertKeyCmd)

	bitcoinKeyCmd.PersistentFlags().String("format", "wif", "Output format: wif|hex")
	viper.BindPFlag("format", bitcoinKeyCmd.PersistentFlags().Lookup("format"))

	bitcoinCmd.AddCommand(bitcoinKeyCmd)
}

func runCreateKey(cmd *cobra.Command, args []string) error {
	format := viper.GetString("format")
	testnet := viper.GetBool("testnet")
	var result string
	switch format {
	case "wif":
		compressed := viper.GetBool("compressed")
		privKey, err := bitcoin.GeneratePrivateKeyWIF(testnet, compressed)
		if err != nil {
			cmd.PrintErrln(err)
			return err
		}
		result = privKey.String()
	case "hex":
		privKey, err := bitcoin.GeneratePrivateKeyHex()
		if err != nil {
			cmd.PrintErrln(err)
			return err
		}
		result = hex.EncodeToString(privKey.Serialize())
	default:
		return fmt.Errorf("invalid format: %s", format)
	}
	cmd.Printf("Private key (%s): %s\n", format, result)
	return nil
}

func runConvertKey(cmd *cobra.Command, args []string) error {
	privKey := args[0]
	result, err := bitcoin.ConvertKey(privKey)
	if err != nil {
		cmd.PrintErrln(err)
		return err
	}
	cmd.Printf("Private key: %s\n", result)
	return nil
}
