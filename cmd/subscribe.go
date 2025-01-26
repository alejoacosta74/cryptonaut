/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/alejoacosta74/cryptonaut/pkg/clients"
	"github.com/alejoacosta74/cryptonaut/pkg/subscription"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// subscribeCmd represents the subscribe command
var subscribeCmd = &cobra.Command{
	Use:   "subscribe",
	Short: "Initiates a subscription to a blockchain",
	Long: `A longer description that spans multiple lines and likely contains examples
	
	Usage:
		cryptonaut subscribe mempool --chain ethereum
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("subscribe called")
	},
}

var subscribeMempoolCmd = &cobra.Command{
	Use:   "mempool",
	Short: "Subscribe to mempool transactions",
	Long: `Subscribe to mempool transactions
	
	Usage:
		cryptonaut subscribe mempool --chain ethereum --to-address 0x0000000000000000000000000000000000000000 --ws-url wss://mainnet.infura.io/ws/v3/YOUR_PROJECT_ID
	`,
	RunE: runSubscribeMempoolCmd,
}

func init() {
	subscribeCmd.AddCommand(subscribeMempoolCmd)
	subscribeMempoolCmd.Flags().StringP("to-address", "t", "", "Filter transactions by to address")
	viper.BindPFlag("to-address", subscribeMempoolCmd.Flags().Lookup("to-address"))
	subscribeMempoolCmd.Flags().StringP("ws-url", "w", "", "Websocket URL")
	subscribeMempoolCmd.MarkFlagRequired("ws-url")
	viper.BindPFlag("ws-url", subscribeMempoolCmd.Flags().Lookup("ws-url"))
	rootCmd.AddCommand(subscribeCmd)

}

func runSubscribeMempoolCmd(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithCancel(cmd.Context())
	defer cancel()

	wsURL := viper.GetString("ws-url")
	toAddress := viper.GetString("to-address")

	// Create the client
	client, err := clients.NewEthereumClient(wsURL)
	if err != nil {
		return err
	}
	defer client.Close()

	// Create and start the subscription
	sub, err := subscription.NewMempoolSubscription(client, toAddress)
	if err != nil {
		return err
	}
	defer sub.Stop()

	if err := sub.Start(ctx); err != nil {
		return err
	}

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	return nil
}
