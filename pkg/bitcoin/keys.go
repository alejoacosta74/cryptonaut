package bitcoin

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/alejoacosta74/go-logger"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcutil"
)

func GeneratePrivateKeyHex() *btcec.PrivateKey {
	privKey, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		logger.Fatalf("Failed to create private key: %v", err)
	}
	return privKey
}

func GeneratePrivateKeyWIF(testnet, compressed bool) *btcutil.WIF {
	privKey := GeneratePrivateKeyHex()
	net := getChainParams(testnet)
	privKeyWIF, err := btcutil.NewWIF(privKey, net, compressed)
	if err != nil {
		logger.Fatalf("Failed to create private key: %v", err)
	}
	return privKeyWIF
}

func DerivePublicKey(wif *btcutil.WIF) []byte {
	return wif.PrivKey.PubKey().SerializeCompressed()
}

func GenerateAddress(privKey string, testnet bool) string {
	net := getChainParams(testnet)
	var key *btcec.PrivateKey

	if isWIF(privKey) {
		wif, _ := btcutil.DecodeWIF(privKey)
		key = wif.PrivKey
	} else if isHex(privKey) {
		b, _ := hex.DecodeString(strings.TrimPrefix(privKey, "0x"))
		key, _ = btcec.PrivKeyFromBytes(btcec.S256(), b)
	} else {
		logger.Fatal("Invalid private key format. Must be WIF or hex.")
	}

	addr, err := btcutil.NewAddressPubKey(key.PubKey().SerializeCompressed(), net)
	if err != nil {
		logger.Fatalf("Failed to generate address: %v", err)
	}
	return addr.EncodeAddress()
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
		priv, _ := btcec.PrivKeyFromBytes(btcec.S256(), privKey)
		wif, err := btcutil.NewWIF(priv, getChainParams(false), true)
		if err != nil {
			return "", err
		}
		return wif.String(), nil
	}
	return "", fmt.Errorf("invalid private key format (must be WIF or hex)")
}
