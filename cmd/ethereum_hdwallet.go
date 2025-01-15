package cmd

import (
	"encoding/hex"
	"fmt"

	"github.com/alejoacosta74/cryptonaut/pkg/ethereum"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ethereumHDWalletCmd = &cobra.Command{
	Use:   "hdwallet",
	Short: "HD Wallet operations",
	Long:  "Create and manage hierarchical deterministic wallets. Supports BIP39 mnemonic generation and BIP32 key derivation.",
}

var deriveEthereumCmd = &cobra.Command{
	Use:   "derive",
	Short: "Derive keys from a mnemonic phrase",
	Long: `Derive cryptographic keys from a BIP39 mnemonic phrase.
    
    This command implements BIP32 hierarchical deterministic wallets to derive
    child keys from a master seed. It supports Ethereum BIP44 (m/44'/60'/0'/0) derivation path.

    The derived keys can be used to generate addresses and sign transactions
    for the respective blockchain networks.

    Usage:
        ethereum hdwallet derive --mnemonic "your mnemonic phrase" --index 0

    The index parameter determines which child key to derive.
	If not specified, the first child key is derived.
	`,
	RunE: runDeriveEthereumCmd,
}

func init() {
	ethereumHDWalletCmd.AddCommand(deriveEthereumCmd)

	// Flags for derive command
	deriveEthereumCmd.Flags().String("mnemonic", "", "Mnemonic phrase")
	deriveEthereumCmd.Flags().Int("index", 0, "Derivation index")
	viper.BindPFlag("ethereum.hdwallet.mnemonic", deriveEthereumCmd.Flags().Lookup("mnemonic"))
	viper.BindPFlag("ethereum.hdwallet.index", deriveEthereumCmd.Flags().Lookup("index"))

	ethereumCmd.AddCommand(ethereumHDWalletCmd)
}

func runDeriveEthereumCmd(cmd *cobra.Command, args []string) error {
	mnemonic := viper.GetString("ethereum.hdwallet.mnemonic")
	index := viper.GetInt("ethereum.hdwallet.index")

	hdNode, err := ethereum.CreateHDNode(mnemonic)
	if err != nil {
		return fmt.Errorf("failed to create hdnode: %v", err)
	}

	privKey, err := ethereum.DerivePrivateKey(hdNode, uint32(index))
	if err != nil {
		return fmt.Errorf("failed to derive private key: %v", err)
	}

	pubKey, err := ethereum.DerivePublicKey(hdNode, uint32(index))
	if err != nil {
		return fmt.Errorf("failed to derive public key: %v", err)
	}

	address := ethereum.GetAddress(pubKey)

	cmd.Println("Private Key:", hex.EncodeToString(privKey.D.Bytes()))
	cmd.Println("Public Key:", hex.EncodeToString(pubKey.X.Bytes()))
	cmd.Println("Address:", address.Hex())
	return nil
}
