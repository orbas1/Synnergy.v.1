# Consensus Guide

Synnergy employs a hybrid consensus engine that blends Proof‑of‑Work (PoW),
Proof‑of‑Stake (PoS) and Proof‑of‑History (PoH). The weights assigned to each
mechanism adjust dynamically according to network conditions so the chain can
favour security, throughput or decentralisation as required.

## Core Concepts

The implementation lives in `core/consensus.go` and related files. Key data
structures include:

- **ConsensusWeights** – percentage weight for PoW, PoS and PoH. Weights always
  sum to `1.0`.
- **SynnergyConsensus** – stateful engine that holds the current weights and
  tuning parameters used for transitions.

### Weight Adjustment

`AdjustWeights(D, S)` recalculates the internal weights based on network demand
`D` and stake concentration `S`. An adjustment factor `Gamma` smooths changes and
ensures no algorithm falls below a minimum floor unless unavailable. The method
normalises the weights after applying availability flags:

```go
sc.AdjustWeights(currentDemand, currentStake)
```

### Transition Threshold

The engine exposes helper functions to compute thresholds for switching
strategies:

- `Tload(D)` – influence of current network load.
- `Tsecurity(threat)` – impact of observed security threats.
- `Tstake(S)` – influence of stake concentration.
- `TransitionThreshold(D, threat, S)` – combined threshold used to decide when
  to adjust weights.

### Adaptive Management

Stage 63 introduces the `AdaptiveManager` which keeps a sliding window of recent
demand and stake metrics. Methods like `RecordMetrics` store observations while
`Adjust` and `Threshold` operate on the averaged values, smoothing sudden swings
in network behaviour.

### Difficulty Adjustment

PoW difficulty adjusts according to the time it took to mine the previous window
of blocks:

```go
newDiff := sc.DifficultyAdjust(oldDiff, actualTime, expectedTime)
```

### Validator Selection

For PoS rounds validators are selected randomly in proportion to their stake:

```go
next := sc.SelectValidator(stakeMap)
```

The function returns the address of the chosen validator or an empty string if
no stakes are present.

## CLI Interaction

Consensus parameters can be inspected and tuned via the `consensus` command
group once the `synnergy` CLI is built. Common operations include starting the
service, viewing current weights and simulating difficulty adjustments. See the
CLI help output for a full list of subcommands.

Stage 41 expands these tools with a `connpool release` command for freeing network connections and a `contractopcodes` helper that reports gas costs for contract operations, allowing operators to budget resources during tuning.

## Extending the Engine

The consensus package is designed for experimentation. Potential extensions
include:

- Integrating additional consensus mechanisms such as BFT replicas.
- Persisting weight history for audit and visualisation.
- Adding metrics to track how often transitions occur.

Contributors are encouraged to document significant changes in an Architecture
Decision Record under `docs/adr/`.
