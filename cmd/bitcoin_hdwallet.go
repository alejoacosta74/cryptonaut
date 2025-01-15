package cmd

import (
	"encoding/hex"
	"fmt"

	"github.com/alejoacosta74/cryptonaut/pkg/crypto"
	"github.com/alejoacosta74/cryptonaut/pkg/hdwallet"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var bitcoinHDWalletCmd = &cobra.Command{
	Use:   "hdwallet",
	Short: "HD wallet operations",
	Long: `Create and manage hierarchical deterministic wallets.
Supports BIP39 mnemonic generation and BIP32 key derivation.`,
}

var generateBitcoinCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a new mnemonic phrase",
	RunE:  runGenerate,
}

var deriveBitcoinCmd = &cobra.Command{
	Use:   "derive",
	Short: "Derive keys from a mnemonic phrase",
	Long: `Derive cryptographic keys from a BIP39 mnemonic phrase.
    
    This command implements BIP32 hierarchical deterministic wallets to derive
    child keys from a master seed. It supports Bitcoin BIP44 (m/44'/0'/0'/0/) derivation path.

    The derived keys can be used to generate addresses and sign transactions
    for the respective blockchain networks.

    Usage:
        bitcoin hdwallet derive --mnemonic "your mnemonic phrase" --index 0
        bitcoin hdwallet derive --mnemonic "your mnemonic phrase" --testnet --index 1

    The index parameter determines which child key to derive.
	If not specified, the first child key is derived.
	
	`,
	RunE: runDerive,
}

func init() {
	bitcoinHDWalletCmd.AddCommand(generateBitcoinCmd)
	bitcoinHDWalletCmd.AddCommand(deriveBitcoinCmd)

	// Flags for derive command
	deriveBitcoinCmd.Flags().String("mnemonic", "", "Mnemonic phrase")
	deriveBitcoinCmd.Flags().Int("index", 0, "Derivation index")
	viper.BindPFlag("bitcoin.hdwallet.mnemonic", deriveBitcoinCmd.Flags().Lookup("mnemonic"))
	viper.BindPFlag("bitcoin.hdwallet.index", deriveBitcoinCmd.Flags().Lookup("index"))

	bitcoinCmd.AddCommand(bitcoinHDWalletCmd)
}

func runGenerate(cmd *cobra.Command, args []string) error {
	mnemonic, err := crypto.GenerateMnemonic()
	if err != nil {
		return fmt.Errorf("failed to generate mnemonic: %v", err)
	}
	cmd.Println("Mnemonic:", mnemonic)
	return nil
}

func runDerive(cmd *cobra.Command, args []string) error {
	mnemonic := viper.GetString("bitcoin.hdwallet.mnemonic")
	index := viper.GetInt("bitcoin.hdwallet.index")
	testnet := viper.GetBool("testnet")

	var network *chaincfg.Params
	if testnet {
		network = &chaincfg.TestNet3Params
	} else {
		network = &chaincfg.MainNetParams
	}

	hdNode, err := hdwallet.CreateHDNode(mnemonic, network)
	if err != nil {
		return fmt.Errorf("failed to create hdnode: %v", err)
	}

	privKey, err := hdwallet.DerivePrivateKey(hdNode, hdwallet.DerivationIndex(index))
	if err != nil {
		cmd.PrintErrf("failed to derive private key: %v", err)
		return err
	}

	pubKey, err := hdwallet.DerivePublicKey(hdNode, hdwallet.DerivationIndex(index))
	if err != nil {
		cmd.PrintErrf("failed to derive public key: %v", err)
		return err
	}
	address, err := hdwallet.DeriveAddress(hdNode, hdwallet.DerivationIndex(index), network)
	if err != nil {
		cmd.PrintErrf("failed to derive address: %v", err)
		return err
	}

	cmd.Println("Private Key:", hex.EncodeToString(privKey.Serialize()))
	cmd.Println("Public Key:", hex.EncodeToString(pubKey.SerializeCompressed()))
	cmd.Println("Address:", address.EncodeAddress())
	return nil
}
