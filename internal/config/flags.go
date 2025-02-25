package config

// Flag names used throughout the application
const (
	// Common flags
	FlagLogLevel   = "log-level"
	FlagConfigFile = "config"

	// Key-related flags
	FlagPrivateKey          = "private-key"
	FlagPrivateKeyFormat    = "format"
	FlagCosmosAddressPrefix = "cosmos-address-prefix"
	FlagPublicKey           = "public-key"
	FlagPubKeyCompressed    = "compressed"

	// Operation flags
	FlagMessage   = "message"
	FlagSignature = "signature"
	FlagAlgorithm = "algo"

	// Network flags
	FlagChain    = "chain"
	FlagTestnet  = "testnet"
	FlagNetwork  = "network"
	FlagEndpoint = "endpoint"
	FlagWsUrl    = "ws-url"

	// BIP44 derivation flags
	FlagMnemonic = "mnemonic"
	FlagIndex    = "index"

	// Bitcoin flags
	FlagBitcoinFormat = "bitcoin-format"

	// ECDSA flags
	FlagSignatureR = "r"
	FlagSignatureS = "s"

	// ZK flags
	FlagCircuit   = "circuit"
	FlagProof     = "proof"
	FlagVk        = "vk" // verification key
	FlagBirthYear = "birth-year"
)

// Flag default values
const (
	DefaultLogLevel = "info"
	DefaultNetwork  = "mainnet"
)
