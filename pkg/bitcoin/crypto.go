package bitcoin

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
)

func ConvertPrivateKeyToWIF(privKey *btcec.PrivateKey, testnet, compressed bool) (*btcutil.WIF, error) {
	net := getChainParams(testnet)
	privKeyWIF, err := btcutil.NewWIF(privKey, net, compressed)
	if err != nil {
		return nil, fmt.Errorf("failed to create private key: %v", err)
	}
	return privKeyWIF, nil
}

func GeneratePublicKeyFromWIF(wif *btcutil.WIF) ([]byte, error) {
	return wif.PrivKey.PubKey().SerializeCompressed(), nil
}

func GenerateAddressFromPrivateKey(privKey string, testnet bool) (string, error) {
	net := getChainParams(testnet)
	var key *btcec.PrivateKey

	if isWIF(privKey) {
		wif, _ := btcutil.DecodeWIF(privKey)
		key = wif.PrivKey
	} else if isHex(privKey) {
		b, _ := hex.DecodeString(strings.TrimPrefix(privKey, "0x"))
		key, _ = btcec.PrivKeyFromBytes(b)
	} else {
		return "", fmt.Errorf("invalid private key format (must be WIF or hex)")
	}

	addr, err := btcutil.NewAddressPubKey(key.PubKey().SerializeCompressed(), net)
	if err != nil {
		return "", fmt.Errorf("failed to generate address: %v", err)
	}
	return addr.EncodeAddress(), nil
}

func ConvertKey(privKey string) (string, error) {
	if isWIF(privKey) {
		wif, _ := btcutil.DecodeWIF(privKey)
		return hex.EncodeToString(wif.PrivKey.Serialize()), nil
	} else if isHex(privKey) {
		privKey, err := hex.DecodeString(privKey)
		if err != nil {
			return "", err
		}
		priv, _ := btcec.PrivKeyFromBytes(privKey)
		wif, err := btcutil.NewWIF(priv, getChainParams(false), true)
		if err != nil {
			return "", err
		}
		return wif.String(), nil
	}
	return "", fmt.Errorf("invalid private key format (must be WIF or hex)")
}
