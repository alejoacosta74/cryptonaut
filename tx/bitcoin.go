package tx

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/wire"
)

func DecodeBitcoinRawTx(rawTx string) (*wire.MsgTx, error) {
	rawTxBytes, err := hex.DecodeString(rawTx)
	if err != nil {
		return nil, fmt.Errorf("failed to decode transaction: %v", err)
	}

	var tx wire.MsgTx

	reader := bytes.NewReader(rawTxBytes)
	err = tx.Deserialize(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize transaction: %v", err)
	}

	return &tx, nil
}
