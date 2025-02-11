package cmd

import (
	"github.com/spf13/cobra"
)

var ethereumCmd = &cobra.Command{
	Use:   "ethereum",
	Short: "Ethereum commands",
	Long:  "Ethereum commands",
}

func init() {
	rootCmd.AddCommand(ethereumCmd)
}
