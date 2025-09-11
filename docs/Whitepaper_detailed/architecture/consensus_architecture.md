# Consensus Architecture

## Overview
Consensus modules decide how blocks are produced and agreed upon across heterogeneous nodes. The system can hop between proof‑of‑work, proof‑of‑stake and proof‑of‑history modes based on live network conditions.

## Key Modules
- `dynamic_consensus_hopping.go` – selects an optimal consensus mode from metrics.
- `core/consensus` components – manage difficulty, validator sets and block validation.
- `core/consensus.go` – provides VRF-based validator selection, BFT finality and fork-choice.
- `consensus_service.go` – runs mining services that participate in block production.
- `consensus_specific_node.go` – locks a node into a fixed mode when required by policy.
- `consensus_adaptive_management.go` – adjusts weights used for dynamic mode selection.

## Workflow
1. **Metrics collection** – nodes feed latency and stake data to `dynamic_consensus_hopping`.
2. **Mode evaluation** – the hopper chooses a consensus mode and informs validators.
3. **Block production** – the `consensus_service` mines or validates blocks according to the selected mode.
4. **Difficulty adjustment** – core consensus utilities recalculate thresholds to maintain target block times.
5. **Validator management** – `validator_management` registers, slashes or rewards participants.
6. **Finality voting** – `FinalizeBlock` aggregates validator votes, seals a block once two thirds agree and credits stake rewards to contributing validators.

## Security Considerations
- Mode switches require thresholds to prevent rapid oscillation and potential exploits.
- Validator slashing protects against double-signing or extended downtime.
- Availability checks ensure no single mode becomes a bottleneck during network splits.

## CLI Integration
- `synnergy consensus` – inspect weights, adjust parameters and view thresholds.
- `synnergy consensus-service` – start or stop a mining service.
- `synnergy consensus-node` – create a node fixed to a particular mode.
- `synnergy consensus-adaptive` – query and tune adaptive management weights.
