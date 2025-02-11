/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/hex"
	"fmt"

	"github.com/alejoacosta74/cryptonaut/internal/config"
	"github.com/alejoacosta74/cryptonaut/pkg/crypto"
	"github.com/alejoacosta74/cryptonaut/pkg/hd"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// hdwalletCmd represents the hdwallet command
var hdDerivationCmd = &cobra.Command{
	Use:   "hd",
	Short: "Derive keys from a mnemonic phrase",
	Long: `Create and manage hierarchical deterministic wallets.
Supports BIP39 mnemonic generation and BIP32 key derivation.

Usage:
    cryptonaut hd derive bitcoin --mnemonic "your mnemonic phrase" --index 0 --testnet
    cryptonaut hd derive ethereum --mnemonic "your mnemonic phrase"  --index 1

The index parameter determines which child key to derive.
If not specified, the first child key is derived.

	`,
}

var generateMnemonicCmd = &cobra.Command{
	Use:   "mnemonic",
	Short: "Generate a new mnemonic phrase",
	Long: `Generate a new mnemonic phrase using BIP39
	
	Usage:
		cryptonaut hd mnemonic
	`,
	RunE: runGenerateMnemonicCmd,
}

var deriveBitcoinKeysCmd = &cobra.Command{
	Use:   "bitcoin",
	Short: "Derive keys from a mnemonic phrase for bitcoin",
	Long: `Derive cryptographic keys from a BIP39 mnemonic phrase for bitcoin.

    Usage:
        cryptonaut hd bitcoin --mnemonic "your mnemonic phrase" --index 0 --testnet
        cryptonaut hd bitcoin --mnemonic "your mnemonic phrase"  --index 1

    The index parameter determines which child key to derive.
	If not specified, the first child key is derived.
	`,
	RunE: runDeriveBitcoinKeysCmd,
}

var deriveEthereumKeysCmd = &cobra.Command{
	Use:   "ethereum",
	Short: "Derive keys from a mnemonic phrase for ethereum",
	Long: `Derive cryptographic keys from a BIP39 mnemonic phrase for ethereum.

    This command implements BIP32 hierarchical deterministic wallets to derive
    child keys from a master seed. It supports Ethereum BIP44 (m/44'/60'/0'/0) derivation path.

    Usage:
        cryptonaut hd ethereum --mnemonic "your mnemonic phrase" --index 0
        cryptonaut hd ethereum --mnemonic "your mnemonic phrase"  --index 1

    The index parameter determines which child key to derive.
	If not specified, the first child key is derived.
	`,
	RunE: runDeriveEthereumKeysCmd,
}

func init() {
	hdDerivationCmd.AddCommand(deriveBitcoinKeysCmd)
	hdDerivationCmd.AddCommand(deriveEthereumKeysCmd)
	hdDerivationCmd.AddCommand(generateMnemonicCmd)

	hdDerivationCmd.PersistentFlags().String(config.FlagMnemonic, "", "Mnemonic phrase")
	viper.BindPFlag(config.FlagMnemonic, hdDerivationCmd.PersistentFlags().Lookup(config.FlagMnemonic))
	deriveBitcoinKeysCmd.MarkPersistentFlagRequired(config.FlagMnemonic)
	deriveEthereumKeysCmd.MarkPersistentFlagRequired(config.FlagMnemonic)

	hdDerivationCmd.PersistentFlags().Int(config.FlagIndex, 0, "Derivation index")
	viper.BindPFlag(config.FlagIndex, hdDerivationCmd.PersistentFlags().Lookup(config.FlagIndex))

	rootCmd.AddCommand(hdDerivationCmd)
	rootCmd.AddCommand(hdDerivationCmd)

}

func runDeriveBitcoinKeysCmd(cmd *cobra.Command, args []string) error {
	mnemonic := viper.GetString(config.FlagMnemonic)
	index := viper.GetInt(config.FlagIndex)
	testnet := viper.GetBool(config.FlagTestnet)

	var network *chaincfg.Params
	if testnet {
		network = &chaincfg.TestNet3Params
	} else {
		network = &chaincfg.MainNetParams
	}
	hdNode, err := hd.CreateBitcoinHDNode(mnemonic, network)
	if err != nil {
		return fmt.Errorf("failed to create hdnode: %v", err)
	}

	privKey, err := hd.DeriveBitcoinPrivateKey(hdNode, crypto.DerivationIndex(index))
	if err != nil {
		cmd.PrintErrf("failed to derive private key: %v", err)
		return err
	}

	pubKey, err := hd.DeriveBitcoinPublicKey(hdNode, crypto.DerivationIndex(index))
	if err != nil {
		cmd.PrintErrf("failed to derive public key: %v", err)
		return err
	}
	address, err := hd.DeriveBitcoinAddress(hdNode, crypto.DerivationIndex(index), network)
	if err != nil {
		cmd.PrintErrf("failed to derive address: %v", err)
		return err
	}

	cmd.Println("Private Key:", hex.EncodeToString(privKey.Serialize()))
	cmd.Println("Public Key:", hex.EncodeToString(pubKey.SerializeCompressed()))
	cmd.Println("Address:", address.EncodeAddress())
	return nil

}

func runDeriveEthereumKeysCmd(cmd *cobra.Command, args []string) error {
	mnemonic := viper.GetString(config.FlagMnemonic)
	index := viper.GetInt(config.FlagIndex)

	hdNode, err := hd.CreateEthereumHDNode(mnemonic)
	if err != nil {
		return fmt.Errorf("failed to create hdnode: %v", err)
	}

	privKey, err := hd.DeriveEthereumPrivateKey(hdNode, uint32(index))
	if err != nil {
		return fmt.Errorf("failed to derive private key: %v", err)
	}

	pubKey, err := hd.DeriveEthereumPublicKey(hdNode, uint32(index))
	if err != nil {
		return fmt.Errorf("failed to derive public key: %v", err)
	}

	address := hd.DeriveEthereumAddress(pubKey)

	cmd.Println("Private Key:", hex.EncodeToString(privKey.D.Bytes()))
	cmd.Println("Public Key:", hex.EncodeToString(pubKey.X.Bytes()))
	cmd.Println("Address:", address.Hex())
	return nil
}

func runGenerateMnemonicCmd(cmd *cobra.Command, args []string) error {
	mnemonic, err := crypto.GenerateMnemonic()
	if err != nil {
		return fmt.Errorf("failed to generate mnemonic: %v", err)
	}
	cmd.Println("Mnemonic:", mnemonic)
	return nil
}
