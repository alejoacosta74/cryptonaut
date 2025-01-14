/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/hex"
	"fmt"

	"github.com/alejoacosta74/cryptonaut/internal/bitcoin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// bitcoinCmd represents the bitcoin command
var bitcoinCmd = &cobra.Command{
	Use:   "bitcoin",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("bitcoin called")
	},
}

var createKeyWIFCmd = &cobra.Command{
	Use:   "create-key-wif",
	Short: "Create a new Bitcoin private key in WIF format",
	Long:  `Create a new Bitcoin private key in WIF format`,
	Run:   runCreateKeyWIFCmd,
}

var createKeyHexCmd = &cobra.Command{
	Use:   "create-key-hex",
	Short: "Create a new Bitcoin private key in hex format",
	Long:  `Create a new Bitcoin private key in hex format`,
	Run:   runCreateKeyHexCmd,
}

var generateAddressCmd = &cobra.Command{
	Use:   "generate-address",
	Short: "Generate a Bitcoin address from a private key",
	Long:  `Generate a Bitcoin address from a private key. The private key can be provided in WIF or HEX format.`,
	Args:  cobra.ExactArgs(1),
	Run:   runGenerateAddressCmd,
}

func init() {
	bitcoinCmd.AddCommand(createKeyWIFCmd)
	bitcoinCmd.AddCommand(createKeyHexCmd)
	bitcoinCmd.AddCommand(generateAddressCmd)
	bitcoinCmd.PersistentFlags().BoolP("testnet", "t", false, "Use testnet")
	bitcoinCmd.PersistentFlags().Bool("compressed", false, "Use compressed key")
	viper.BindPFlag("testnet", bitcoinCmd.PersistentFlags().Lookup("testnet"))
	viper.BindPFlag("compressed", bitcoinCmd.PersistentFlags().Lookup("compressed"))

	rootCmd.AddCommand(bitcoinCmd)
}

func runCreateKeyWIFCmd(cmd *cobra.Command, args []string) {
	testnet := viper.GetBool("testnet")
	compressed := viper.GetBool("compressed")
	privKey := bitcoin.GeneratePrivateKeyWIF(testnet, compressed)
	fmt.Println("Private key (WIF):", privKey.String())

}

func runCreateKeyHexCmd(cmd *cobra.Command, args []string) {
	privKey := bitcoin.GeneratePrivateKeyHex()
	fmt.Println("Private key (HEX):", hex.EncodeToString(privKey.Serialize()))
}

func runGenerateAddressCmd(cmd *cobra.Command, args []string) {
	testnet := viper.GetBool("testnet")
	privKey := args[0]
	address := bitcoin.GenerateAddress(privKey, testnet)
	fmt.Println("Address:", address)
}
