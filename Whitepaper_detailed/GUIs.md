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
