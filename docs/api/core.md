# Core Package Reference

The core package implements consensus, transaction processing, networking and other primitives that underpin Synnergy. It is the foundation on which CLI commands, services and smart contracts are built.

## Modules

- **Consensus** – adaptive PoW/PoS/PoH engine with dynamic weighting
- **Ledger** – in‑memory and persistent block storage with WAL replay
- **Transactions** – structures and helpers for signing and validation
- **VM** – lightweight virtual machine that executes smart‑contract bytecode
- **Bridging** – cross‑chain registries and transfer managers

## Key Types

- `Block` – representation of a blockchain block
- `Transaction` – basic transaction structure
- `BridgeTransferManager` – coordinates cross‑chain deposits and claims
- `ProtocolRegistry` – tracks cross‑chain protocol definitions
- `ChainConnectionManager` – opens and monitors inter-chain connections
- `CrossChainRegistry` – stores contract address mappings across networks
- `ConsensusNetworkManager` – registers cross-consensus scaling networks
- `CustodialNode` – maintains off-chain asset custody records
- `EnterpriseSpecialNode` – aggregates heterogeneous node plugins into a combined enterprise control surface

## Interfaces

Several interfaces allow modules to be mocked or extended:

- `StateRW` – read/write access to the ledger state
- `VirtualMachine` – pluggable contract execution environment
- `OpContext` – context passed to VM opcodes

## Example Usage

Initialise a ledger and process a transaction:

```go
led := core.NewLedger()
tx := core.NewTransaction("alice", "bob", 10, 1, 0)
if err := led.ApplyTransaction(tx); err != nil {
    log.Fatal(err)
}
```

## Error Handling

Errors returned by core functions follow the conventions in [`docs/reference/errors_list.md`](../reference/errors_list.md). Always check returned errors and propagate context with `fmt.Errorf` when wrapping.

For additional details, read the source code under `core/`.

