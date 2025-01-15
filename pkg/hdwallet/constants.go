package hdwallet

type DerivationIndex uint32

func (d DerivationIndex) ToUint32() uint32 {
	return uint32(d)
}

const (
	PurposePath          DerivationIndex = 44
	BitcoinCoinTypePath  DerivationIndex = 0
	EthereumCoinTypePath DerivationIndex = 60
	AccountPath          DerivationIndex = 0
	ChainPath            DerivationIndex = 0
)
