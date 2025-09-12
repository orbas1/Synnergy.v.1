# Block and Sub-Blocks

> *Neto Solaris – Synnergy Network Technical Whitepaper*

## Overview
Our blockchain architecture organizes activity into **sub-blocks** that are validated through Proof of Stake (PoS) and timestamped via Proof of History (PoH). These sub-blocks are then aggregated into full **blocks**, which are sealed using Proof of Work (PoW) to finalize state across the network.

## Transaction Intake and Validator Staking
Nodes accept transactions into a mempool only after verifying that the sender possesses sufficient balance and the transaction is well‑formed【F:core/node.go†L37-L56】. Stakeholders register validator stakes through `SetStake`, which enforces a minimum bonding requirement and tracks slashed validators to exclude them from selection【F:core/node.go†L98-L115】. Misconduct such as double‑signing or downtime can be reported to slash half the validator’s stake, while `Rehabilitate` allows recovery after penalties have been addressed【F:core/node.go†L117-L138】.

## Sub-Block Formation
Sub-blocks encapsulate ordered transaction batches signed by their originating validator. Each sub-block records the transactions it contains, the validator's identity, a PoH hash, timestamp, and cryptographic signature【F:core/block.go†L10-L17】. The `NewSubBlock` constructor computes a deterministic hash of the transaction IDs, validator string, and timestamp before signing the result, binding the validator to the data it proposes【F:core/block.go†L19-L37】. Validators can later verify authenticity by recomputing and comparing signatures using `VerifySignature`【F:core/block.go†L39-L43】.

## Block Assembly, Fee Distribution, and Mining
`MineBlock` orchestrates the full block lifecycle. The node selects a validator using weighted stake, converts the mempool into a single sub-block, validates it, and links the new block to the previous hash before initiating PoW【F:core/node.go†L58-L83】. Transaction fees are aggregated, partitioned among development, charity, loan pools, validators and other stakeholders, then adjusted for block utilization and distributed proportionally through a contract that credits each recipient on-ledger【F:core/node.go†L85-L92】【F:core/fees.go†L115-L127】【F:core/fees.go†L240-L255】. The final block hash and nonce are discovered via SHA‑256 PoW, securing the block’s contents【F:core/consensus.go†L190-L207】.

## Block Composition and Header Hashing
Blocks gather validated sub-blocks, reference the previous block hash, and track a nonce, timestamp, and final hash value【F:core/block.go†L45-L52】. The `HeaderHash` function constructs the PoW target by hashing the previous hash, each sub-block's PoH hash, and the timestamp–nonce pair, yielding the value miners test during PoW computations【F:core/block.go†L54-L69】.

## Consensus Interaction
`ValidateSubBlock` enforces PoS and PoH checks by confirming that the sub-block is populated and signed by its declared validator【F:core/consensus.go†L180-L188】. Beyond mining, the `SynnergyConsensus` engine dynamically adjusts PoW, PoS, and PoH weightings in response to network demand and stake concentration, exposes threshold calculations, and can toggle method availability or mining rewards for enterprise policy control【F:core/consensus.go†L55-L89】【F:core/consensus.go†L121-L142】.

## Adaptive Difficulty and Validator Selection
Consensus provides a weighted‑random validator selection algorithm that prevents dominance by large stakeholders and returns no result if a single party exceeds half of total stake【F:core/consensus.go†L144-L177】. Mining difficulty can be tuned using `DifficultyAdjust`, which recomputes the target based on actual versus expected block times, allowing operators to maintain predictable block intervals【F:core/consensus.go†L121-L128】.

## Genesis Block Initialization
Network bootstrapping occurs through `InitGenesis`, which credits the creator wallet, constructs an empty sub-block, assembles it into the inaugural block, and mines it with minimal difficulty. The function records chain height, block hash, circulating supply, remaining supply, and consensus weights for auditability【F:core/genesis_block.go†L19-L40】.

## Smart Contract and CLI Interfaces
Developers and operators interact with block construction via smart contract opcodes and command-line tooling. Opcodes such as `ProposeSubBlock`, `ValidatePoH`, and `SealMainBlockPOW` expose sub-block proposal and block sealing capabilities within on-chain logic【F:contracts_opcodes.go†L120-L132】【F:core/opcode.go†L340-L345】. The `block` CLI provides utilities to craft sub-blocks, verify signatures, assemble blocks, and calculate header hashes, while `consensus` commands mine blocks, tune weightings, and adjust difficulty for operational testing【F:cli/block.go†L16-L90】【F:cli/consensus.go†L19-L138】.

## Ledger Integration and Persistence
The ledger maintains account balances, a UTXO view, and block history. `AddBlock` appends new blocks while persisting them to a write-ahead log; `replayWAL` restores prior state on startup, and helpers such as `Head` and `GetBlock` enable efficient inspection of chain state【F:core/ledger.go†L35-L113】. This persistence layer underpins enterprise recovery strategies and supports auditing of every committed block.

## Synchronization and Integrity
A dedicated `SyncManager` coordinates block downloads and maintains synchronization status, allowing nodes to track whether they are up-to-date with the network【F:core/blockchain_synchronization.go†L1-L54】. The repository also includes a placeholder `block_integrity_check.sh` script, signalling ongoing efforts toward automated verification of stored block data【F:scripts/block_integrity_check.sh†L1-L17】.

## Security Considerations
Sub-block signatures tie proposed data to validator identities, while slashing functions deter misbehavior and can rehabilitate validators after penalties have been served【F:core/node.go†L117-L138】. PoW sealing safeguards against retroactive manipulation, and persistent ledger plus synchronization tools provide a robust audit trail for Neto Solaris stakeholders.

## Conclusion
By modularizing validation into sub-blocks and finalizing state through PoW, Neto Solaris's Synnergy Network delivers a layered security model that balances performance, auditability, and decentralization across its blockchain infrastructure.

