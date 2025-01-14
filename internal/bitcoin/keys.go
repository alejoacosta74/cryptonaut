package bitcoin

import (
	"encoding/hex"
	"strings"

	"github.com/alejoacosta74/go-logger"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
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

func getChainParams(testnet bool) *chaincfg.Params {
	if testnet {
		return &chaincfg.TestNet3Params
	}
	return &chaincfg.MainNetParams
}

func isWIF(privKey string) bool {
	// WIF must be base58 encoded and have a minimum length
	if len(privKey) < 51 || len(privKey) > 52 {
		return false
	}

	// Attempt to decode WIF - this will validate base58 format and checksum
	wif, err := btcutil.DecodeWIF(privKey)
	if err != nil {
		return false
	}

	// Ensure the decoded key is valid
	return wif != nil && wif.IsForNet(&chaincfg.MainNetParams) || wif.IsForNet(&chaincfg.TestNet3Params)
}

func isHex(privKey string) bool {
	cleanKey := strings.TrimPrefix(privKey, "0x")
	if len(cleanKey) != 64 {
		return false
	}

	b, err := hex.DecodeString(cleanKey)
	if err != nil {
		return false
	}

	priv, _ := btcec.PrivKeyFromBytes(btcec.S256(), b)
	return priv != nil
}
