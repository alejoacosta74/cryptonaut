/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/hex"
	"fmt"

	"github.com/alejoacosta74/cryptonaut/internal/config"
	"github.com/alejoacosta74/cryptonaut/pkg/bitcoin"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/spf13/cobra"
)

// deriveCmd represents the derive command
var deriveCmd = &cobra.Command{
	Use:   "derive",
	Short: "Derive new crypto artifacts",
	Long: `Derive new crypto artifacts. Examples:
		cryptonaut derive address --from-key <key>
		cryptonaut derive address --from-key <key> --chain bitcoin
		cryptonaut derive address --from-key <key> --chain ethereum
		cryptonaut derive address --from-key <key> --chain bitcoin --testnet
	`,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("derive called")
	// },
}

var derivePublicKeyCmd = &cobra.Command{
	Use:   "pubkey",
	Short: "Derive a new public key",
	Long: `Derive a new public key from a private key in hex format
	
	Usage:
		cryptonaut derive pubkey --from-key <key>
	`,
	RunE: runDerivePublicKeyCmd,
}

var deriveAddressCmd = &cobra.Command{
	Use:   "address",
	Short: "Derive a new address",
	Long:  `Derive a new address from a private key`,
	RunE:  runDeriveAddressCmd,
}

func init() {
	deriveCmd.AddCommand(derivePublicKeyCmd)
	deriveCmd.AddCommand(deriveAddressCmd)

	deriveCmd.PersistentFlags().String(config.FlagPrivateKey, "", "Private key in hex format")
	deriveCmd.MarkPersistentFlagRequired(config.FlagPrivateKey)

	deriveAddressCmd.Flags().String(config.FlagChain, "", "Chain to derive the address from [bitcoin, ethereum]")
	deriveAddressCmd.MarkFlagRequired(config.FlagChain)

	deriveAddressCmd.PersistentFlags().Bool(config.FlagTestnet, false, "use Bitcoin testnet")

	rootCmd.AddCommand(deriveCmd)

}

func runDerivePublicKeyCmd(cmd *cobra.Command, args []string) error {
	privateKeyHex, err := cmd.Parent().PersistentFlags().GetString(config.FlagPrivateKey)
	if err != nil {
		return fmt.Errorf("invalid private key hex: %w", err)
	}

	privKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return fmt.Errorf("invalid private key hex: %w", err)
	}

	privKey, pubKey := btcec.PrivKeyFromBytes(privKeyBytes)

	cmd.Println("Private key:", hex.EncodeToString(privKey.Serialize()))
	cmd.Println("Public key:", hex.EncodeToString(pubKey.SerializeCompressed()))
	return nil
}

func runDeriveAddressCmd(cmd *cobra.Command, args []string) error {
	chain, err := cmd.Flags().GetString(config.FlagChain)
	if err != nil {
		return fmt.Errorf("invalid chain: %w", err)
	}

	if chain != "bitcoin" && chain != "ethereum" {
		return fmt.Errorf("invalid chain: %s", chain)
	}

	privateKeyHex, err := cmd.Parent().PersistentFlags().GetString(config.FlagPrivateKey)
	if err != nil {
		return fmt.Errorf("invalid private key hex: %w", err)
	}

	switch chain {
	case "bitcoin":
		// read the testnet flag
		testnet, err := cmd.PersistentFlags().GetBool(config.FlagTestnet)
		if err != nil {
			return fmt.Errorf("invalid network type for bitcoin: %w", err)
		}
		address, err := bitcoin.GenerateAddress(privateKeyHex, testnet)
		if err != nil {
			return fmt.Errorf("invalid address: %w", err)
		}
		cmd.Println("Address:", address)
	case "ethereum":
		// TODO: implement ethereum address derivation
		cmd.Println("Ethereum address")
	}
	return nil
}
