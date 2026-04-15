# Polymarket Key Generator

Polymarket does not always expose their CLOB (Central Limit Order Book) API keys via their web interface for non-institutional wallets. 

To generate your `POLY_API_KEY`, `POLY_API_SECRET`, and `POLY_API_PASSPHRASE` required by this plugin, you must sign a cryptographic message using your wallet's private key to prove ownership to the Polymarket servers.

This standalone script automates that process.

## Prerequisites

1. You must have Go installed.
2. You need your Coinbase Wallet's **Public Address** (0x...)
3. You need your Coinbase Wallet's **Polygon Private Key** (0x...)
   - *In Coinbase Wallet: Settings > Security/Developer > Show Private Key*

## Usage

Run the script locally on your machine, passing your address and private key as arguments:

```bash
go run keygen.go <PUBLIC_ADDRESS> <POLYGON_PRIVATE_KEY>
```

**Example:**
```bash
go run keygen.go 0x1234567890abcdef1234567890abcdef12345678 0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890
```

The script will ping the Polymarket API and print out your three API keys. You can then copy and paste them into the Allele Dashboard for the `allele-exchange-polymarket` plugin.

*Security Note: This script runs entirely locally. Your private key is only used to generate a mathematical signature and is never transmitted over the internet.*
