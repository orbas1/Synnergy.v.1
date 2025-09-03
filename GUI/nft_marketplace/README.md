# NFT Marketplace GUI

This reference interface demonstrates how NFTs can be minted and traded
through the Synnergy CLI.

## Usage

```bash
# Mint an NFT
node -e "import('./src/main.ts').then(m => m.mintNFT('id1','alice','meta','100'))"

# Buy an NFT
node -e "import('./src/main.ts').then(m => m.buyNFT('id1','bob'))"
```

The implementation spawns the local `synnergy` binary and therefore inherits all
CLI configuration such as network endpoints and gas table settings.
