package ethereum

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

func DecodeEthereumRawTx(rawTx string) (*types.Transaction, error) {
	rawTxBytes, err := hex.DecodeString(rawTx)
	if err != nil {
		return nil, fmt.Errorf("failed to decode transaction: %v", err)
	}

	var tx types.Transaction

	err = rlp.DecodeBytes(rawTxBytes, &tx)
	if err != nil {
		return nil, fmt.Errorf("failed to decode transaction: %v", err)
	}

	return &tx, nil
}
