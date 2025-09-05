# Core Package Reference

The core package implements consensus, transaction processing, and network primitives.

## Key Types
- `Block` – representation of a blockchain block
- `Transaction` – basic transaction structure
- `BridgeTransferManager` – coordinates cross-chain deposits and claims
- `ProtocolRegistry` – tracks cross-chain protocol definitions
- `ChainConnectionManager` – opens and monitors inter-chain connections
- `CrossChainRegistry` – stores contract address mappings across networks
- `ConsensusNetworkManager` – registers cross-consensus scaling networks
- `CustodialNode` – maintains off-chain asset custody records

## Functions
- `ValidateBlock` – verifies block integrity
- `ApplyTransaction` – applies a transaction to the current state

For additional details, read the source code under `core/`.
