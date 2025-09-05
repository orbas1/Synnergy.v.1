# Opcodes and Gas

## Overview
At Blackridge Group Ltd., the Synnergy network executes every function through a deterministic opcode. Each opcode is a 24‑bit identifier that maps directly to a smart‑contract or system routine, allowing wallets and dashboards to reason about behaviour and cost before execution. Gas fees are priced centrally so that network actors share a consistent view of execution costs across releases.

## Opcode Encoding
Every instruction occupies three bytes. The high byte denotes the protocol category and the remaining two bytes form a sequential ordinal, yielding the canonical layout `0xCCNNNN`【F:core/opcode.go†L6-L8】. The dispatcher exposes helpers for converting these values to hexadecimal or raw bytecode and for validating inbound byte streams, ensuring tooling can safely serialise or parse opcodes【F:core/opcode.go†L1660-L1676】.

### Conversion Utilities
Synnergy ships multiple helpers to translate between human‑readable names and their binary representation. `ParseOpcode` converts a three‑byte slice into a numeric opcode and rejects malformed lengths, while `MustParseOpcode` provides a panic‑on‑error variant for tooling【F:core/opcode.go†L1675-L1690】. For outbound workflows, `ToBytecode` and `HexDump` expose canonical encodings so dashboards can embed opcodes directly in transactions or documentation. `DebugDump` offers a lexicographically sorted list of name‑to‑hex pairs for on‑chain auditing tools and CLI inspection【F:core/opcode.go†L1692-L1728】.

## Opcode Catalogue
### Synnergy Virtual Machine
The Synnergy Virtual Machine (SNVM) exposes a consolidated catalogue of opcodes. The catalogue is continually expanded to incorporate new modules—such as consensus, node management and cross‑chain flows—and assigns a stable code to each exported function【F:snvm._opcodes.go†L3-L15】.

#### Catalogue Coverage
`snvm._opcodes.go` enumerates every function compiled into the VM, spanning AI marketplaces, node orchestration, cross‑chain bridges and hundreds of token utilities. Each entry declares a descriptive identifier and its 24‑bit code so external systems can reason about capability coverage before execution【F:snvm._opcodes.go†L1-L16】.

### Normalisation and Registration
During initialisation the opcode catalogue is normalised. Sequential identifiers are generated per category and handlers are registered with the dispatcher, guaranteeing unique values at runtime. Once the catalogue is finalised, the gas table is constructed so that each opcode has an associated cost【F:core/opcode.go†L1624-L1653】.

### Module Constants
Critical opcodes are exposed as variables so dependants can reference stable values without embedding magic numbers. `contracts_opcodes.go` resolves contract‑management operations—such as deployment and pause/resume—through the catalogue at start‑up, panicking if a name is missing to catch mismatches during CI【F:core/contracts_opcodes.go†L1-L24】.

### Catalogue Introspection
The dispatcher publishes the full mapping for external auditors. `Catalogue` returns a slice of opcodes with their current gas price, while `Opcodes` exposes a map of numeric codes to human‑readable names for wallet discovery【F:core/opcode.go†L71-L85】. When the virtual machine executes an instruction, `Dispatch` resolves the handler and charges its base gas before invocation, guaranteeing consistent metering across clients【F:core/opcode.go†L106-L118】.

### Reference Lists
For external tooling, the project maintains human‑readable tables of opcodes under `docs/reference`. These lists document function names alongside their hex codes to simplify integration with monitoring and wallet software【F:docs/reference/opcodes_list.md†L1-L5】.

## Gas Pricing Framework
### Central Gas Table
Gas pricing is derived from a canonical table that maps opcode names to their base cost. The table is loaded from `docs/reference/gas_table_list.md` and cached on first use. Any opcode not explicitly priced receives the `DefaultGasCost` of one, ensuring new features remain executable without governance intervention【F:core/gas_table.go†L15-L33】【F:docs/reference/gas_table_list.md†L1-L6】.

### Loading and Caching
The table is parsed once and stored in memory using a `sync.Once` gate so that all threads share a consistent view of pricing. The parser streams `gas_table_list.md` line by line, trimming Markdown syntax and extracting backtick‑quoted names with numeric costs. Malformed rows silently fall back to the default cost so an incomplete guide never halts node start‑up【F:core/gas_table.go†L34-L88】. `LoadGasTable` then returns the cached map under read locks to keep lookups race‑free【F:core/gas_table.go†L90-L96】.

### Runtime API
The core package exposes utilities to adjust and inspect pricing at runtime. Operators or tests can set costs for individual opcodes, generate snapshots for audit, serialise the complete schedule to JSON, or query pricing by exported function name for high‑level tooling【F:core/gas_table.go†L70-L122】. `GasTableSnapshotJSON` emits hex‑keyed maps for downstream analytics, while `GasCostByName` resolves friendly function names through the opcode catalogue to retrieve current prices【F:core/gas_table.go†L93-L122】.

### Concurrency and Overrides
Gas pricing is guarded by read/write mutexes so concurrent queries remain safe【F:core/gas.go†L5-L28】. Helper functions report whether an opcode is priced, inject overrides at runtime and reset the cache for test isolation, giving governance processes fine‑grained control over execution fees【F:gas_table.go†L98-L133】.

## Operational Tooling
A dedicated CLI command allows authorised personnel to tune gas values or capture the current schedule. Costs can be modified with `synnergy gas set`, and snapshots can be exported in plain text or JSON to feed dashboards and regulatory systems【F:cli/gas_table.go†L11-L55】.

## Testing and Audit
Extensive tests confirm that gas overrides and snapshots behave deterministically. The suite verifies that unpriced opcodes default to zero gas, that updating a cost reflects immediately, that snapshots are immutable copies, and that named lookups fall back to default pricing when unknown【F:core/gas_test.go†L5-L12】【F:core/gas_table_test.go†L80-L99】. These guarantees ensure external systems can trust published gas schedules.

## Reference Guides
For comprehensive listings of every opcode and its baseline cost, the repository maintains Markdown references under `docs/reference`. `opcodes_list.md` maps human‑readable function names to their canonical hex codes, while `gas_table_list.md` enumerates the default pricing schedule【F:docs/reference/opcodes_list.md†L1-L8】【F:docs/reference/gas_table_list.md†L1-L8】. A more narrative explanation of operational categories and pricing rationales lives in `docs/Whitepaper_detailed/guide/opcode_and_gas_guide.md`, providing deep context for auditors and integrators【F:docs/Whitepaper_detailed/guide/opcode_and_gas_guide.md†L1-L5】.

## Summary
Through disciplined opcode cataloguing and transparent gas accounting, Blackridge Group Ltd. delivers predictable execution semantics across the Synnergy ecosystem. Deterministic pricing, runtime configurability and comprehensive tooling allow developers, operators and regulators to validate costs with confidence.
