# Tokens

## Token Families
Synnergy supports multiple token standards covering payment, utility, and compliance use cases. Stage 80 aligns the runtime and documentation by pre-registering all token opcodes during CLI bootstrap, ensuring wallets and dashboards can calculate costs deterministically.【F:cmd/synnergy/main.go†L110-L143】 Synthron Coin acts as the canonical example with ed25519 signatures, nonce enforcement, and audit-ready ledger entries.【F:docs/Whitepaper_detailed/Synthron Coin.go†L1-L167】

The SYN500 utility token extends the suite with enterprise service entitlements. Controllers can grant windowed usage tiers, the CLI tracks consumption and telemetry, and audit logs capture every allocation for compliance review.【F:core/syn500.go†L1-L150】【F:cli/syn500.go†L1-L172】

## Lifecycle Operations
- **Issuance** – authority wallets sign mint requests; the VM validates signatures and increments supply under governance control.
- **Transfer** – token transfers run through the VM sandbox with gas metering and replay protection, generating immutable ledger records.
- **Redemption** – treasury functions burn or lock tokens as part of regulatory actions. Reversal workflows are available when governance approves compensating transactions.

## Integration Points
APIs and CLI commands expose balance queries, transfer operations, and compliance hooks. Because gas snapshots propagate through subscriptions, partners integrating with the token suite receive live fee updates without redeploying clients.【F:gas_table.go†L20-L152】

## Security Considerations
Every token action is signed and logged. Automated tests cover mint, transfer, and double-spend scenarios to ensure the VM, gas registry, and ledger stay in sync across upgrades.【F:virtual_machine_test.go†L1-L96】【F:core/virtual_machine_test.go†L1-L108】
