# NFT Marketplace Architecture

The NFT marketplace provides a thin coordination layer between the CLI, core
runtime and the JavaScript GUI. It maintains an in-memory catalogue of NFTs and
relies on the shared gas table for fee estimation.

* **Virtual machine integration** – NFT minting uses the same VM sandboxing
  mechanisms as smart contracts ensuring deterministic execution.
* **Consensus and wallet integration** – operations consume gas via
  `MintNFT`, `ListNFT` and `BuyNFT` opcodes so wallets and consensus nodes can
  account for costs.
* **Fault tolerance** – the marketplace is concurrency-safe and can be wrapped
  by higher level services for persistence or sharding in enterprise
deployments.

The accompanying GUI invokes the CLI which in turn calls the core module,
allowing existing authentication, node and authority policies to apply without
additional plumbing.
