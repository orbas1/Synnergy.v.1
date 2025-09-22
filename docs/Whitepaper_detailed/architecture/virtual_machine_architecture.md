# Virtual Machine Architecture

## Overview
The Synnergy Virtual Machine (SVM) executes WebAssembly smart contracts in a deterministic, resource‑constrained environment. It supports sandboxed execution, opcode metering and context-aware timeouts.

## Key Modules
- `virtual_machine.go` – core VM implementation handling instruction dispatch and gas metering.
- `vm_sandbox_management.go` – creates and tears down sandboxes for contract calls.
- `snvm._opcodes.go` – opcode table mapping functions to gas costs.
- `contracts_opcodes.go` – exposes opcodes for contract deployment and invocation.

## Workflow
1. **Instantiation** – the VM spins up a sandbox with defined memory and time limits.
2. **Execution** – opcodes from the contract are interpreted; gas usage is tracked in `gas_table.go`.
3. **Completion** – `vm_sandbox_management` cleans up memory and returns results.
4. **Snapshotting** – gas table snapshots can be written for analysis or billing.

## Security Considerations
- Sandboxes prevent contracts from accessing host resources directly.
- Timeouts stop runaway code and free resources for other transactions.
- Opcodes are versioned to ensure deterministic behaviour across nodes.

## CLI Integration
- `synnergy vm` – inspect VM settings and manage snapshots.

## Enterprise Diagnostics
- `synnergy integration status` runs the VM in heavy mode and exposes `mode` plus `concurrency` through its JSON payload, enabling external automation to verify execution capacity.
- The command also registers enterprise opcode defaults ensuring VM and gas table documentation stay synchronized.
