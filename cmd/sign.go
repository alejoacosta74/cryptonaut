/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/hex"
	"fmt"

	"github.com/alejoacosta74/cryptonaut/internal/config"
	"github.com/alejoacosta74/cryptonaut/pkg/crypto"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/spf13/cobra"
)

// signCmd represents the sign command
var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "Sign a message using a private key",
	Long:  `Sign a message using a private key`,
}

var signSchnorrCmd = &cobra.Command{
	Use:   "schnorr",
	Short: "Sign a message using Schnorr",
	Long:  `Sign a message using Schnorr signatures (BIP340).`,
	RunE:  runSignSchnorrCmd,
}

func init() {
	signCmd.AddCommand(signSchnorrCmd)

	// Add persistent flags for signCmd for FlagMessage
	signCmd.PersistentFlags().String(config.FlagMessage, "", "Message to sign")
	signCmd.MarkPersistentFlagRequired(config.FlagMessage)

	// Add persistent flags for signCmd for FlagPrivateKey
	signCmd.PersistentFlags().String(config.FlagPrivateKey, "", "Private key to sign the message")
	signCmd.MarkPersistentFlagRequired(config.FlagPrivateKey)

	rootCmd.AddCommand(signCmd)
}

func runSignSchnorrCmd(cmd *cobra.Command, args []string) error {
	message, err := cmd.Parent().PersistentFlags().GetString(config.FlagMessage)
	if err != nil {
		return fmt.Errorf("failed to get message: %w", err)
	}

	privateKeyHex, err := cmd.Parent().PersistentFlags().GetString(config.FlagPrivateKey)
	if err != nil {
		return fmt.Errorf("failed to get private key: %w", err)
	}

	messageBytes := []byte(message)
	privKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return fmt.Errorf("failed to decode private key: %w", err)
	}

	privKey, _ := btcec.PrivKeyFromBytes(privKeyBytes)

	signature, err := crypto.SignMessage(privKey, messageBytes)
	if err != nil {
		return fmt.Errorf("failed to sign message: %w", err)
	}

	cmd.Println("Signature:", hex.EncodeToString(signature.Serialize()))
	return nil
}
