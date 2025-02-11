package cmd

import (
	"math/big"

	"github.com/alejoacosta74/cryptonaut/internal/config"
	"github.com/alejoacosta74/cryptonaut/pkg/crypto/ecdsa/secp256r1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ecdsaCmd = &cobra.Command{
	Use:   "ecdsa",
	Short: "ECDSA commands",
	Long:  "ECDSA commands",
}

var ecdsaGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a ECDSA private key in curve P-256",
	Long:  "Generate a ECDSA private key in curve P-256",
	RunE:  runEcdsaGenerateCmd,
}

var ecdsaPubkeyCmd = &cobra.Command{
	Use:   "pubkey",
	Short: "Derive a ECDSA public key from a private key",
	Long:  "Derive a ECDSA public key from a private key",
	RunE:  runEcdsaPubkeyCmd,
	PreRun: func(cmd *cobra.Command, args []string) {
		// mark private key as required
		cmd.MarkFlagRequired(config.FlagPrivateKey)
	},
}

var ecdsaSignCmd = &cobra.Command{
	Use:   "sign",
	Short: "Sign a message with a ECDSA private key",
	Long:  "Sign a message with a ECDSA private key",
	RunE:  runEcdsaSignCmd,
	Args:  cobra.ExactArgs(1), // the message to sign
	PreRun: func(cmd *cobra.Command, args []string) {
		cmd.MarkFlagRequired(config.FlagPrivateKey)
	},
}

var ecdsaVerifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify a signature with r and s signature components",
	Long:  "Verify a signature with r and s signature components",
	RunE:  runEcdsaVerifyCmd,
	Args:  cobra.ExactArgs(1), // the message to verify
}

func init() {
	ecdsaCmd.AddCommand(ecdsaGenerateCmd)
	ecdsaCmd.AddCommand(ecdsaPubkeyCmd)
	ecdsaCmd.AddCommand(ecdsaSignCmd)
	ecdsaCmd.AddCommand(ecdsaVerifyCmd)
	ecdsaVerifyCmd.Flags().StringP(config.FlagSignatureR, "r", "", "R signature component")
	ecdsaVerifyCmd.MarkFlagRequired(config.FlagSignatureR)
	viper.BindPFlag(config.FlagSignatureR, ecdsaVerifyCmd.Flags().Lookup(config.FlagSignatureR))
	ecdsaVerifyCmd.Flags().StringP(config.FlagSignatureS, "s", "", "S signature component")
	ecdsaVerifyCmd.MarkFlagRequired(config.FlagSignatureS)
	viper.BindPFlag(config.FlagSignatureS, ecdsaVerifyCmd.Flags().Lookup(config.FlagSignatureS))

	rootCmd.AddCommand(ecdsaCmd)
}

func runEcdsaGenerateCmd(cmd *cobra.Command, args []string) error {
	privateKey, err := secp256r1.GeneratePrivateKey()
	if err != nil {
		return err
	}
	privKeyStr, err := secp256r1.SerializePrivateKey(privateKey)
	if err != nil {
		return err
	}
	cmd.Println("Private key:", privKeyStr)
	return nil
}

func runEcdsaPubkeyCmd(cmd *cobra.Command, args []string) error {
	privKeyStr := viper.GetString(config.FlagPrivateKey)

	privKey, err := secp256r1.ParseECDSAPrivateKeyFromHex(privKeyStr)
	if err != nil {
		return err
	}

	pubKey, err := secp256r1.DerivePublicKey(privKey)
	if err != nil {
		return err
	}
	cmd.Println("Public key:", secp256r1.SerializePublicKey(pubKey))

	return nil
}

func runEcdsaSignCmd(cmd *cobra.Command, args []string) error {
	privKeyStr := viper.GetString(config.FlagPrivateKey)
	privKey, err := secp256r1.ParseECDSAPrivateKeyFromHex(privKeyStr)
	if err != nil {
		return err
	}
	message := args[0]

	r, s, err := secp256r1.SignMessage(privKey, []byte(message))
	if err != nil {
		return err
	}
	cmd.Printf("Signature: r=%x\n, s=%x\n", r, s)
	return nil
}

func runEcdsaVerifyCmd(cmd *cobra.Command, args []string) error {
	r := viper.GetString(config.FlagSignatureR)
	s := viper.GetString(config.FlagSignatureS)
	pubKeyStr := viper.GetString(config.FlagPublicKey)
	pubKey, err := secp256r1.ParseECDSAPublicKeyFromHex(pubKeyStr)
	if err != nil {
		return err
	}
	message := args[0]

	// convert r and s to big.Int
	rInt, _ := new(big.Int).SetString(r, 16)
	sInt, _ := new(big.Int).SetString(s, 16)

	valid, err := secp256r1.VerifySignature(pubKey, []byte(message), rInt, sInt)
	if err != nil {
		return err
	}
	cmd.Println("Signature is valid:", valid)
	return nil
}
