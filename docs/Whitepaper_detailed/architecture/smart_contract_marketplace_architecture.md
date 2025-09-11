# Smart-Contract Marketplace Architecture

## Overview
The smart-contract marketplace lets developers publish and trade reusable WebAssembly contracts. It leverages the existing contract registry and virtual machine so uploaded bytecode executes with the same determinism and gas accounting as native contracts.

## Key Modules
- `cli/contracts.go` – deploys templates and manages marketplace listings.
- `smart-contracts/templates` – collection of example contracts available for reuse.
- Core contract registry routines accessed through opcodes.

## Workflow
1. **Template selection** – developers pick a template or upload custom WASM bytecode.
2. **Deployment** – `synnergy contracts deploy` registers the contract and stores metadata.
3. **Listing** – optional marketplace listing makes the contract discoverable to others.
4. **Invocation** – consumers use standard contract calls; gas usage is tracked by the VM.

## Security Considerations
- Uploaded bytecode is validated and sandboxed by the virtual machine.
- Listings reference immutable hashes so users receive the expected code.
- Gas limits prevent runaway execution and are surfaced in CLI output.

## CLI Integration
- `synnergy contracts` – deploy, list and remove marketplace contracts.
