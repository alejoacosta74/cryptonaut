/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/hex"
	"fmt"

	"github.com/alejoacosta74/cryptonaut/internal/config"
	"github.com/alejoacosta74/cryptonaut/pkg/crypto"
	"github.com/spf13/cobra"
)

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify a message using a signature",
	Long:  `Verify a message using a signature`,
	// RunE:  runVerifyCmd,
}

var verifySchnorrCmd = &cobra.Command{
	Use:   "schnorr",
	Short: "Verify a message using a Schnorr signature",
	Long:  `Verify a message using a Schnorr signature`,
	RunE:  runVerifySchnorrCmd,
}

func init() {
	verifyCmd.AddCommand(verifySchnorrCmd)

	// Add persistent flags for verifyCmd for FlagMessage
	verifyCmd.PersistentFlags().String(config.FlagMessage, "", "Message to verify")
	verifyCmd.MarkPersistentFlagRequired(config.FlagMessage)

	// Add persistent flags for verifyCmd for FlagSignature
	verifyCmd.PersistentFlags().String(config.FlagSignature, "", "Signature to verify")
	verifyCmd.MarkPersistentFlagRequired(config.FlagSignature)

	// Add persistent flags for verifyCmd for FlagPublicKey
	verifyCmd.PersistentFlags().String(config.FlagPublicKey, "", "Public key to verify")
	verifyCmd.MarkPersistentFlagRequired(config.FlagPublicKey)

	rootCmd.AddCommand(verifyCmd)

}

func runVerifySchnorrCmd(cmd *cobra.Command, args []string) error {
	//get the flags
	message, err := cmd.Parent().PersistentFlags().GetString(config.FlagMessage)
	if err != nil {
		return fmt.Errorf("failed to get message: %w", err)
	}

	signature, err := cmd.Parent().PersistentFlags().GetString(config.FlagSignature)
	if err != nil {
		return fmt.Errorf("failed to get signature: %w", err)
	}

	publicKey, err := cmd.Parent().PersistentFlags().GetString(config.FlagPublicKey)
	if err != nil {
		return fmt.Errorf("failed to get public key: %w", err)
	}

	signatureBytes, err := hex.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("failed to decode signature: %w", err)
	}

	publicKeyBytes, err := hex.DecodeString(publicKey)
	if err != nil {
		return fmt.Errorf("failed to decode public key: %w", err)
	}

	messageBytes := []byte(message)

	valid, err := crypto.VerifyMessage(publicKeyBytes, messageBytes, signatureBytes)
	if err != nil {
		return fmt.Errorf("failed to verify message: %w", err)
	}

	cmd.Println("Signature is valid:", valid)
	return nil
}
