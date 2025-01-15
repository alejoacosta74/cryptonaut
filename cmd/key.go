/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// cryptoCmd represents the crypto command
var cryptoCmd = &cobra.Command{
	Use:   "crypto",
	Short: "Cryptographic operations",
	Long:  `Cryptographic operations.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("crypto called")
	},
}

var generatePublicKeyCmd = &cobra.Command{
	Use:   "pubkey",
	Short: "Generate a public key from a private key",
	Long: `Generate a public key from a private key.
	
	Usage:
		cryptonaut crypto pubkey --private-key "your-private-key-hex"
	`,
	RunE: runGeneratePublicKeyCmd,
}

func init() {
	cryptoCmd.AddCommand(generatePublicKeyCmd)

	generatePublicKeyCmd.Flags().String("private-key", "", "Private key in hex format")
	generatePublicKeyCmd.MarkFlagRequired("private-key")

	viper.BindPFlag("crypto.private_key", generatePublicKeyCmd.Flags().Lookup("private-key"))

	rootCmd.AddCommand(cryptoCmd)
}

func runGeneratePublicKeyCmd(cmd *cobra.Command, args []string) error {
	privateKeyHex := viper.GetString("crypto.private_key")

	privKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return fmt.Errorf("invalid private key hex: %w", err)
	}

	privKey, pubKey := btcec.PrivKeyFromBytes(privKeyBytes)

	cmd.Println("Private key:", hex.EncodeToString(privKey.Serialize()))
	cmd.Println("Public key:", hex.EncodeToString(pubKey.SerializeCompressed()))
	return nil
}
