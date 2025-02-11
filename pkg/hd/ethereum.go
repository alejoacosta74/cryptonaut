package hd

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

// CreateHDNode creates a new HD node for Ethereum from a mnemonic
func CreateEthereumHDNode(mnemonic string) (*bip32.Key, error) {
	// validate the mnemonic
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, fmt.Errorf("invalid mnemonic: '%s'", mnemonic)
	}

	// Step 2: Generate a seed from the mnemonic
	// The password here can be "" (empty) or a passphrase for additional security
	seed := bip39.NewSeed(mnemonic, "")

	// Step 3: Create a master key from the seed
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, fmt.Errorf("failed to create master key: %v", err)
	}

	// Step 4: Derive the Ethereum HD wallet path: m/44'/60'/0'/0/0
	// BIP44 Path for Ethereum:
	// - m/44'    : BIP44 purpose
	// - 60'      : Ethereum coin type
	// - 0'       : Account
	// - 0        : External chain
	// - 0        : Address index
	purposeKey, err := masterKey.NewChildKey(bip32.FirstHardenedChild + 44) // 44'
	if err != nil {
		return nil, fmt.Errorf("failed to derive purpose key: %v", err)
	}

	coinTypeKey, err := purposeKey.NewChildKey(bip32.FirstHardenedChild + 60) // 60'
	if err != nil {
		return nil, fmt.Errorf("failed to derive coin type key: %v", err)
	}

	accountKey, err := coinTypeKey.NewChildKey(bip32.FirstHardenedChild + 0) // 0'
	if err != nil {
		return nil, fmt.Errorf("failed to derive account key: %v", err)
	}

	changeKey, err := accountKey.NewChildKey(0) // 0
	if err != nil {
		return nil, fmt.Errorf("failed to derive change key: %v", err)
	}
	return changeKey, nil
}

// DerivePrivateKey derives a private key at a specific index
func DeriveEthereumPrivateKey(changeKey *bip32.Key, index uint32) (*ecdsa.PrivateKey, error) {
	addressKey, err := changeKey.NewChildKey(index)
	if err != nil {
		return nil, fmt.Errorf("failed to derive address key: %v", err)
	}

	privateKeyBytes := addressKey.Key
	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to private key: %v", err)
	}
	return privateKey, nil
}

// DerivePublicKey gets the public key from a private key
func DeriveEthereumPublicKey(changeKey *bip32.Key, index uint32) (*ecdsa.PublicKey, error) {
	privateKey, err := DeriveEthereumPrivateKey(changeKey, index)
	if err != nil {
		return nil, fmt.Errorf("failed to derive private key: %v", err)
	}

	publicKey, ok := privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("failed to get public key")
	}

	return publicKey, nil
}

// GetAddress derives Ethereum address from public key
func DeriveEthereumAddress(publicKey *ecdsa.PublicKey) common.Address {
	return crypto.PubkeyToAddress(*publicKey)
}
