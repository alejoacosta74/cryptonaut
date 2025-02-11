package cmd

import (
	"encoding/hex"
	"fmt"

	"github.com/alejoacosta74/cryptonaut/internal/config"
	"github.com/alejoacosta74/cryptonaut/pkg/bitcoin"
	"github.com/alejoacosta74/cryptonaut/pkg/crypto/ecdsa/secp256k1"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var bitcoinCmd = &cobra.Command{
	Use:   "bitcoin",
	Short: "Bitcoin key and address generation",
	Long:  "Bitcoin key and address generation",
}

var bitcoinGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a Bitcoin private key",
	Long:  "Generate a Bitcoin private key",
	RunE:  runBitcoinGenerateCmd,
}

var bitcoinPubkeyCmd = &cobra.Command{
	Use:   "pubkey",
	Short: "Get the public key from a Bitcoin private key",
	Long:  "Get the public key from a Bitcoin private key",
	RunE:  runBitcoinPubkeyCmd,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Parent().Parent().MarkPersistentFlagRequired(config.FlagPrivateKey)
	},
}

var bitcoinAddressCmd = &cobra.Command{
	Use:   "address",
	Short: "Get the Bitcoin address from a public key",
	Long:  "Get the Bitcoin address from a public key",
	RunE:  runBitcoinAddressCmd,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Parent().Parent().MarkPersistentFlagRequired(config.FlagPrivateKey)
	},
}

var convertKeyCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert a Bitcoin private key from/to WIF and Hex",
	Long: `Convert a Bitcoin private key from/to WIF and Hex
	Usage:
	cryptonaut convert 6abd31bf5fe56e1aa5a49b8430a2bcaa276b4cd352b3d7072e89bb9a8a204cc1                 
	Private key: KzoCR5BTboQXqG9ah8HiHtigrK2DkrpgouYg94m4ZRWiCVEybGoy

	cryptonaut convert KzoCR5BTboQXqG9ah8HiHtigrK2DkrpgouYg94m4ZRWiCVEybGoy
	Private key: 6abd31bf5fe56e1aa5a49b8430a2bcaa276b4cd352b3d7072e89bb9a8a204cc1

	`,
	Args: cobra.ExactArgs(1),
	RunE: runConvertKey,
}

func init() {
	bitcoinCmd.AddCommand(bitcoinGenerateCmd)
	bitcoinCmd.AddCommand(bitcoinPubkeyCmd)
	bitcoinCmd.AddCommand(bitcoinAddressCmd)
	bitcoinCmd.AddCommand(convertKeyCmd)
	// flag to specify the format of the private key
	bitcoinGenerateCmd.Flags().StringP(config.FlagBitcoinFormat, "f", "hex", "Format of the private key (hex or wif)")
	viper.BindPFlag(config.FlagBitcoinFormat, bitcoinGenerateCmd.Flags().Lookup(config.FlagBitcoinFormat))
	// flat to specify the testnet flag
	bitcoinAddressCmd.Flags().BoolP(config.FlagTestnet, "t", false, "Use testnet")
	viper.BindPFlag(config.FlagTestnet, bitcoinAddressCmd.Flags().Lookup(config.FlagTestnet))

	rootCmd.AddCommand(bitcoinCmd)
}

func runBitcoinGenerateCmd(cmd *cobra.Command, args []string) error {
	format := viper.GetString(config.FlagBitcoinFormat)
	var privKey *btcec.PrivateKey
	var err error
	switch format {
	case "hex":
		privKey, err = secp256k1.GeneratePrivateKey()
		if err != nil {
			return err
		}
		cmd.Println("Private key:", hex.EncodeToString(privKey.Serialize()))
	case "wif":
		isTestnet := viper.GetBool(config.FlagTestnet)
		privKeyWIF, err := bitcoin.ConvertPrivateKeyToWIF(privKey, isTestnet, true)
		if err != nil {
			return err
		}
		cmd.Println("Private key:", privKeyWIF.String())
	default:
		return fmt.Errorf("invalid format: %s", format)
	}
	return nil
}

func runBitcoinPubkeyCmd(cmd *cobra.Command, args []string) error {
	privateKeyHex := viper.GetString(config.FlagPrivateKey)
	privKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return fmt.Errorf("invalid private key hex: %w", err)
	}

	privKey, pubKey := btcec.PrivKeyFromBytes(privKeyBytes)

	cmd.Println("Private key:", hex.EncodeToString(privKey.Serialize()))
	cmd.Println("Public key:", hex.EncodeToString(pubKey.SerializeCompressed()))

	return nil
}

func runBitcoinAddressCmd(cmd *cobra.Command, args []string) error {
	privateKeyHex := viper.GetString(config.FlagPrivateKey)
	isTestnet := viper.GetBool(config.FlagTestnet)
	address, err := bitcoin.GenerateAddressFromPrivateKey(privateKeyHex, isTestnet)
	if err != nil {
		return fmt.Errorf("invalid address: %w", err)
	}
	cmd.Println("Address:", address)
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
