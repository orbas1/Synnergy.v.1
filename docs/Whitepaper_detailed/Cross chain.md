# Cross-Chain Transactions

Stage 42 expands Synnergy's cross-chain capabilities with a dedicated CLI for executing asset transfers between networks.  The `cross_tx` module exposes two primary operations:

- **lockmint** – locks native assets on the source chain and mints wrapped tokens on the destination.
- **burnrelease** – burns wrapped tokens and releases the original assets back to a native chain.

Both commands report deterministic gas usage and optionally emit JSON so dashboards and automation tools can coordinate transfers through the function web.  Transfers are tracked by the `CrossChainTxManager`, which records bridge identifiers, participant addresses and completion status for auditability.

These primitives pair with bridge registries and connection managers to provide fault-tolerant, enterprise-grade interoperability with external chains.
