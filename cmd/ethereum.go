package cmd

import (
	"encoding/hex"
	"fmt"

	"github.com/alejoacosta74/cryptonaut/internal/config"
	"github.com/alejoacosta74/cryptonaut/pkg/ethereum"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ethereumCmd = &cobra.Command{
	Use:   "ethereum",
	Short: "Ethereum commands",
	Long:  "Ethereum commands",
}

var ethereumGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a Ethereum private key",
	Long:  "Generate a Ethereum private key",
	RunE:  runEthereumGenerateCmd,
}

var ethereumPubkeyCmd = &cobra.Command{
	Use:   "pubkey",
	Short: "Get the public key from a Ethereum private key",
	Long:  "Get the public key from a Ethereum private key",
	RunE:  runEthereumPubkeyCmd,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Parent().Parent().MarkPersistentFlagRequired(config.FlagPrivateKey)
	},
}

var ethereumAddressCmd = &cobra.Command{
	Use:   "address",
	Short: "Get the Ethereum address from a private key",
	Long:  "Get the Ethereum address from a private key",
	RunE:  runEthereumAddressCmd,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Parent().Parent().MarkPersistentFlagRequired(config.FlagPrivateKey)
	},
}

func init() {
	ethereumCmd.AddCommand(ethereumGenerateCmd)
	ethereumCmd.AddCommand(ethereumPubkeyCmd)
	ethereumCmd.AddCommand(ethereumAddressCmd)
	rootCmd.AddCommand(ethereumCmd)
}

func runEthereumGenerateCmd(cmd *cobra.Command, args []string) error {
	privateKey, err := ethereum.GeneratePrivateKey()
	if err != nil {
		return fmt.Errorf("failed to generate private key: %v", err)
	}
	cmd.Println("Private Key:", hex.EncodeToString(privateKey.D.Bytes()))
	return nil
}

func runEthereumPubkeyCmd(cmd *cobra.Command, args []string) error {
	privateKeyString := viper.GetString(config.FlagPrivateKey)
	privateKey, err := ethereum.ParsePrivateKeyFromString(privateKeyString)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %v", err)
	}
	pubKey, err := ethereum.DerivePublicKey(privateKey)
	if err != nil {
		return fmt.Errorf("failed to derive public key: %v", err)
	}
	cmd.Println("Public Key:", hex.EncodeToString(pubKey.X.Bytes()))
	return nil
}

func runEthereumAddressCmd(cmd *cobra.Command, args []string) error {
	privateKeyString := viper.GetString(config.FlagPrivateKey)
	privateKey, err := ethereum.ParsePrivateKeyFromString(privateKeyString)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %v", err)
	}
	address, err := ethereum.GenerateAddress(privateKey)
	if err != nil {
		return fmt.Errorf("failed to generate address: %v", err)
	}
	cmd.Println("Address:", address)
	return nil
}
