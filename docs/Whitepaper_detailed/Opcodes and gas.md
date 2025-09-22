# Opcodes and Gas

## Deterministic Execution Model
Synnergy executes every runtime capability through a deterministic 24-bit opcode. Each entry in the catalogue couples a human-readable identifier with a binary constant so that governance, wallets, and automation pipelines can reason about behaviour before a single byte is dispatched. The dispatcher validates every opcode on boot, panicking on collisions so regression tests catch catalogue drift long before production deploys.

## Live Gas Registry
Stage 80 introduced an enterprise-grade gas registry. `gas_table.go` now exposes immutable `GasSnapshot` objects, subscription APIs, and hot reload support so that authority nodes, CLI tooling, and the new web console all stay synchronised without restarts. Updates propagate under writer locks, emit structured telemetry, and include version counters and timestamps for audit trails.【F:gas_table.go†L20-L152】【F:cmd/synnergy/main.go†L51-L96】 Runtime overrides—whether injected by governance or pre-flight tests—fan out to every subscriber and are persisted in the ledger documentation packs for post-incident review.

To prevent transient resets from erasing bespoke pricing, the VM layer mirrors the registry’s latest snapshot and persists a cached override map. Whenever a new snapshot arrives the override cache is refreshed, but if a subsequent reset clears the table the cached costs continue to apply until governance publishes a replacement, guaranteeing deterministic execution for long-running workflows.【F:virtual_machine.go†L57-L186】【F:core/virtual_machine.go†L48-L188】

## Catalogue Services
`snvm._opcodes.go` maintains the canonical opcode index. Stage 80 hardens the catalogue with build-time duplicate detection, O(1) lookup tables, and the new `SNVMOpcodeByCode` helper so both the VM and external auditors can recover human-readable metadata from machine code. Legacy entries that differ only in case or reuse an existing numeric code are gracefully ignored, allowing historical catalogues to coexist while tooling consumes a single authoritative mapping. The API is shared with `core/virtual_machine.go`, ensuring node operators and integrators observe the exact same mapping across CLI, VM, and web environments.【F:snvm._opcodes.go†L1-L1343】

## CLI and Automation Integration
The `synnergy` CLI now subscribes to gas updates at boot. It emits JSON-formatted telemetry for every change, re-registers critical opcodes to guard against documentation drift, and refuses to start if the VM bootstrap fails. Automation suites can therefore capture deterministic logs for compliance reports while developers gain immediate feedback when documentation and implementation diverge.【F:cmd/synnergy/main.go†L51-L163】

## Governance and Testing
Gas pricing is enforced end-to-end. Unit tests stress subscription fan-out, concurrency, and reload semantics, while the VM test suite proves that gas limits halt execution deterministically and respect hot reloads. Governance tooling consumes the same APIs through `core/gas_table.go`, guaranteeing parity between CLI inspection, dashboard snapshots, and ledger exports.【F:gas_table_test.go†L4-L99】【F:virtual_machine_test.go†L1-L96】【F:core/virtual_machine_test.go†L1-L108】

## Operational Playbook
Regulated operators manage pricing through a three-step workflow:

1. **Model** – update `docs/reference/gas_table_list.md` with proposed costs and submit the change for review.
2. **Dry-run** – use `synnergy gas set` in a staging cluster; subscribers validate the new snapshot through telemetry.
3. **Adopt** – once approved, governance publishes a signed snapshot, watchers propagate the update, and downstream wallets recalculate fee estimates automatically.

Every update is recorded in the ledger, signed by the issuing authority, and backed by unit and integration tests so that regulators and enterprise partners can audit the entire lifecycle.
