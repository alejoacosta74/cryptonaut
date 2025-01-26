package clients

import (
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type EthereumClient struct {
	rpcClient  *rpc.Client
	ethClient  *ethclient.Client
	gethClient *gethclient.Client
	wsURL      string
}

// NewEthereumClient creates a new instance of EthereumClient
func NewEthereumClient(wsURL string) (*EthereumClient, error) {
	rpcClient, err := rpc.Dial(wsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum node: %w", err)
	}

	return &EthereumClient{
		rpcClient:  rpcClient,
		ethClient:  ethclient.NewClient(rpcClient),
		gethClient: gethclient.New(rpcClient),
		wsURL:      wsURL,
	}, nil
}

// Close cleanly shuts down the client connections
func (c *EthereumClient) Close() {
	if c.rpcClient != nil {
		c.rpcClient.Close()
	}
}

func (c *EthereumClient) GetGethClient() *gethclient.Client {
	return c.gethClient
}

func (c *EthereumClient) GetEthClient() *ethclient.Client {
	return c.ethClient
}
