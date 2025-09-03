# Cross-Chain Architecture

Cross-chain modules enable interoperability with external networks through bridges, agnostic protocols, and transaction relays.

**Key Modules**
- cross_chain.go
- cross_chain_bridge.go
- cross_chain_connection.go
- cross_chain_contracts.go
- cross_chain_transactions.go
- cross_chain_agnostic_protocols.go
- cross_consensus_scaling_networks.go

**Related CLI Files**
- cli/cross_chain.go
- cli/cross_chain_bridge.go
- cli/cross_chain_connection.go
- cli/cross_chain_contracts.go
- cli/cross_chain_transactions.go
- cli/cross_chain_agnostic_protocols.go
- cli/cross_consensus_scaling_networks.go

These components coordinate communication and asset transfers across multiple blockchains.

Stage 8 hardens these modules for production use.  Each manager is concurrency
safe, exposes deterministic gas‑priced opcodes and is accessible through the
`synnergy` CLI.  Registries and bridge transfers persist in memory but are
designed to be swapped with database backends for fault‑tolerant deployments.
