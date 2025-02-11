# Cryptonaut CLI 🚀

Your Swiss Army Knife for Blockchain and Cryptocurrency Operations

## Overview

Cryptonaut is a powerful command-line tool that provides a comprehensive suite of cryptocurrency and blockchain utilities. Whether you're working with Bitcoin, Ethereum, Cosmos, or exploring cryptographic operations, Cryptonaut has got you covered.

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

#### Cosmos

- Generate a new private key:

```bash
cryptonaut cosmos generate 
Cosmos private key: dc31eed917590637849c14a1c032372201b823e622b86491c913a2133ffb7a3b0d797b9c272a6a962b2d107e95c8347197bca8c87bdc50f081db934d56482bfb
```

- Get the public key from a private key:

```bash
cryptonaut cosmos pubkey --private-key dc31eed917590637849c14a1c032372201b823e622b86491c913a2133ffb7a3b0d797b9c272a6a962b2d107e95c8347197bca8c87bdc50f081db934d56482bfb
Cosmos public key: 0D797B9C272A6A962B2D107E95C8347197BCA8C87BDC50F081DB934D56482BFB
```

- Get the address from a private key:

```bash
cryptonaut cosmos address --private-key dc31eed917590637849c14a1c032372201b823e622b86491c913a2133ffb7a3b0d797b9c272a6a962b2d107e95c8347197bca8c87bdc50f081db934d56482bfb
Cosmos address: cosmos102qqme5y2unezjlg8xghzmas0jdyfdely479z
```

#### Bitcoin

- Generate a new private key:

```bash
cryptonaut bitcoin generate 
Private key: 5887c2df0c75bc44dd1e33f3f45c08f39a0970a8fda69f1aa241831ee983dc71
```

- Get the public key from a private key:

```bash
cryptonaut bitcoin pubkey --private-key 5887c2df0c75bc44dd1e33f3f45c08f39a0970a8fda69f1aa241831ee983dc71
Private key: 5887c2df0c75bc44dd1e33f3f45c08f39a0970a8fda69f1aa241831ee983dc71
Public key: 03709cea5a5fb840d2d5c85b76db27f5f5969693140ac5e9d904914fc891129b7f
```

- Get the address from a private key:

```bash
cryptonaut bitcoin address --private-key 5887c2df0c75bc44dd1e33f3f45c08f39a0970a8fda69f1aa241831ee983dc71
Address: 1P3Ykb3ZZnEMKAhb6NhW4Lex7h24qdfwuH
```
#### Ethereum

- Generate a new private key:

```bash
cryptonaut ethereum generate
Private Key: 45b4314b4f5964ed2b6030909da0fdd7bcf8e653dfef438233ea45b1b59f0d0f
```

- Get the public key from a private key

```bash
cryptonaut ethereum pubkey --private-key 45b4314b4f5964ed2b6030909da0fdd7bcf8e653dfef438233ea45b1b59f0d0f
Public Key: 907e5d2b3c5e4cdef69a8c547aa1280df55ed82903475ca10f035c1c4bd27bd2
```

- Get the address from a private key

```bash
cryptonaut ethereum address --private-key 45b4314b4f5964ed2b6030909da0fdd7bcf8e653dfef438233ea45b1b59f0d0f
Address: 0x6e91d895Cd7c010fbA616260FeCe1FC1d4AA4a85
```

### Managing keys and digital Signatures

#### ECDSA signatures (P-256 i.e. curve secp256r1, equation y² = x³ - 3x + b)

```bash
# Generate a ECDSA private key
cryptonaut ecdsa generate
Private key: 3077020101042045471f8f388f79438402107d6904db836f81445a840464561b6eab9918b1646ca00a06082a8648ce3d030107a144034200040f343af9c39b929a542866089cac57d167f85900ea5b891b4aeb5a6dc624c35170053107601dc60409c9e24b7ecc2e0624cf069ec2e94c791b07123abf07e1f1

# Generate a ECDSA public key from a private key
cryptonaut ecdsa pubkey --private-key 3077020101042045471f8f388f79438402107d6904db836f81445a840464561b6eab9918b1646ca00a06082a8648ce3d030107a144034200040f343af9c39b929a542866089cac57d167f85900ea5b891b4aeb5a6dc624c35170053107601dc60409c9e24b7ecc2e0624cf069ec2e94c791b07123abf07e1f1
Public key: 040f343af9c39b929a542866089cac57d167f85900ea5b891b4aeb5a6dc624c35170053107601dc60409c9e24b7ecc2e0624cf069ec2e94c791b07123abf07e1f1

# Sign a message
cryptonaut ecdsa sign "hello world" --private-key 3077020101042045471f8f388f79438402107d6904db836f81445a840464561b6eab9918b1646ca00a06082a8648ce3d030107a144034200040f343af9c39b929a542866089cac57d167f85900ea5b891b4aeb5a6dc624c35170053107601dc60409c9e24b7ecc2e0624cf069ec2e94c791b07123abf07e1f1

Signature: 
r=18f137a519e7c505bdf8aefbbe67045b5138c1f123e79c4578df719acb59ea63
s=d0caa13a4e35c9fee6a471f0b830e30d236b6799760fe99dbb0666a125554b10

# Verify a signature
cryptonaut ecdsa verify "hello world" --r=18f137a519e7c505bdf8aefbbe67045b5138c1f123e79c4578df719acb59ea63 --s=d0caa13a4e35c9fee6a471f0b830e30d236b6799760fe99dbb0666a125554b10 --public-key 040f343af9c39b929a542866089cac57d167f85900ea5b891b4aeb5a6dc624c35170053107601dc60409c9e24b7ecc2e0624cf069ec2e94c791b07123abf07e1f1
Signature is valid: true
```

#### Schnorr signatures (secp256k1)
  
```bash
# Generate a private key
cryptonaut schnorr generate
Private key: c201901e446e3707a6bd3f4c1de85939aee1fd68fc2926690a33d4fc88063d78

# Get the public key from a private key
cryptonaut schnorr pubkey --private-key c201901e446e3707a6bd3f4c1de85939aee1fd68fc2926690a33d4fc88063d78
Public key: 03df7d7f349f4f63e79c3aea2227c2b96ef5a9aa51960a7962a3e0d95571fcc84a

# Get the public key from a private key in un compressed format
cryptonaut schnorr pubkey --private-key c201901e446e3707a6bd3f4c1de85939aee1fd68fc2926690a33d4fc88063d78 --compressed=false
Public key: 04df7d7f349f4f63e79c3aea2227c2b96ef5a9aa51960a7962a3e0d95571fcc84a4d2e21eb89130a9c96952d191505620351a369cb47462e3d7e09e79fa49d07d7

# Sign a message
cryptonaut schnorr sign "hello world" --private-key '479408efb759a4fcf8f482a45ecc8e6185fbe24ff4ee5deca8d390e4bcddd947' 
Signature: 047894e7ca77a4f5597136ac015396ae6098a258c3f918630acd06b0e444485e6630d7179f43f289e4ad5f05f6f424ab15c99b4d11f33c4ab38a664ddef4825a

# Verify a signature
cryptonaut schnorr verify "hello world" --public-key 03d43aa64ab048f935da807d95f8efc7e8f3425c3b4da8f7cffc2721b14a1dd666 --signature 047894e7ca77a4f5597136ac015396ae6098a258c3f918630acd06b0e444485e6630d7179f43f289e4ad5f05f6f424ab15c99b4d11f33c4ab38a664ddef4825a 
Signature is valid:  true
```

#### BLS signatures (BLS12-381)

```bash
# Generate a private key
cryptonaut bls generate
Private key: 40b04d19c23c45469e8b875a7c7cc00368b41555470854a403e0ea059135d7e9 

# Get the public key from the BLS private key
cryptonaut bls pubkey --private-key 40b04d19c23c45469e8b875a7c7cc00368b41555470854a403e0ea059135d7e9
Public key: b0b90706e41978b0d5400b51306b1e54c78e4d8b569baeb4939159635a9576d8d413caec63bfd664f256df3b77dab69f

# Sign a message
cryptonaut bls sign "hello world" --private-key 40b04d19c23c45469e8b875a7c7cc00368b41555470854a403e0ea059135d7e9
Signature: b88c8beef218a88179dd6bfda698d934cd011326265eefa106547ce4fd8d310d4b4175c7d14ac513b6d88dcdd95ebeb017785d3b45d5eeff0304074d2642760b2dc7743f4aee381b9db089f96e2232da2004605ecbfd2e62e786a7e2d196b093

# Verify a signature
cryptonaut bls verify "hello world" --public-key b0b90706e41978b0d5400b51306b1e54c78e4d8b569baeb4939159635a9576d8d413caec63bfd664f256df3b77dab69f --signature b88c8beef218a88179dd6bfda698d934cd011326265eefa106547ce4fd8d310d4b4175c7d14ac513b6d88dcdd95ebeb017785d3b45d5eeff0304074d2642760b2dc7743f4aee381b9db089f96e2232da2004605ecbfd2e62e786a7e2d196b093
Signature is valid: true
```

### HD Wallet Operations

#### Generate mnemonic:
```bash
cryptonaut hd mnemonic
Mnemonic: genius unique bicycle wood bullet cross economy move bulb canvas nurse extend flight urge account island please people angry length snap foil brick congress
```

#### Derive keys:
```bash
# For Bitcoin
cryptonaut hd bitcoin --mnemonic 'legend rude glance must update smooth fever alone clarify stool harbor dutch swarm casual brisk odor capital good strong ensure wreck hybrid chalk ketchup' --index 0            
Private Key: f5b58ecb663dcc8e648876d335804dfe7de8542467746a35b64d9aa7ab260b41
Public Key: 03f93a8a9f7934eb5f60e3dee14d97aefa37d20b51df387f0faf7069be490d1bd1
Address: 1MjFFWJC6L3qzXhDQgNmttad77Qcn8mVyb

# For Ethereum
cryptonaut hd ethereum --mnemonic 'legend rude glance must update smooth fever alone clarify stool harbor dutch swarm casual brisk odor capital good strong ensure wreck hybrid chalk ketchup' --index 0
Address: 0x6099f0f046D843d6AD6a7daeC35c55b1D92A8cC8
```

### Transaction Management

#### Decode raw transactions:

- Ethereum
  
```json
cryptonaut ethereum tx decode f86b01843b9aca00825208941234567890123456789012345678901234567890880de0b6b3a76400008025a0b40bc16dbe93b2fd8698af2cbb2cd10ae64e15a1922d842153cf09fc1f26033da0429b5caf480e7840843f9451bd8f5cbb14f6cebb081dabfe6663c88dbfa56f8b

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
- Bitcoin

```json
cryptonaut bitcoin tx decode 010000000134129078563412907856341290785634129078563412907856341290785634120000000000ffffffff0100e1f505000000001976a914bade2cc53d518a756148ca179894efba4089a44888ac00000000
{
    "hash": "a378a99a0a32f789cea579179db4fe697375baa0436adcc053724a07bb254f4e",
    "version": 1,
    "locktime": 0,
    "size": 85,
    "inputs": [
        {
            "txid": "1234567890123456789012345678901234567890123456789012345678901234",
            "vout": 0,
            "scriptSig": "",
            "sequence": 4294967295
        }
    ],
    "outputs": [
        {
            "value": 100000000,
            "scriptPubKey": "76a914bade2cc53d518a756148ca179894efba4089a44888ac"
        }
    ]
}
```

### Subscription

Subscribe to mempool transactions:

- Ethereum

```bash
cryptonaut ethereum tx mempool --to-address 0x0000000000000000000000000000000000000000 --ws-url wss://mainnet.infura.io/ws/v3/YOUR_PROJECT_ID
```

## Roadmap 🗺️

### Coming Soon
- [ ] Ethereum node interaction (balance checks, transaction broadcasting)
- [ ] Bitcoin node interaction (balance checks, transaction broadcasting)
- [ ] Vanity address generation
- [ ] Smart contract deployment and interaction
- [ ] ERC-20 and ERC-721 token operations
- [ ] Proof-of-Work simulation
- [ ] Secure storage for crypto artifacts

### Planned Features
- [ ] Support for additional blockchain networks
- [ ] Advanced transaction building features
- [ ] Multi-signature wallet support
- [ ] Support for additional signature schemes

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

[MIT License](LICENSE)