/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/alejoacosta74/cryptonaut/pkg/bitcoin"
	"github.com/spf13/cobra"
)

var convertKeyCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert a Bitcoin private key from/to WIF and Hex",
	Long: `Convert a Bitcoin private key from/to WIF and Hex
	Usage:
	cryptonaut convert 6abd31bf5fe56e1aa5a49b8430a2bcaa276b4cd352b3d7072e89bb9a8a204cc1                 
	Private key: KzoCR5BTboQXqG9ah8HiHtigrK2DkrpgouYg94m4ZRWiCVEybGoy

	cryptonaut convert KzoCR5BTboQXqG9ah8HiHtigrK2DkrpgouYg94m4ZRWiCVEybGoy
	Private key: 6abd31bf5fe56e1aa5a49b8430a2bcaa276b4cd352b3d7072e89bb9a8a204cc1

	`,
	Args: cobra.ExactArgs(1),
	RunE: runConvertKey,
}

func init() {
	rootCmd.AddCommand(convertKeyCmd)
}

func runConvertKey(cmd *cobra.Command, args []string) error {
	privKey := args[0]
	result, err := bitcoin.ConvertKey(privKey)
	if err != nil {
		cmd.PrintErrln(err)
		return err
	}
	cmd.Printf("Private key: %s\n", result)
	return nil
}
