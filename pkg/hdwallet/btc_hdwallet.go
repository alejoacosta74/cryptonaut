package hdwallet

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/tyler-smith/go-bip39"
)

func CreateHDNode(mnemonic string, network *chaincfg.Params) (*hdkeychain.ExtendedKey, error) {

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
	purposeKey, err := masterKey.Child(PurposePath.ToUint32() + HardenedOffset)
	if err != nil {
		return nil, fmt.Errorf("failed to derive purpose key: %v", err)
	}

	// Derive 0' (Coin type for Bitcoin)
	coinTypeKey, err := purposeKey.Child(BitcoinCoinTypePath.ToUint32() + HardenedOffset)
	if err != nil {
		return nil, fmt.Errorf("failed to derive coin type key: %v", err)
	}

	// Derive 0' (Account)
	accountKey, err := coinTypeKey.Child(AccountPath.ToUint32() + HardenedOffset)
	if err != nil {
		return nil, fmt.Errorf("failed to derive account key: %v", err)
	}

	// Derive 0 (External chain)
	chainKey, err := accountKey.Child(ChainPath.ToUint32())
	if err != nil {
		return nil, fmt.Errorf("failed to derive chain key: %v", err)
	}

	return chainKey, nil
}

func DerivePrivateKey(hdNode *hdkeychain.ExtendedKey, index DerivationIndex) (*btcec.PrivateKey, error) {
	// derive the child key at the given index
	child, err := hdNode.Child(index.ToUint32())
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

func DerivePublicKey(hdNode *hdkeychain.ExtendedKey, index DerivationIndex) (*btcec.PublicKey, error) {
	child, err := hdNode.Child(index.ToUint32())
	if err != nil {
		return nil, fmt.Errorf("failed to derive child for index %d: %v", index, err)
	}

	pubKey, err := child.ECPubKey()
	if err != nil {
		return nil, fmt.Errorf("failed to derive public key: %v", err)
	}

	return pubKey, nil
}

func DeriveAddress(hdNode *hdkeychain.ExtendedKey, index DerivationIndex, network *chaincfg.Params) (*btcutil.AddressPubKeyHash, error) {
	child, err := hdNode.Child(index.ToUint32())
	if err != nil {
		return nil, fmt.Errorf("failed to derive child for index %d: %v", index, err)
	}

	address, err := child.Address(network)
	if err != nil {
		return nil, fmt.Errorf("failed to derive address: %v", err)
	}

	return address, nil
}
