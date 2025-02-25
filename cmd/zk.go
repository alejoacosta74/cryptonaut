package cmd

import (
	"fmt"

	"github.com/alejoacosta74/cryptonaut/internal/config"
	"github.com/alejoacosta74/cryptonaut/pkg/zk/snarkproof"
	"github.com/alejoacosta74/go-logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var zkCmd = &cobra.Command{
	Use:   "zk",
	Short: "Zero-knowledge proof",
}

var zkSnarkCmd = &cobra.Command{
	Use:   "snark",
	Short: "Zero-knowledge proof using snark",
}

var zkSnarkProveCmd = &cobra.Command{
	Use:   "prove",
	Short: "Prove a statement using snark",
	Long: `Prove a statement using snark using Groth16 algorithm as the backend
Example:
cryptonaut zk snark prove --circuit age --birth-year 1990
`,
	RunE: runZkSnarkProve,
}

var zkSnarkVerifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify a statement using snark",
	Long: `Verify a statement using snark using Groth16 algorithm as the backend
Example:
cryptonaut zk snark verify --circuit age --birth-year 1990 --proof proof.data --vk vk.data
`,
	RunE: runZkSnarkVerify,
}

func init() {
	zkCmd.AddCommand(zkSnarkCmd)

	zkSnarkCmd.PersistentFlags().StringP(config.FlagCircuit, "c", "", "Circuit to use")
	zkSnarkCmd.MarkPersistentFlagRequired(config.FlagCircuit)
	viper.BindPFlag(config.FlagCircuit, zkSnarkCmd.PersistentFlags().Lookup(config.FlagCircuit))
	zkSnarkCmd.PersistentFlags().StringP(config.FlagProof, "p", "proof.data", "Filename for the proof")
	viper.BindPFlag(config.FlagProof, zkSnarkCmd.PersistentFlags().Lookup(config.FlagProof))
	zkSnarkCmd.PersistentFlags().StringP(config.FlagVk, "v", "vk.data", "Filename for the verification key")
	viper.BindPFlag(config.FlagVk, zkSnarkCmd.PersistentFlags().Lookup(config.FlagVk))

	zkSnarkProveCmd.Flags().IntP(config.FlagBirthYear, "b", 0, "Birth year of the user")
	viper.BindPFlag(config.FlagBirthYear, zkSnarkProveCmd.Flags().Lookup(config.FlagBirthYear))

	zkSnarkCmd.AddCommand(zkSnarkProveCmd)
	zkSnarkCmd.AddCommand(zkSnarkVerifyCmd)

	rootCmd.AddCommand(zkCmd)
}

func runZkSnarkProve(cmd *cobra.Command, args []string) error {
	circuit := viper.GetString("circuit")
	switch circuit {
	case "age":
		birthYear := viper.GetInt("birth-year")
		err := snarkproof.GenerateAgeProof(birthYear)
		if err != nil {
			return fmt.Errorf("Failed to generate proof: %v", err)
		}
		logger.Infof("✅ Proof successfully generated and saved in file %s, and verification key in file %s \n", viper.GetString("proof"), viper.GetString("vk"))
	default:
		return fmt.Errorf("Circuit %s not supported", circuit)
	}
	return nil
}

func runZkSnarkVerify(cmd *cobra.Command, args []string) error {
	circuit := viper.GetString("circuit")
	proofFile := viper.GetString("proof")
	vkFile := viper.GetString("vk")
	switch circuit {
	case "age":
		logger.Infof("🔍 Verifying proof from file %s using verification key from file %s", proofFile, vkFile)
		err := snarkproof.VerifyAgeProof(proofFile, vkFile)
		if err != nil {
			return fmt.Errorf("Failed to verify proof: %v", err)
		}
		logger.Infof("✅ Proof verification result: %v\n", true)
	}
	return nil
}
