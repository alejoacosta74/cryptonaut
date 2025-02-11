package cmd

import (
	"encoding/hex"

	"github.com/alejoacosta74/cryptonaut/internal/config"
	"github.com/alejoacosta74/cryptonaut/pkg/crypto/ecdsa/secp256k1"
	"github.com/alejoacosta74/cryptonaut/pkg/crypto/schnorr"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var schnorrCmd = &cobra.Command{
	Use:   "schnorr",
	Short: "Schnorr signature",
}

var schnorrGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a Schnorr private key",
	RunE:  runSchnorrGenerate,
}

var schnorrPubkeyCmd = &cobra.Command{
	Use:   "pubkey",
	Short: "Generate a Schnorr public key from a private key",
	RunE:  runSchnorrPubKey,
	PreRun: func(cmd *cobra.Command, args []string) {
		cmd.MarkFlagRequired(config.FlagPrivateKey)
	},
}

var schnorrSignCmd = &cobra.Command{
	Use:   "sign",
	Short: "Sign a message using Schnorr",
	Args:  cobra.ExactArgs(1), // message to sign
	RunE:  runSchnorrSign,
	PreRun: func(cmd *cobra.Command, args []string) {
		cmd.MarkFlagRequired(config.FlagPrivateKey)
	},
}

var schnorrVerifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify a Schnorr signature",
	Args:  cobra.ExactArgs(1), // message to verify
	RunE:  runSchnorrVerify,
	PreRun: func(cmd *cobra.Command, args []string) {
		cmd.MarkFlagRequired(config.FlagPublicKey)
		cmd.MarkFlagRequired(config.FlagSignature)
	},
}

func init() {
	schnorrCmd.AddCommand(schnorrGenerateCmd)
	schnorrCmd.AddCommand(schnorrPubkeyCmd)
	schnorrCmd.AddCommand(schnorrSignCmd)
	schnorrCmd.AddCommand(schnorrVerifyCmd)

	rootCmd.AddCommand(schnorrCmd)
}

func runSchnorrGenerate(cmd *cobra.Command, args []string) error {
	privKey, err := secp256k1.GeneratePrivateKeyHex()
	if err != nil {
		return err
	}
	cmd.Println("Private key: ", privKey)
	return nil
}

func runSchnorrPubKey(cmd *cobra.Command, args []string) error {
	privKeyStr := viper.GetString(config.FlagPrivateKey)
	isCompressed := viper.GetBool(config.FlagPubKeyCompressed)
	privKey, err := secp256k1.ParsePrivateKeyFromString(privKeyStr)
	if err != nil {
		return err
	}

	pubKey := secp256k1.DerivePublicKey(privKey)
	pubKeyStr := secp256k1.SerializePublicKeyToString(pubKey, isCompressed)
	cmd.Println("Public key: ", pubKeyStr)
	return nil
}

func runSchnorrSign(cmd *cobra.Command, args []string) error {
	privKeyStr := viper.GetString(config.FlagPrivateKey)
	privKey, err := secp256k1.ParsePrivateKeyFromString(privKeyStr)
	if err != nil {
		return err
	}
	message := args[0]
	signature, err := schnorr.SignMessage(privKey, []byte(message))
	if err != nil {
		return err
	}
	cmd.Println("Signature: ", hex.EncodeToString(signature.Serialize()))
	return nil
}

func runSchnorrVerify(cmd *cobra.Command, args []string) error {
	pubKeyStr := viper.GetString(config.FlagPublicKey)
	pubKey, err := secp256k1.ParsePublicKeyFromString(pubKeyStr)
	if err != nil {
		return err
	}
	signatureStr := viper.GetString(config.FlagSignature)
	message := args[0]
	signature, err := hex.DecodeString(signatureStr)
	if err != nil {
		return err
	}
	valid, err := schnorr.VerifyMessage(pubKey.SerializeCompressed(), []byte(message), signature)
	if err != nil {
		return err
	}
	cmd.Println("Signature is valid: ", valid)
	return nil
}
