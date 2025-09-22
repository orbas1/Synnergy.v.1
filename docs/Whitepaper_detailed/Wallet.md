# Wallet

## Architecture
Synnergy wallets are permissioned clients that sign every operation with ed25519 keys. Stage 80 aligns wallet services with the central gas registry so fee estimates and signature workflows reflect live pricing without manual refresh. During CLI bootstrap the wallet module initialises alongside the VM, secrets manager, and content node stack, guaranteeing that signatures, gas pricing, and storage attestations share the same trust boundary.【F:cmd/synnergy/main.go†L110-L143】

## Security
- **Key Management** – secrets are loaded through the internal secrets manager; hardware modules can be integrated using the same interface.
- **Signature Workflow** – every mint, transfer, and reversal is signed, including nonce metadata to prevent replay.【F:docs/Whitepaper_detailed/Synthron Coin.go†L1-L167】
- **Telemetry** – wallets subscribe to gas updates to present accurate fee prompts to users.【F:gas_table.go†L20-L152】

## User Experience
Wallet clients expose:
- Transaction builders that integrate with the opcode catalogue, providing human-readable descriptions of every action.
- Reversal status dashboards that reflect governance decisions in real time.【F:docs/Whitepaper_detailed/Reversing and cancelling transactions.md†L1-L19】
- Ledger explorers that combine balances, nonce history, and gas version metadata.
- Stage 73 and Stage 80 CLI flows require wallet provenance for governance-sensitive operations, ensuring `syn3700`, `syn3800`, `syn3900`, and utility token entitlements are all executed with authenticated credentials.【F:cli/syn3700_token.go†L15-L229】【F:cli/syn3800.go†L1-L203】【F:cli/syn3900.go†L1-L140】
- The shared Stage 73 state file (`--stage73-state`) can be supplied alongside wallet arguments so automated workflows operate on the same persisted snapshots that controllers and auditors review.【F:cli/stage73_state.go†L1-L142】

## Testing
End-to-end tests exercise wallet scenarios through the Synthron Coin reference implementation and VM hot reload flows, ensuring signatures, gas metering, and ledger updates stay consistent across upgrades.【F:virtual_machine_test.go†L1-L96】【F:docs/Whitepaper_detailed/Synthron Coin_test.go†L1-L49】
