# GUIs

Synnergy ships lightweight web interfaces that interact with the CLI.  The
wallet GUI (Stage 31) generates encrypted wallets and queries balances.  Stage
32 introduces an Explorer GUI that issues read-only CLI RPC calls to display
blockchain state such as chain height and block details.  These tools provide a
foundation for richer dashboards and can be extended to visualise transactions,
addresses and node health.

Stage 33 adds an AI Marketplace GUI that deploys AI-enhanced contracts through the CLI. Users supply a WebAssembly module, model hash, manifest and gas limit, and the interface returns the on-chain contract address. This marketplace demonstrates how advanced contract workflows can be wrapped in thin, CLI-driven front ends.
Stage 34 introduces a Smart-Contract Marketplace GUI that leverages the new
`marketplace` CLI commands to deploy generic WebAssembly contracts and trade
their ownership. It serves as a reference for integrating contract workflows
into web applications.
Stage 35 introduces a Storage Marketplace GUI enabling users to list and lease storage through the CLI while inheriting the runtime's sandboxing and gas accounting.
Stage 36 debuts an NFT Marketplace GUI for minting and trading NFTs via CLI commands, showcasing opcode-priced asset workflows.
Stage 37 adds a DAO Explorer GUI that manages decentralised autonomous organisations by spawning the `synnergy dao` commands to create, join, leave and inspect governance groups.
Stage 38 introduces a Token Creation Tool GUI that spawns the CLI to create new token contracts with deterministic gas costs.
Stage 39 adds a DEX Screener GUI that queries liquidity pools through the CLI so dashboards can monitor reserves and pricing in real time.

Stage 40 introduces Administrative Dashboards for authority node indexing and cross-chain management, allowing operators to supervise nodes and bridges via CLI-backed views.
