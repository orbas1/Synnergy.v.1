# How to Use the Synnergy Network Consensus

## Introduction
Blackridge Group Ltd. created the Synnergy Network to provide a resilient foundation for compliant decentralised finance. Its consensus layer combines proof-of-work (PoW), proof-of-stake (PoS) and proof-of-history (PoH) so operators can tune performance, security and energy profile without redeploying nodes. This guide explains how to interact with the consensus engine, configure its behaviour and participate in block production.

## Hybrid Consensus Architecture
Synnergy’s engine assigns a relative weight to each algorithm and continually re‑balances them using network demand and stake concentration. A fresh engine starts with a 40/30/30 weighting for PoW, PoS and PoH and exposes parameters \(Alpha, Beta, Gamma, Dmax, Smax\) to govern threshold calculus and smoothing【F:core/consensus.go†L11-L88】.

### Weight Parameters and Thresholds
`SynnergyConsensus` exposes `Threshold(D,S)` and `TransitionThreshold(D, threat, S)` to model load, security risk and stake concentration. `AdjustWeights` smooths changes using `Gamma`, clamps each mechanism to a 7.5\% floor and renormalises totals so weights always sum to one【F:core/consensus.go†L55-L119】.

### Dynamic Mode Hopping
The `ConsensusHopper` monitors transactions per second, average latency and the number of active validators, selecting the most appropriate mode for the next epoch. High throughput with sub‑second latency favours PoS, small validator sets tilt toward PoH, and all other conditions fall back to PoW【F:dynamic_consensus_hopping.go†L17-L70】. A read–write mutex protects mode selection and metric history so concurrent callers can query or override the hopper without races【F:dynamic_consensus_hopping.go†L24-L55】.

### Manual Overrides and Metric Inspection
Enterprise operators may temporarily pin the network to a single mode for maintenance windows or incident response. `SetMode` bypasses automatic heuristics and immediately switches to the requested mechanism, while `LastMetrics` returns the data that informed the most recent evaluation so administrators can justify overrides【F:dynamic_consensus_hopping.go†L43-L55】.

### Sub-block Workflow and Double-Lock Finality
During PoS or PoH epochs, block intervals are split into four sub‑steps: proposal, pre‑finalisation, PoW aggregation and finalisation. Pre‑finalised sub‑blocks provide near‑instant confirmation, while PoW aggregation secures the chain against long‑range attacks. Minted rewards are distributed 50 % to PoS validators, 10 % to PoH schedulers and 40 % to PoW miners, with validator payouts proportional to stake and participation【F:Synnergy_Network_Future_Of_Blockchan.md†L380-L416】.

Transaction fees follow a separate allocation path and are apportioned among validators, miners and other network funds according to the fee distribution contract【F:core/fees.go†L115-L125】.

Sub‑blocks encapsulate ordered transactions, validator identity and a PoH‑derived hash which is signed by the proposer; the consensus engine verifies that each sub‑block is non‑empty and bears a valid signature before inclusion in the final PoW block【F:core/block.go†L10-L43】【F:core/consensus.go†L180-L188】.

## Getting Started
### Prerequisites
- Go toolchain and access to the Synnergy CLI.
- Network connectivity to a ledger peer.
- Appropriate staking tokens if participating as a validator.

### Creating a Consensus-Specific Node
Use the `consensus-node` command to lock a node to a single consensus mode. The `create` subcommand accepts the mode (pow, pos or poh), an identifier and a network address. Subsequent helpers expose the node state through `info`, mine blocks with `mine`, and allocate stake on the local ledger via `stake <address> <amount>` so the node can participate in validator selection【F:cli/consensus_specific_node.go†L13-L76】.

### Deploying Validator Nodes
`ValidatorNode` bundles a base node with a `ValidatorManager` and `QuorumTracker`. Use `AddValidator`, `RemoveValidator` or `SlashValidator` to manage membership, and `HasQuorum` to ensure sufficient active stake before proceeding【F:core/validator_node.go†L5-L47】. `ValidatorManager` enforces minimum stake requirements and records slashing events, halving stake balances when `Slash` is invoked【F:core/consensus_validator_management.go†L11-L65】.


## Operating the Consensus Engine
### Inspecting and Adjusting Weights
The `consensus` command group exposes operational controls. Use `weights` to display the current distribution, `adjust [demand] [stake]` to recalculate weights, and `threshold` or `transition` to compute switching thresholds based on demand, threat level and stake concentration【F:cli/consensus.go†L15-L114】.

### Difficulty and Availability Controls
`difficulty [old] [actual] [expected]` recalculates PoW difficulty for a target block time using `old * (actual/expected)`【F:core/consensus.go†L121-L128】. `availability [pow] [pos] [poh]` toggles validator availability flags, and `powrewards [enabled]` enables or disables mining incentives without altering other modes【F:cli/consensus.go†L115-L180】.

### Mining a Block
For development or testing, `mine [difficulty]` creates a dummy block with an adjustable difficulty target and reports the discovered nonce【F:cli/consensus.go†L19-L36】.

## Validator Management
Synnergy selects validators using a weighted random algorithm that rejects sets where the largest stake exceeds half of the total, preventing dominance. Each sub‑block must include transactions and a valid validator signature before inclusion in a block【F:core/consensus.go†L144-L188】. The validator management layer maintains per‑address stake and slashing state and instruments add, remove and slash operations with telemetry spans for auditability【F:core/consensus_validator_management.go†L11-L65】. Validator nodes wrap these controls and quorum tracking so governance decisions immediately affect participation【F:core/validator_node.go†L23-L47】.


## Stake Management and Slashing
Operators may apply discretionary penalties or rewards using the `StakePenaltyManager`. `AdjustStake` modifies balances, `Penalize` appends timestamped records and `Info` returns the current stake with full history. The CLI exposes `stake_penalty slash <addr> <amount>` and `stake_penalty reward <addr> <amount>` for on‑chain enforcement【F:stake_penalty.go†L8-L61】【F:cli/stake_penalty.go†L15-L49】.

## High Availability and Failover
For enterprise resilience, `FailoverManager` tracks node heartbeats and promotes the freshest backup when the primary misses a configurable timeout. CLI helpers initialise the manager, register backups, record heartbeats and report the active node for orchestration systems【F:high_availability.go†L8-L69】【F:cli/high_availability.go†L14-L72】.

## System Health Monitoring
Blackridge ships a `SystemHealthLogger` for telemetry beyond consensus events. It samples goroutine counts, memory allocation, peer connectivity and block height, storing the latest snapshot for dashboards or alerting systems【F:core/system_health_logging.go†L11-L44】.

## Security, Transparency and Telemetry
All mode changes and validator metrics are logged for public inspection. Administrators may pin modes or adjust weights, but every decision is recorded, enabling external auditors to replay consensus events and verify fairness【F:Synnergy_Network_Future_Of_Blockchan.md†L384-L436】. Core validator operations emit OpenTelemetry spans so monitoring systems can trace stake changes and slashing actions in real time【F:core/consensus_validator_management.go†L28-L57】【F:internal/telemetry/telemetry.go†L1-L10】.

## Best Practices
- Monitor `ConsensusHopper` metrics to anticipate mode shifts and plan infrastructure capacity.
- Maintain diverse validator participation to avoid triggering the selection guardrail and to maximise reward share.
- Review telemetry feeds and logs regularly for signs of latency spikes or validator churn.
- Configure failover heartbeats and promote backups automatically to keep block production uninterrupted.
- Audit stake penalty history to detect recurring misbehaviour and refine governance policies.
- Stream system health snapshots to watchtower or equivalent observability stacks to forecast resource needs.

## Conclusion
The Synnergy Network consensus, engineered by Blackridge Group Ltd., delivers flexible, observable and secure block production. By understanding the hybrid architecture and operational tooling described above, network participants can confidently deploy nodes, manage staking positions and tune performance to meet evolving requirements.
