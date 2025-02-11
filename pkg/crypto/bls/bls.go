package bls

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/herumi/bls-eth-go-binary/bls"
)

func GeneratePrivateKey() (*bls.SecretKey, error) {
	if err := bls.Init(bls.BLS12_381); err != nil {
		return nil, err
	}
	var sec bls.SecretKey
	sec.SetByCSPRNG()
	return &sec, nil
}

func ParsePrivateKeyFromString(privateKeyStr string) (*bls.SecretKey, error) {
	if err := bls.Init(bls.BLS12_381); err != nil {
		return nil, err
	}
	var sec bls.SecretKey
	privateKeyBytes, err := hex.DecodeString(privateKeyStr)
	if err != nil {
		return nil, err
	}
	sec.Deserialize(privateKeyBytes)
	return &sec, nil
}

func DerivePublicKey(privateKey *bls.SecretKey) *bls.PublicKey {
	if err := bls.Init(bls.BLS12_381); err != nil {
		return nil
	}
	return privateKey.GetPublicKey()
}

func ParsePublicKeyFromString(publicKeyStr string) (*bls.PublicKey, error) {
	if err := bls.Init(bls.BLS12_381); err != nil {
		return nil, err
	}
	var pub bls.PublicKey
	publicKeyBytes, err := hex.DecodeString(publicKeyStr)
	if err != nil {
		return nil, err
	}
	pub.Deserialize(publicKeyBytes)
	return &pub, nil
}

func SignMessage(privateKey *bls.SecretKey, message string) (*bls.Sign, error) {
	if err := bls.Init(bls.BLS12_381); err != nil {
		return nil, err
	}
	digest := sha256.Sum256([]byte(message))
	return privateKey.SignByte(digest[:]), nil

}

func ParseSignatureFromString(signatureStr string) (*bls.Sign, error) {
	if err := bls.Init(bls.BLS12_381); err != nil {
		return nil, err
	}
	var sig bls.Sign
	signatureBytes, err := hex.DecodeString(signatureStr)
	if err != nil {
		return nil, err
	}
	sig.Deserialize(signatureBytes)
	return &sig, nil
}

func VerifyMessage(publicKey *bls.PublicKey, message string, signature *bls.Sign) (bool, error) {
	if err := bls.Init(bls.BLS12_381); err != nil {
		return false, err
	}
	digest := sha256.Sum256([]byte(message))
	return signature.VerifyByte(publicKey, digest[:]), nil
}
