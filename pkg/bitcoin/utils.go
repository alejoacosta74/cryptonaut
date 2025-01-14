package bitcoin

import (
	"encoding/hex"
	"strings"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
)

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
