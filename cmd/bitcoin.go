/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// bitcoinCmd represents the bitcoin command
var bitcoinCmd = &cobra.Command{
	Use:   "bitcoin",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("bitcoin called")
	},
}

func init() {
	bitcoinCmd.PersistentFlags().BoolP("testnet", "t", false, "Use testnet")
	bitcoinCmd.PersistentFlags().Bool("compressed", false, "Use compressed key")
	viper.BindPFlag("testnet", bitcoinCmd.PersistentFlags().Lookup("testnet"))
	viper.BindPFlag("compressed", bitcoinCmd.PersistentFlags().Lookup("compressed"))

	rootCmd.AddCommand(bitcoinCmd)
}
