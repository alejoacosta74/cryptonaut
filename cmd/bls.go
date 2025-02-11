package cmd

import (
	"fmt"

	"github.com/alejoacosta74/cryptonaut/internal/config"
	"github.com/alejoacosta74/cryptonaut/pkg/crypto/bls"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var blsCmd = &cobra.Command{
	Use:   "bls",
	Short: "BLS operations",
}

var blsGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a BLS private key",
	RunE:  runBlsGenerateCmd,
}

var blsPubKeyCmd = &cobra.Command{
	Use:   "pubkey",
	Short: "Generate a BLS public key from a private key",
	RunE:  runBlsPubKeyCmd,
	PreRun: func(cmd *cobra.Command, args []string) {
		cmd.MarkFlagRequired(config.FlagPrivateKey)
	},
}

var blsSignCmd = &cobra.Command{
	Use:   "sign",
	Short: "Sign a message using BLS",
	Args:  cobra.ExactArgs(1), // message to sign
	RunE:  runBlsSignCmd,
	PreRun: func(cmd *cobra.Command, args []string) {
		cmd.MarkFlagRequired(config.FlagPrivateKey)
	},
}

var blsVerifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify a BLS signature",
	Args:  cobra.ExactArgs(1), // message to verify
	RunE:  runBlsVerifyCmd,
	PreRun: func(cmd *cobra.Command, args []string) {
		cmd.MarkFlagRequired(config.FlagPublicKey)
		cmd.MarkFlagRequired(config.FlagSignature)
	},
}

func init() {
	blsCmd.AddCommand(blsGenerateCmd)
	blsCmd.AddCommand(blsPubKeyCmd)
	blsCmd.AddCommand(blsSignCmd)
	blsCmd.AddCommand(blsVerifyCmd)

	rootCmd.AddCommand(blsCmd)
}

func runBlsGenerateCmd(cmd *cobra.Command, args []string) error {
	privateKey, err := bls.GeneratePrivateKey()
	if err != nil {
		return fmt.Errorf("failed to generate private key: %w", err)
	}
	cmd.Println("Private key:", privateKey.SerializeToHexStr())
	return nil
}

func runBlsPubKeyCmd(cmd *cobra.Command, args []string) error {
	privateKeyStr := viper.GetString(config.FlagPrivateKey)
	privateKey, err := bls.ParsePrivateKeyFromString(privateKeyStr)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}
	publicKey := bls.DerivePublicKey(privateKey)
	cmd.Println("Public key:", publicKey.SerializeToHexStr())
	return nil
}

func runBlsSignCmd(cmd *cobra.Command, args []string) error {
	privateKeyStr := viper.GetString(config.FlagPrivateKey)
	privateKey, err := bls.ParsePrivateKeyFromString(privateKeyStr)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}
	message := args[0]
	signature, err := bls.SignMessage(privateKey, message)
	if err != nil {
		return fmt.Errorf("failed to sign message: %w", err)
	}
	cmd.Println("Signature:", signature.SerializeToHexStr())
	return nil
}

func runBlsVerifyCmd(cmd *cobra.Command, args []string) error {
	publicKeyStr := viper.GetString(config.FlagPublicKey)
	publicKey, err := bls.ParsePublicKeyFromString(publicKeyStr)
	if err != nil {
		return fmt.Errorf("failed to parse public key: %w", err)
	}
	message := args[0]
	signatureStr := viper.GetString(config.FlagSignature)
	signature, err := bls.ParseSignatureFromString(signatureStr)
	if err != nil {
		return fmt.Errorf("failed to parse signature: %w", err)
	}
	valid, err := bls.VerifyMessage(publicKey, message, signature)
	if err != nil {
		return fmt.Errorf("failed to verify message: %w", err)
	}
	cmd.Println("Signature is valid:", valid)
	return nil
}
