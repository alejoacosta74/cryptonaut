package hd

import (
	"fmt"

	"github.com/alejoacosta74/cryptonaut/pkg/crypto"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/tyler-smith/go-bip39"
)

func CreateBitcoinHDNode(mnemonic string, network *chaincfg.Params) (*hdkeychain.ExtendedKey, error) {

	// validate the mnemonic
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, fmt.Errorf("invalid mnemonic: '%s'", mnemonic)
	}

	// derive the hdnode
	seed := bip39.NewSeed(mnemonic, "")

	// Create the master key (root key)
	masterKey, err := hdkeychain.NewMaster(seed, network) // Use TestNetParams for testnet
	if err != nil {
		return nil, fmt.Errorf("failed to create master key: %v", err)
	}

	// Derive the path m/44'/0'/0'/0/
	// Hardened key derivations are offset by HardenedKeyStart (2^31)
	const HardenedOffset uint32 = hdkeychain.HardenedKeyStart

	// Derive 44' (Purpose)
	purposeKey, err := masterKey.Derive(crypto.PurposePath.ToUint32() + HardenedOffset)
	if err != nil {
		return nil, fmt.Errorf("failed to derive purpose key: %v", err)
	}

	// Derive 0' (Coin type for Bitcoin)
	coinTypeKey, err := purposeKey.Derive(crypto.BitcoinCoinTypePath.ToUint32() + HardenedOffset)
	if err != nil {
		return nil, fmt.Errorf("failed to derive coin type key: %v", err)
	}

	// Derive 0' (Account)
	accountKey, err := coinTypeKey.Derive(crypto.AccountPath.ToUint32() + HardenedOffset)
	if err != nil {
		return nil, fmt.Errorf("failed to derive account key: %v", err)
	}

	// Derive 0 (External chain)
	chainKey, err := accountKey.Derive(crypto.ChainPath.ToUint32())
	if err != nil {
		return nil, fmt.Errorf("failed to derive chain key: %v", err)
	}

	return chainKey, nil
}

func DeriveBitcoinPrivateKey(hdNode *hdkeychain.ExtendedKey, index crypto.DerivationIndex) (*btcec.PrivateKey, error) {
	// derive the child key at the given index
	child, err := hdNode.Derive(index.ToUint32())
	if err != nil {
		return nil, fmt.Errorf("failed to derive child for index %d: %v", index, err)
	}

	// get the private key
	privKey, err := child.ECPrivKey()
	if err != nil {
		return nil, fmt.Errorf("failed to derive private key: %v", err)
	}

	return privKey, nil
}

func DeriveBitcoinPublicKey(hdNode *hdkeychain.ExtendedKey, index crypto.DerivationIndex) (*btcec.PublicKey, error) {
	child, err := hdNode.Derive(index.ToUint32())
	if err != nil {
		return nil, fmt.Errorf("failed to derive child for index %d: %v", index, err)
	}

	pubKey, err := child.ECPubKey()
	if err != nil {
		return nil, fmt.Errorf("failed to derive public key: %v", err)
	}

	return pubKey, nil
}

func DeriveBitcoinAddress(hdNode *hdkeychain.ExtendedKey, index crypto.DerivationIndex, network *chaincfg.Params) (*btcutil.AddressPubKeyHash, error) {
	child, err := hdNode.Derive(index.ToUint32())
	if err != nil {
		return nil, fmt.Errorf("failed to derive child for index %d: %v", index, err)
	}

	address, err := child.Address(network)
	if err != nil {
		return nil, fmt.Errorf("failed to derive address: %v", err)
	}

	return address, nil
}
