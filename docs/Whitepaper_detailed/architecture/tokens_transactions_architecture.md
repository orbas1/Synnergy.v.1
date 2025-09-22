# Tokens and Transactions Architecture

## Overview
This architecture governs token standards, contract execution and the lifecycle of transactions on the Synnergy Network. It ensures interoperability among numerous SYN token variants and maintains deterministic gas accounting.

## Key Modules
- Token implementations under `internal/tokens/` such as `syn20`, `syn200`, `syn300` and others.
- `contracts.go` and `contracts_opcodes.go` – manage contract deployment and invocation.
- `cross_chain_transactions.go` – relays transactions across networks.
- `gas_table.go` – records gas costs for operations.

## Workflow
1. **Token creation** – CLI commands like `synnergy syn20` mint standard tokens.
2. **Transaction assembly** – clients construct transactions referencing token contracts and desired methods.
3. **Execution** – the virtual machine runs opcodes from `contracts_opcodes.go` and charges gas from `gas_table`.
4. **Cross-chain relay** – `cross_chain_transactions` forwards messages to other networks when required.

## Security Considerations
- Token standards include allowance checks to prevent unauthorized transfers.
- Gas accounting protects nodes from infinite loops or heavy computation.
- Cross-chain operations require relayer authorization and proof verification.

## CLI Integration
- `synnergy syn20`, `synnergy syn200`, etc. – manage specific token standards.
- `synnergy tx` – build and broadcast transactions.

## Enterprise Diagnostics
- New enterprise opcodes (`IntegrationDiagnostics`, `IntegrationConsensusProbe`, `IntegrationAuthoritySync`) are registered automatically during the integration probe, guaranteeing the gas table is aligned with documentation.
- Diagnostic transactions mined by the probe validate that ledger fee distribution and token transfers behave correctly before real workloads resume.
