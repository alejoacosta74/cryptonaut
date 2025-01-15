/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/hex"
	"fmt"

	"github.com/alejoacosta74/cryptonaut/internal/config"
	"github.com/alejoacosta74/cryptonaut/pkg/bip44"
	"github.com/alejoacosta74/cryptonaut/pkg/crypto"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/spf13/cobra"
)

// hdwalletCmd represents the hdwallet command
var bip44DerivationCmd = &cobra.Command{
	Use:   "bip44",
	Short: "Derive keys from a mnemonic phrase",
	Long: `Create and manage hierarchical deterministic wallets.
Supports BIP39 mnemonic generation and BIP32 key derivation.

Usage:
    cryptonaut hdwallet derive bitcoin --mnemonic "your mnemonic phrase" --index 0 --testnet
    cryptonaut hdwallet derive ethereum --mnemonic "your mnemonic phrase"  --index 1

The index parameter determines which child key to derive.
If not specified, the first child key is derived.

	`,
}

var deriveBitcoinKeysCmd = &cobra.Command{
	Use:   "bitcoin",
	Short: "Derive keys from a mnemonic phrase for bitcoin",
	Long: `Derive cryptographic keys from a BIP39 mnemonic phrase for bitcoin.

    Usage:
        cryptonaut bip44 bitcoin --mnemonic "your mnemonic phrase" --index 0 --testnet
        cryptonaut bip44 bitcoin --mnemonic "your mnemonic phrase"  --index 1

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
        cryptonaut bip44 ethereum --mnemonic "your mnemonic phrase" --index 0
        cryptonaut bip44 ethereum --mnemonic "your mnemonic phrase"  --index 1

    The index parameter determines which child key to derive.
	If not specified, the first child key is derived.
	`,
	RunE: runDeriveEthereumKeysCmd,
}

func init() {
	bip44DerivationCmd.AddCommand(deriveBitcoinKeysCmd)
	bip44DerivationCmd.AddCommand(deriveEthereumKeysCmd)

	bip44DerivationCmd.PersistentFlags().String(config.FlagMnemonic, "", "Mnemonic phrase")
	bip44DerivationCmd.MarkPersistentFlagRequired(config.FlagMnemonic)
	bip44DerivationCmd.PersistentFlags().Int(config.FlagIndex, 0, "Derivation index")

	deriveBitcoinKeysCmd.PersistentFlags().Bool(config.FlagTestnet, false, "Use testnet")

	rootCmd.AddCommand(bip44DerivationCmd)

}

func runDeriveBitcoinKeysCmd(cmd *cobra.Command, args []string) error {
	mnemonic, err := cmd.Parent().PersistentFlags().GetString(config.FlagMnemonic)
	if err != nil {
		return fmt.Errorf("failed to get mnemonic: %v", err)
	}
	index, err := cmd.Parent().PersistentFlags().GetInt(config.FlagIndex)
	if err != nil {
		return fmt.Errorf("failed to get index: %v", err)
	}
	testnet, err := cmd.PersistentFlags().GetBool(config.FlagTestnet)
	if err != nil {
		return fmt.Errorf("failed to get testnet: %v", err)
	}

	var network *chaincfg.Params
	if testnet {
		network = &chaincfg.TestNet3Params
	} else {
		network = &chaincfg.MainNetParams
	}
	hdNode, err := bip44.CreateBitcoinHDNode(mnemonic, network)
	if err != nil {
		return fmt.Errorf("failed to create hdnode: %v", err)
	}

	privKey, err := bip44.DeriveBitcoinPrivateKey(hdNode, crypto.DerivationIndex(index))
	if err != nil {
		cmd.PrintErrf("failed to derive private key: %v", err)
		return err
	}

	pubKey, err := bip44.DeriveBitcoinPublicKey(hdNode, crypto.DerivationIndex(index))
	if err != nil {
		cmd.PrintErrf("failed to derive public key: %v", err)
		return err
	}
	address, err := bip44.DeriveBitcoinAddress(hdNode, crypto.DerivationIndex(index), network)
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
	mnemonic, err := cmd.Parent().PersistentFlags().GetString(config.FlagMnemonic)
	if err != nil {
		return fmt.Errorf("failed to get mnemonic: %v", err)
	}
	index, err := cmd.Parent().PersistentFlags().GetInt(config.FlagIndex)
	if err != nil {
		return fmt.Errorf("failed to get index: %v", err)
	}

	hdNode, err := bip44.CreateEthereumHDNode(mnemonic)
	if err != nil {
		return fmt.Errorf("failed to create hdnode: %v", err)
	}

	privKey, err := bip44.DeriveEthereumPrivateKey(hdNode, uint32(index))
	if err != nil {
		return fmt.Errorf("failed to derive private key: %v", err)
	}

	pubKey, err := bip44.DeriveEthereumPublicKey(hdNode, uint32(index))
	if err != nil {
		return fmt.Errorf("failed to derive public key: %v", err)
	}

	address := bip44.DeriveEthereumAddress(pubKey)

	cmd.Println("Private Key:", hex.EncodeToString(privKey.D.Bytes()))
	cmd.Println("Public Key:", hex.EncodeToString(pubKey.X.Bytes()))
	cmd.Println("Address:", address.Hex())
	return nil
}
