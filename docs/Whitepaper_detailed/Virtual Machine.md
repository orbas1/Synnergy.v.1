# Virtual Machine

## Overview
Stage 80 delivers a production-grade execution engine for Synnergy. The VM interprets 24-bit opcodes, enforces gas limits, and dispatches handlers with deterministic timing. Both the root VM (`virtual_machine.go`) and the core runtime (`core/virtual_machine.go`) share the same architecture: concurrency-limited execution, hot-reloadable opcode handlers, and subscription-driven gas lookups.【F:virtual_machine.go†L1-L186】【F:core/virtual_machine.go†L1-L310】

## Key Capabilities
- **Concurrency Control** – configurable limiter channels guarantee predictable resource usage across heavy, light, and super-light profiles.
- **Gas Subscriptions** – VM instances consume live gas snapshots and expose `Close()` semantics so operators can rotate nodes without leaking goroutines or stale pricing data.【F:virtual_machine.go†L93-L186】【F:core/virtual_machine.go†L62-L188】
- **Override Preservation** – cached overrides survive global gas table resets so long-running nodes and tests never fall back to default pricing mid-execution.【F:virtual_machine.go†L57-L186】【F:core/virtual_machine.go†L48-L188】
- **Opcode Metadata** – integration with `SNVMOpcodeByCode` allows tooling and auditors to resolve opcode names from machine code instantly while gracefully ignoring legacy duplicates so catalogue drift never reaches the dispatcher.【F:snvm._opcodes.go†L1-L1343】
- **Fallback Dispatch** – the core VM integrates with the opcode dispatcher so unregistered opcodes invoke catalogue handlers while still accruing deterministic gas.

## Tooling
Developers and operators can:
- Register custom handlers for testing or specialised workflows.
- Query VM gas snapshots to confirm synchronisation with governance decisions.
- Simulate execution with context cancellation to test failure scenarios and timeout behaviour.

## Testing and Assurance
Unit tests cover lifecycle management, custom handlers, gas hot reloads, and replay protection. Combined with the Synthron Coin integration tests, they provide strong guarantees that opcode execution remains deterministic across upgrades.【F:virtual_machine_test.go†L1-L96】【F:core/virtual_machine_test.go†L1-L108】【F:docs/Whitepaper_detailed/Synthron Coin_test.go†L1-L49】
