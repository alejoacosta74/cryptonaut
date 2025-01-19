/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/hex"
	"fmt"

	"github.com/alejoacosta74/cryptonaut/internal/config"
	"github.com/alejoacosta74/cryptonaut/pkg/bitcoin"
	"github.com/alejoacosta74/cryptonaut/pkg/cosmos"
	"github.com/alejoacosta74/cryptonaut/pkg/crypto"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate new crypto artifacts",
	Long: `	Generate new crypto artifacts. Examples:
		cryptonaut generate mnemonic
		cryptonaut generate key
		cryptonaut generate address [--chain bitcoin|ethereum] --from-key <key>
		cryptonaut generate key --algo schnorr
		cryptonaut generate key --algo ecdsa
		cryptonaut generate key --algo bls`,
}

var generateMnemonicCmd = &cobra.Command{
	Use:   "mnemonic",
	Short: "Generate a new mnemonic phrase",
	Long: `Generate a new mnemonic phrase using BIP39
	
	Usage:
		crypto mnemonic
	`,
	RunE: runGenerateMnemonicCmd,
}

var generatePrivateKeyCmd = &cobra.Command{
	Use:   "key",
	Short: "Generate a new private key",
	Long: `Generate a new private key in hex format.
	
	Usage:
		cryptonaut generate key
	`,
	RunE: runGeneratePrivateKeyCmd,
}

func init() {
	generateCmd.AddCommand(generateMnemonicCmd)
	generateCmd.AddCommand(generatePrivateKeyCmd)

	generatePrivateKeyCmd.Flags().StringP(config.FlagPrivateKeyFormat, "f", "hex", "Private key format (hex|wif)")
	generatePrivateKeyCmd.Flags().BoolP(config.FlagTestnet, "t", false, "Generate a private key for the testnet network")

	rootCmd.AddCommand(generateCmd)
}

func runGenerateMnemonicCmd(cmd *cobra.Command, args []string) error {
	mnemonic, err := crypto.GenerateMnemonic()
	if err != nil {
		return fmt.Errorf("failed to generate mnemonic: %v", err)
	}
	cmd.Println("Mnemonic:", mnemonic)
	return nil
}

func runGeneratePrivateKeyCmd(cmd *cobra.Command, args []string) error {
	chain := viper.GetString(config.FlagChain)
	format, err := cmd.Flags().GetString(config.FlagPrivateKeyFormat)
	if err != nil {
		return fmt.Errorf("failed to get private key format: %v", err)
	}
	switch chain {
	case "bitcoin":
		switch format {
		case "hex":
			privKey, err := crypto.GeneratePrivateKey()
			if err != nil {
				return fmt.Errorf("failed to generate private key: %v", err)
			}
			cmd.Println("Private key:", hex.EncodeToString(privKey.Serialize()))
		case "wif":
			testnet, err := cmd.Flags().GetBool(config.FlagTestnet)
			if err != nil {
				return fmt.Errorf("failed to get testnet flag: %v", err)
			}
			privKey, err := bitcoin.GeneratePrivateKeyWIF(testnet, true)
			if err != nil {
				return fmt.Errorf("failed to generate private key: %v", err)
			}
			cmd.Println("Private key:", privKey.String())
		default:
			return fmt.Errorf("invalid private key format: %s", format)
		}
	case "ethereum":
		// TODO: implement ethereum private key generation
		cmd.Println("Ethereum private key")
	case "cosmos":
		privKey := cosmos.GeneratePrivateKey()
		cmd.Println("Cosmos private key:", hex.EncodeToString(privKey.Bytes()))
	default:
		return fmt.Errorf("invalid chain: %s", chain)
	}
	return nil
}
