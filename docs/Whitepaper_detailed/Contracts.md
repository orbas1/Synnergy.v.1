# Contracts

## Introduction
Blackridge Group Ltd. designs the Synnergy Network to provide a deterministic and enterprise‑grade smart contract platform. This document details the contract framework, from compilation and deployment to management, security, and cross‑chain interoperability.

## Contract Lifecycle
### Compilation
- Contracts are authored in a supported high‑level language and compiled to WebAssembly (WASM).
- `CompileWASM` validates the bytecode and produces a SHA‑256 hash used for deterministic addressing.

### Deployment
- `Deploy` registers a contract by its bytecode hash, stores optional Ricardian manifest data, and enforces per‑contract gas limits.
- The resulting address is collision‑resistant and derived solely from bytecode, ensuring predictable deployments across environments.

### Invocation
- `Invoke` executes a method on the contract through the configured virtual machine, automatically enforcing the gas ceiling and checking the paused state.
- Output data and consumed gas are returned for auditability.

### Deterministic Hex Encoding
- `CompileWASM` emits a 32‑byte SHA‑256 digest and exposes it as a 64‑character hex string that serves as the canonical contract identifier.
- The address is deterministic: redeploying identical bytecode always yields the same 64‑character value, allowing offline build verification and protecting against address squatting.
- `Deploy` persists this hex hash as the on‑chain address and stores the original bytecode, enabling reproducible builds and content‑addressable lookups across nodes.
- Registries such as `ContractRegistry`, `AIContractRegistry`, and `XContractRegistry` index contracts by this hex string, making it the universal reference for upgrades and cross‑chain mappings.
- CLI helpers like `contracts info` surface the hex address alongside metadata, while `opcodes hex` resolves function names to their 32‑bit hex opcodes for deterministic debugging.
- Opcode tables in `contracts_opcodes.go` and `snvm._opcodes.go` encode every VM function as a `0x` prefixed value. The `opcodeByName` helper and CLI tooling use these hex codes for gas accounting and interoperability audits.

## Contract Registry
- `ContractRegistry` maintains all deployed contracts in a concurrency‑safe map and exposes helper methods for listing or retrieving entries.
- `VirtualMachine` defines a minimal execution interface (`Execute`, `Start`, `Stop`, `Status`), allowing modular VM implementations.

## Contract Management
Administrative control is provided through `ContractManager`:
- **Transfer** – changes ownership of a contract without redeployment.
- **Pause/Resume** – toggles execution to mitigate emergencies or upgrades.
- **Upgrade** – replaces WASM bytecode and optionally adjusts gas limits while preserving the original address.
- **Info** – returns metadata including owner, manifest, gas limit, and paused status.

## Supported Languages and Tooling
`SupportedContractLanguages` and `IsLanguageSupported` provide a canonical whitelist that the VM recognises for compilation and execution. The CLI selects the appropriate toolchain based on file extension and rejects anything outside this list.

| Language | Toolchain | Notes | Repository Examples |
| --- | --- | --- | --- |
| **Wasm/WAT** | Native ingestion | One‑to‑one mapping to opcodes for low‑level optimisation and testing | Pre‑compiled `.wasm` templates under `smart-contracts/` |
| **Go** | TinyGo → `wasm32-wasi` | Reuses Go libraries with deterministic builds and static analysis support | Example modules compiled into `.wasm` artifacts |
| **JavaScript/AssemblyScript** | AssemblyScript compiler | Enables rapid prototyping with TypeScript‑like syntax while emitting deterministic binaries | N/A (compiled outputs only) |
| **Solidity** | Solang cross‑compiler | Allows existing EVM contracts to target the Synnergy VM without major rewrites | Source library in `smart-contracts/solidity/` |
| **Rust** | Cargo (`wasm32-unknown-unknown`) | Memory‑safe, high‑performance contracts; integrates with the `smart-contracts/rust` workspace | `smart-contracts/rust/` |
| **Python** | Pyodide/MicroPython transpilers | Favoured for AI and data‑science workloads where dynamic typing aids experimentation | Compiled WASM examples such as `ai_model_market.wasm` |
| **Yul** | Solidity Yul IR → WASM | Minimal intermediate representation for audit‑friendly, gas‑efficient bytecode | Used for specialised low‑level modules |

### Benefits of Multi‑Language Support

Polyglot contract authoring removes linguistic barriers and lets teams choose the toolchain that best matches their expertise. Web developers can contribute AssemblyScript modules, financial engineers can port Solidity code through Solang, and data scientists can ship Python‑driven AI models—all while targeting the same runtime.

Standardising on WebAssembly and a 64‑character hex address ensures these diverse languages interoperate seamlessly. Regardless of source language, `CompileWASM` produces identical bytecode and the same deterministic address, allowing audit tools and registries to treat every contract uniformly. This uniformity reduces migration friction from other chains, protects existing developer investments, and invites a global contributor base, ultimately making the Synnergy Network more universally accessible and resilient.

The command‑line interface enables lifecycle operations via `contracts compile|deploy|invoke|list|info` and surfaces syntax‑specific tips based on the selected language. Extending `SupportedContractLanguages` allows future addition of languages such as Move or Vyper following the same deterministic pipeline.

## Opcode Architecture and Gas
- `ContractOpcodes` enumerates a stable opcode set spanning AI, governance, financial primitives, compliance, consensus, and bridging.
- Each opcode maps to a deterministic numeric identifier used in gas accounting and tooling, facilitating clear cost analysis across modules.
- `snvm._opcodes.go` mirrors this table for the broader VM, ensuring CLI `opcodes hex` queries and internal gas metering reference a single source of truth.

## AI‑Enhanced Contracts
- `AIContractRegistry` extends the base registry by associating deployed contracts with model hashes for inference.
- `DeployAIContract` stores both the WASM and corresponding AI model hash, while `InvokeAIContract` standardizes calling the `infer` method for on‑chain machine‑learning workflows.

## Cross‑Chain Contract Mapping
- `XContractRegistry` links local contract addresses to remote chains, supporting bidirectional interoperability.
- Mappings can be registered, listed, queried, and removed, enabling bridge modules to orchestrate cross‑chain calls with deterministic references.

## Sample Contract Library
- The repository ships with an extensive `smart-contracts/solidity` suite demonstrating tokens, DeFi primitives, governance systems, data oracles, AI marketplaces, and more.
- Examples include `AIModelMarket.sol`, `ArbitratedEscrow.sol`, `StakingRewards.sol`, `WeatherOracle.sol`, and `YieldFarm.sol`, offering production‑ready templates for rapid development.

## Security and Governance
- Deterministic addressing, per‑contract gas limits, and the ability to pause or upgrade contracts provide strong operational controls.
- Opcode segregation and language whitelisting assist regulators and auditors in verifying contract behavior within Blackridge Group Ltd.’s compliance framework.

## Conclusion
Through a modular registry, rich opcode vocabulary, AI integrations, and cross‑chain mapping, the Synnergy Network delivers a comprehensive contract ecosystem aligned with Blackridge Group Ltd.’s mission to power secure, interoperable, and intelligent blockchain solutions.
