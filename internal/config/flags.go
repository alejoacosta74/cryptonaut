package config

// Flag names used throughout the application
const (
	// Common flags
	FlagLogLevel   = "log-level"
	FlagConfigFile = "config"

	// Key-related flags
	FlagPrivateKey       = "key"
	FlagPrivateKeyFormat = "format"
	FlagPublicKey        = "pubkey"

	// Operation flags
	FlagMessage   = "message"
	FlagSignature = "signature"
	FlagAlgorithm = "algo"

	// Network flags
	FlagChain    = "chain"
	FlagTestnet  = "testnet"
	FlagNetwork  = "network"
	FlagEndpoint = "endpoint"

	// BIP44 derivation flags
	FlagMnemonic = "mnemonic"
	FlagIndex    = "index"
)

// Flag default values
const (
	DefaultLogLevel = "info"
	DefaultNetwork  = "mainnet"
)
