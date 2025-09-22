# NFT Marketplace Architecture

## Overview
The NFT marketplace allows creators to mint and trade tokens representing unique assets. It reuses the core token and contract systems, exposing a browser interface that wraps the `nft_marketplace` CLI.

## Key Modules
- `core/nft_marketplace.go` – in-memory catalogue and trade execution logic.
- `cli/nft_marketplace.go` – create listings, buy and sell NFTs.
- `smart-contracts/nft_minting.wasm` – sample contract template for minting.
- GUI under `GUI/nft_marketplace` – React frontend for user interaction.

## Workflow
1. **Minting** – creators deploy the `nft_minting.wasm` template via the contracts CLI.
2. **Listing** – `nft_marketplace` registers the token with price and owner information.
3. **Purchase** – buyers invoke the CLI or GUI to transfer ownership atomically.
4. **Catalogue update** – the marketplace updates its in-memory index and broadcasts events.

## Security Considerations
- Listings reference on-chain token IDs to prevent spoofing.
- Trades execute atomically to avoid partial transfers.
- The marketplace can enforce whitelists or royalties via smart contract hooks.

## CLI Integration
- `synnergy nft-marketplace` – manage listings and execute trades.

## Enterprise Diagnostics
- Marketplace operations rely on VM availability and ledger settlement; both are validated when `synnergy integration status` mines its diagnostic block and reports VM concurrency.
- Automation can gate NFT releases on a clean diagnostics report, aligning CLI and GUI expectations before creators schedule drops.
