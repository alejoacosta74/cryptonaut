# Cryptonaut CLI 🚀

Your Swiss Army Knife for Blockchain and Cryptocurrency Operations

## Overview

Cryptonaut is a powerful command-line tool that provides a comprehensive suite of cryptocurrency and blockchain utilities. Whether you're working with Bitcoin, Ethereum, or exploring cryptographic operations, Cryptonaut has got you covered.

## Features

- 🔑 Key Management (Bitcoin & Ethereum)
- 🔐 Multiple Signature Schemes (Schnorr, ECDSA*, BLS*)
- 👛 HD Wallet Support (BIP44)
- 📝 Transaction Management
- 🔗 Blockchain Node Interaction*
- 🎯 Vanity Address Generation*
- 🌳 Merkle Tree Operations*
- ⛏️ PoW Simulation*

(*) Coming soon

## Installation

```bash
go install github.com/yourusername/cryptonaut@latest
```

## Available Commands

### Key and Address Management

Generate a new private key:
```bash
# In hex format
cryptonaut generate key                                                                               
Private key: 18d4998b89bc1cf229c9b729ed39106c23c3dbf0637f89b2384fe3836c4e3247

# In WIF format (for Bitcoin)
cryptonaut generate key --format wif
Private key: Ky4LV32SbdWQv8uDsvtguVFEbxSBSRG1Rk2mz6DDWzb4uTgU9tJg

# For testnet
cryptonaut generate key --format wif --testnet
Private key: cRYDPznKNZPv2JVcYpWQXc6Zfkn68N22S259NqDPNJ4Pe2Nzdm7C
```

Derive addresses:
```bash
# Bitcoin address from hex private key
cryptonaut derive address --key 6abd31bf5fe56e1aa5a49b8430a2bcaa276b4cd352b3d7072e89bb9a8a204cc1 --chain bitcoin                                                                                           
Address: 1EjXB4qohumD9Tbk4NMekqCkLz1baWChNW

# Bitcoin address from WIF
cryptonaut derive address --key KzoCR5BTboQXqG9ah8HiHtigrK2DkrpgouYg94m4ZRWiCVEybGoy --chain bitcoin 
Address: 1EjXB4qohumD9Tbk4NMekqCkLz1baWChNW
```

### Digital Signatures

Currently supporting Schnorr signatures:
```bash
# Sign a message
cryptonaut sign schnorr --message "hello world" --key '479408efb759a4fcf8f482a45ecc8e6185fbe24ff4ee5deca8d390e4bcddd947' 
Signature: 047894e7ca77a4f5597136ac015396ae6098a258c3f918630acd06b0e444485e6630d7179f43f289e4ad5f05f6f424ab15c99b4d11f33c4ab38a664ddef4825a

# Verify a signature
cryptonaut verify schnorr --message "hello world" --pubkey 03d43aa64ab048f935da807d95f8efc7e8f3425c3b4da8f7cffc2721b14a1dd666 --signature 047894e7ca77a4f5597136ac015396ae6098a258c3f918630acd06b0e444485e6630d7179f43f289e4ad5f05f6f424ab15c99b4d11f33c4ab38a664ddef4825a 
```

### HD Wallet Operations

Generate mnemonic:
```bash
cryptonaut generate mnemonic
Mnemonic: genius unique bicycle wood bullet cross economy move bulb canvas nurse extend flight urge account island please people angry length snap foil brick congress
```

Derive keys:
```bash
# For Bitcoin
cryptonaut bip44 bitcoin --mnemonic 'legend rude glance must update smooth fever alone clarify stool harbor dutch swarm casual brisk odor capital good strong ensure wreck hybrid chalk ketchup' --index 0            
Private Key: f5b58ecb663dcc8e648876d335804dfe7de8542467746a35b64d9aa7ab260b41
Public Key: 03f93a8a9f7934eb5f60e3dee14d97aefa37d20b51df387f0faf7069be490d1bd1
Address: 1MjFFWJC6L3qzXhDQgNmttad77Qcn8mVyb

# For Ethereum
cryptonaut bip44 ethereum --mnemonic 'legend rude glance must update smooth fever alone clarify stool harbor dutch swarm casual brisk odor capital good strong ensure wreck hybrid chalk ketchup' --index 0
Address: 0x6099f0f046D843d6AD6a7daeC35c55b1D92A8cC8
```

### Transaction Management

Decode raw transactions:
```json
cryptonaut tx decode ethereum f86b01843b9aca00825208941234567890123456789012345678901234567890880de0b6b3a76400008025a0b40bc16dbe93b2fd8698af2cbb2cd10ae64e15a1922d842153cf09fc1f26033da0429b5caf480e7840843f9451bd8f5cbb14f6cebb081dabfe6663c88dbfa56f8b

{
    "hash": "0xc7b6e5e7a83c44651cc0a4ceb33eaeafbb84c7b9f25690443ed3d669a29e0a72",
    "nonce": 1,
    "gasPrice": "1000000000",
    "gas": 21000,
    "to": "0x1234567890123456789012345678901234567890",
    "value": "1000000000000000000",
    "data": "",
    "chainId": "1",
    "type": 0
}
```

## Roadmap 🗺️

### Coming Soon
- [ ] ECDSA signature scheme support
- [ ] BLS signature scheme support
- [ ] Ethereum node interaction (balance checks, transaction broadcasting)
- [ ] Bitcoin node interaction (balance checks, transaction broadcasting)
- [ ] Vanity address generation
- [ ] Smart contract deployment and interaction
- [ ] ERC-20 and ERC-721 token operations
- [ ] Merkle tree utilities
- [ ] Proof-of-Work simulation
- [ ] Secure storage for crypto artifacts

### Planned Features
- [ ] Support for additional blockchain networks
- [ ] Advanced transaction building features
- [ ] Multi-signature wallet support
- [ ] Integration with popular blockchain APIs
- [ ] Support for additional signature schemes

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

[MIT License](LICENSE)