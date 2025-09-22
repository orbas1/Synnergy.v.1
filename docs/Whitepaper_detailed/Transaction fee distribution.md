# Transaction Fee Distribution

## Fee Collection
All transactions incur deterministic gas fees drawn from the central registry. Fees are charged at execution time in the VM, recorded in the ledger, and exposed through CLI telemetry so operators can audit revenue streams in real time.【F:virtual_machine.go†L1-L186】【F:cmd/synnergy/main.go†L51-L163】

## Allocation Model
1. **Validator Rewards** – a configurable percentage rewards validators for securing consensus and hosting archival data.
2. **Treasury Fund** – a governance-controlled treasury receives a portion to finance upgrades, grants, and compliance initiatives.
3. **Insurance Pool** – residual fees accrue to an insurance pool that backs reversal operations and customer protection programmes.

Allocation ratios are codified in governance policy and enforced during distribution runs. Ledger entries include digital signatures and metadata detailing block height, total fees, and allocation percentages, providing a verifiable audit trail.【F:docs/Whitepaper_detailed/Synthron Coin.go†L1-L167】

## Operational Process
- Fee data is exported via signed JSON reports for enterprise accounting systems.
- Authority nodes can simulate adjustments in staging clusters using the same subscription APIs, ensuring changes propagate instantly without service restarts.【F:gas_table.go†L20-L152】
- Tests verify that fee snapshots remain stable across gas reloads and that failed executions never leak unaccounted fees.【F:virtual_machine_test.go†L1-L96】【F:core/virtual_machine_test.go†L1-L108】
- SYN500 telemetry binds service-tier usage to the active gas version so entitlement fees can be reconciled alongside validator, treasury, and insurance distributions.【F:core/syn500.go†L1-L150】【F:cli/syn500.go†L1-L172】
