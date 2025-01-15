## Managing keys and addresses

### Generating a private key

- In hex format

```bash
cryptonaut generate key                                                                               
Private key: 18d4998b89bc1cf229c9b729ed39106c23c3dbf0637f89b2384fe3836c4e3247
```

- In WIF format (for bitcoin)

```bash
cryptonaut generate key --format wif
Private key: Ky4LV32SbdWQv8uDsvtguVFEbxSBSRG1Rk2mz6DDWzb4uTgU9tJg

# specifying the network
cryptonaut generate key --format wif --testnet
Private key: cRYDPznKNZPv2JVcYpWQXc6Zfkn68N22S259NqDPNJ4Pe2Nzdm7C
```

### Generating an address

- For bitcoin

```bash

# from hex
cryptonaut derive address --key 6abd31bf5fe56e1aa5a49b8430a2bcaa276b4cd352b3d7072e89bb9a8a204cc1 --chain bitcoin                                                                                           
Address: 1EjXB4qohumD9Tbk4NMekqCkLz1baWChNW

# convert to wif
cryptonaut convert 6abd31bf5fe56e1aa5a49b8430a2bcaa276b4cd352b3d7072e89bb9a8a204cc1                                                    
Private key: KzoCR5BTboQXqG9ah8HiHtigrK2DkrpgouYg94m4ZRWiCVEybGoy

# from wif
cryptonaut derive address --key  KzoCR5BTboQXqG9ah8HiHtigrK2DkrpgouYg94m4ZRWiCVEybGoy --chain bitcoin 
Address: 1EjXB4qohumD9Tbk4NMekqCkLz1baWChNW
```

## Schnorr signatures

- Signing a message

```bash
cryptonaut sign schnorr --message "hello world" --key '479408efb759a4fcf8f482a45ecc8e6185fbe24ff4ee5deca8d390e4bcddd947' 
Signature: 047894e7ca77a4f5597136ac015396ae6098a258c3f918630acd06b0e444485e6630d7179f43f289e4ad5f05f6f424ab15c99b4d11f33c4ab38a664ddef4825a

# Get public key from private key
cryptonaut derive pubkey --key '479408efb759a4fcf8f482a45ecc8e6185fbe24ff4ee5deca8d390e4bcddd947'                       
Private key: 479408efb759a4fcf8f482a45ecc8e6185fbe24ff4ee5deca8d390e4bcddd947
Public key: 03d43aa64ab048f935da807d95f8efc7e8f3425c3b4da8f7cffc2721b14a1dd666                  

# verify signature
cryptonaut verify schnorr --message "hello world" --pubkey 03d43aa64ab048f935da807d95f8efc7e8f3425c3b4da8f7cffc2721b14a1dd666 --signature 047894e7ca77a4f5597136ac015396ae6098a258c3f918630acd06b0e444485e6630d7179f43f289e4ad5f05f6f424ab15c99b4d11f33c4ab38a664ddef4825a 
Signature is valid: true
```

## HD Wallet derivation

### Generating a mnemonic phrase

```bash
cryptonaut generate mnemonic
Mnemonic: genius unique bicycle wood bullet cross economy move bulb canvas nurse extend flight urge account island please people angry length snap foil brick congress
```

### Deriving keys from a mnemonic phrase

- For bitcoin

```bash
cryptonaut bip44 bitcoin --mnemonic 'legend rude glance must update smooth fever alone clarify stool harbor dutch swarm casual brisk odor capital good strong ensure wreck hybrid chalk ketchup' --index 0            
Private Key: f5b58ecb663dcc8e648876d335804dfe7de8542467746a35b64d9aa7ab260b41
Public Key: 03f93a8a9f7934eb5f60e3dee14d97aefa37d20b51df387f0faf7069be490d1bd1
Address: 1MjFFWJC6L3qzXhDQgNmttad77Qcn8mVyb
```

- For ethereum

```bash
cryptonaut bip44 ethereum --mnemonic 'legend rude glance must update smooth fever alone clarify stool harbor dutch swarm casual brisk odor capital good strong ensure wreck hybrid chalk ketchup' --index 0
Private Key: 964293bf0be5bc0935baa371734574db848b5c7e071c6102aeb10014cec584a6
Public Key: 0e1609df17c4091cdaeb1476ca8123f2d8badcf661fa80c1731806aa6eeea5cd
Address: 0x6099f0f046D843d6AD6a7daeC35c55b1D92A8cC8
```
