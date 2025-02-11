package ethereum

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
)

// MempoolSubscription represents an active mempool subscription
type MempoolSubscription struct {
	client       *EthereumClient
	toAddress    *common.Address // optional filter
	subscription *rpc.ClientSubscription
	txChan       chan *types.Transaction
	txHashChan   chan common.Hash
}

// NewMempoolSubscription creates a new mempool subscription
func NewMempoolSubscription(client *EthereumClient, toAddress string) (*MempoolSubscription, error) {
	var addr *common.Address
	if toAddress != "" {
		parsedAddr := common.HexToAddress(toAddress)
		addr = &parsedAddr
	}

	return &MempoolSubscription{
		client:     client,
		toAddress:  addr,
		txChan:     make(chan *types.Transaction),
		txHashChan: make(chan common.Hash),
	}, nil
}

// Start begins the subscription
func (s *MempoolSubscription) Start(ctx context.Context) error {
	gethClient := s.client.GetGethClient()
	// sub, err := gethClient.SubscribeFullPendingTransactions(ctx, s.txChan)
	sub, err := gethClient.SubscribePendingTransactions(ctx, s.txHashChan)
	if err != nil {
		return fmt.Errorf("failed to subscribe to pending transactions: %w", err)
	}
	s.subscription = sub

	go s.handleTransactions(ctx)
	return nil
}

// Stop cleanly stops the subscription
func (s *MempoolSubscription) Stop() {
	if s.subscription != nil {
		s.subscription.Unsubscribe()
	}
}

func (s *MempoolSubscription) handleTransactions(ctx context.Context) {
	for {
		select {
		case txHash := <-s.txHashChan:
			ethClient := s.client.GetEthClient()
			tx, isPending, err := ethClient.TransactionByHash(ctx, txHash)
			if err != nil {
				log.Printf("Failed to fetch transaction: %v", err)
				continue
			}
			// filter by to address
			if s.toAddress != nil && tx.To() != nil && *s.toAddress != *tx.To() {
				continue
			}
			processTransaction(tx, isPending)
		case err := <-s.subscription.Err():
			log.Printf("Subscription error: %v", err)
			return
		case <-ctx.Done():
			log.Println("Context done")
			return
		}
	}
}

func processTransaction(tx *types.Transaction, isPending bool) {
	fmt.Printf("Transaction found (pending: %t):\n", isPending)
	fmt.Printf("  Hash: %s\n", tx.Hash().Hex())
	fmt.Printf("  To: %s\n", tx.To())
	fmt.Printf("  Value: %f ETH\n", weiToEther(tx.Value()))
	fmt.Printf("  Gas Price: %f Gwei\n", weiToGwei(tx.GasPrice()))
	// fmt.Printf("  Input Data: %s\n", hex.EncodeToString(tx.Data()))
	msg, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
	if err != nil {
		log.Printf("Failed to get sender address for transaction %s: %v", tx.Hash().Hex(), err)
	}
	fmt.Printf("  From: %s\n", msg.Hex())
	fmt.Println()

}

// Helper function to convert Wei to Ether
func weiToEther(wei *big.Int) *big.Float {
	return new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(params.Ether))
}

// Helper function to convert Wei to Gwei if needed
func weiToGwei(wei *big.Int) *big.Float {
	return new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(params.GWei))
}
