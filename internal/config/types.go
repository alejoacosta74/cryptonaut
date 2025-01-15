package config

// Algorithm represents supported signing algorithms
type Algorithm string

const (
	AlgoSchnorr Algorithm = "schnorr"
	AlgoECDSA   Algorithm = "ecdsa"
	AlgoBLS     Algorithm = "bls"
)

// Network represents supported blockchain networks
type Network string

const (
	NetworkMainnet Network = "mainnet"
	NetworkTestnet Network = "testnet"
)
