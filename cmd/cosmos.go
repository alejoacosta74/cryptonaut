package cmd

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/alejoacosta74/cryptonaut/internal/config"
	"github.com/alejoacosta74/cryptonaut/pkg/cosmos"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cosmosCmd = &cobra.Command{
	Use:   "cosmos",
	Short: "Cosmos key generation and manipulation",
	Long:  `Cosmos key generation and manipulation.`,
}

var cosmosGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a new Cosmos private key",
	Long:  `Generate a new Cosmos private key`,
	RunE:  runCosmosGenerateCmd,
}

var cosmosPubkeyCmd = &cobra.Command{
	Use:   "pubkey",
	Short: "Derive a Cosmos public key from a private key",
	Long:  `Derive a Cosmos public key from a private key`,
	RunE:  runCosmosPubkeyCmd,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Parent().Parent().MarkPersistentFlagRequired(config.FlagPrivateKey)
	},
}

var cosmosAddressCmd = &cobra.Command{
	Use:   "address",
	Short: "Gets a Cosmos address from a public key",
	Long:  `Gets a Cosmos address from a public key`,
	RunE:  runCosmosAddressCmd,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Parent().Parent().MarkPersistentFlagRequired(config.FlagPrivateKey)
	},
}

func init() {
	cosmosCmd.AddCommand(cosmosGenerateCmd)
	cosmosCmd.AddCommand(cosmosPubkeyCmd)
	cosmosCmd.AddCommand(cosmosAddressCmd)

	cosmosAddressCmd.Flags().String(config.FlagCosmosAddressPrefix, "cosmos", "Cosmos address prefix")
	viper.BindPFlag(config.FlagCosmosAddressPrefix, cosmosAddressCmd.Flags().Lookup(config.FlagCosmosAddressPrefix))

	rootCmd.AddCommand(cosmosCmd)
}

func runCosmosGenerateCmd(cmd *cobra.Command, args []string) error {
	privKey := cosmos.GeneratePrivateKey()
	cmd.Println("Cosmos private key:", hex.EncodeToString(privKey.Bytes()))
	return nil
}

func runCosmosPubkeyCmd(cmd *cobra.Command, args []string) error {
	privKeyHex := viper.GetString(config.FlagPrivateKey)
	if privKeyHex == "" {
		return fmt.Errorf("private key is required")
	}
	pubKey := cosmos.GeneratePublicKeyFromPrivateKeyHex(privKeyHex)
	pubkeyString := strings.TrimPrefix(pubKey.String(), "PubKeyEd25519{")
	pubkeyString = strings.TrimSuffix(pubkeyString, "}")
	cmd.Println("Cosmos public key:", pubkeyString)
	return nil
}

func runCosmosAddressCmd(cmd *cobra.Command, args []string) error {
	privKeyHex := viper.GetString(config.FlagPrivateKey)
	cosmosAddrPrefix := viper.GetString(config.FlagCosmosAddressPrefix)
	config := cosmos.AddressConfig{
		AccountAddressPrefix: cosmosAddrPrefix,
		AccountPubKeyPrefix:  cosmosAddrPrefix + "pub",
	}
	address := cosmos.GenerateBech32AddressFromPrivateKeyHex(privKeyHex, config)
	cmd.Println("Cosmos address:", address)
	return nil
}
