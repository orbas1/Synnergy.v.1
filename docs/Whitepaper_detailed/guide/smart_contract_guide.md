# Synnergy Smart Contract Guide

## Overview

Synnergy's contract platform executes WebAssembly (WASM) bytecode in a pluggable virtual machine. Contracts are registered through the on-chain `ContractRegistry` and executed deterministically with gas metering. This guide covers authoring, compiling, deploying and invoking contracts on a Synnergy node.

## Stage 82 Execution Improvements

Stage 82 adds safety rails for contract execution. `registerEnterpriseGasMetadata`
ensures contract opcodes share the same descriptive metadata as the rest of the
runtime, so the CLI, documentation and JavaScript control panel expose identical
gas costs when deploying or invoking bytecode. `bootstrapRuntime` registers a VM
execution hook that logs opcode failures with opcode names, gas cost and
remaining gas, allowing auditors to trace problematic contracts without enabling
debug instrumentation. During `synnergy orchestrator bootstrap` the heavy VM is
started, the orchestrator wallet is sealed and the ledger is audited before any
contract deployment occurs, giving developers confidence that WASM execution
targets a verified environment.

Contracts can be written in any language that compiles to WASM64. The repository includes Solidity examples compiled with a compatible toolchain. Rust and AssemblyScript are also common choices. Once compiled, bytecode is hashed to derive a deterministic contract address and may optionally include a Ricardian manifest binding legal terms to the code hash.

## Directory Structure

- `core/contracts.go` – runtime registry and deployment helpers.
- `core/virtual_machine.go` – virtual machine interfaces and in-memory state used during execution.
- `core/contracts_opcodes.go` – opcode constants reserved for contract operations.
- `cli/contracts.go` – CLI commands for compiling, deploying and invoking contracts.
- `cli/ai_contract.go` – Deploy and invoke AI-enhanced contracts with model tracking.
- `cli/audit.go` and `cli/audit_node.go` – Record and query audit logs from the CLI.
- `cmd/smart_contracts/` – example Solidity contracts demonstrating custom opcodes.

Review these files for implementation details and additional inline comments.

## Writing Contracts

1. **Choose a language** that targets WASM. Rust with `wasm-pack` or AssemblyScript are recommended for low level control.
2. **Follow the gas model** defined in `core/gas_table.go`. Each opcode consumes a fixed amount of gas. Custom opcodes such as `Liquidity_AddLiquidity` and `QueryOracle` are assigned specific costs.
3. **Interact with chain services** by calling exported opcodes via inline assembly. Examples in `cmd/smart_contracts` show how to mint tokens, add liquidity or query an oracle using assembly `call` instructions.
4. **Keep code deterministic** – avoid sources of randomness or floating point arithmetic that could diverge across nodes.

### Ricardian Manifest

A Ricardian contract is a JSON file that links legal prose to a specific code hash. When deploying a contract you may supply a manifest containing fields such as:

```json
{
  "name": "TokenMinter",
  "version": "1.0",
  "author": "Example Corp",
  "terms": "<link to PDF or legal text>",
  "hash": "<sha256 of wasm bytecode>"
}
```

The CLI stores this manifest on the ledger so anyone can audit the deployed code and its associated terms.

## Compiling

Use the `contracts compile` command to convert a WAT or WASM file into a deterministic byte blob. The command invokes `wat2wasm` if a `.wat` source is provided. Output is written to the directory specified by `WASM_OUT_DIR` (defaults to `./wasm`). Example:

```bash
synnergy contracts compile ./examples/hello.wat
```

This prints the output path and hash of the compiled artifact. The hash is later used to derive the contract address.

## Deployment

Deploy a compiled WASM binary with optional Ricardian manifest:

```bash
synnergy contracts deploy --wasm ./wasm/hello.wasm --ric ./manifest.json --gas 3000000
```

Deployment registers the contract in the `ContractRegistry`, persists bytecode to the ledger and sets the maximum gas limit. The resulting address is printed to stdout.

## Invocation

Invoke methods using the `contracts invoke` subcommand. Arguments are passed as hex bytes and a gas limit must be specified:

```bash
synnergy contracts invoke 0xabc... --method greet --args 48656c6c6f --gas 200000
```

The registry locates the contract, executes it inside the VM and returns any bytes produced by the call. Use `contracts list` to see all deployed addresses or `contracts info <addr>` to display stored Ricardian metadata.

## Example Contracts

Several Solidity examples live under `cmd/smart_contracts`:

- **cross_chain_eth.sol** – simple Ether bridge emitting events for off‑chain relayers.
- **liquidity_adder.sol** – demonstrates calling the custom `Liquidity_AddLiquidity` opcode.
- **multi_sig_wallet.sol** – basic multisig wallet requiring multiple confirmations.
- **oracle_reader.sol** – queries a Synnergy oracle via opcode `QueryOracle`.
- **token_minter.sol** – mints tokens using the `MintToken` opcode.

These files showcase how inline assembly can access Synnergy-specific opcodes that extend beyond standard Ethereum functionality.

## Testing

Unit tests for the contract system reside in `tests/contracts_test.go`. They illustrate basic deployment and invocation flows using the in-memory ledger provided by `core`. Run tests with:

```bash
go test ./tests -run Contracts
```

Ensure all tests pass before deploying contracts on a live network.

## Best Practices

- **Version Control** – commit both source and compiled WASM to guarantee reproducible builds.
- **Gas Limits** – estimate execution cost and set generous gas limits when deploying. Excess gas is refunded.
- **Security Reviews** – audit contracts thoroughly; once deployed code cannot be replaced except through a migration process.
- **Logging** – use the logging facilities in `core` to emit structured logs from your contract for easier debugging.
- **Static Calls** – where possible, read contract state with `StaticCall` to avoid unintended mutations.

## Further Reading

- [`cli_guide.md`](cli_guide.md) – overview of all CLI command groups.
- [`module_guide.md`](module_guide.md) – descriptions of core modules including the VM and ledger.
- [`config_guide.md`](config_guide.md) – details on configuring a Synnergy node.
- `tests/` – unit tests demonstrating module usage.

With these resources you can author robust smart contracts on Synnergy, deploy them to the ledger and interact with them via the provided CLI.
