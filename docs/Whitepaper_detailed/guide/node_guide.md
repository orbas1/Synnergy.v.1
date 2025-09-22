# Node Guide

This guide describes the node implementations available in the Synnergy project. Each node is a focused Go module that exposes a small API and can be composed with other components to build complex behaviour.

## Stage 82 Bootstrap Changes

Stage 82 ensures node lifecycles participate in the enterprise bootstrap flow.
When `synnergy orchestrator bootstrap` runs it starts the shared VM, seals the
wallet, performs a ledger audit and re-registers the orchestrator as an authority
node before any mining, staking or custody services initialise. The CLI and web
dashboard now expose the orchestrator's consensus relayer count and authority role
distribution, giving operators immediate visibility into node registration
health. Gas metadata for node operations—mining, staking, logistics, watchtower
and warfare—remains synchronised across documentation and runtime through the new
`registerEnterpriseGasMetadata` helper.

## Mining Node

The `MiningNode` simulates proof-of-work mining. It can start or stop hashing, submit block hashes, and attempt synchronous mining with adjustable difficulty. The node keeps minimal state so it can be used in tests and simulations.

**Key methods**
- `Start` / `Stop` – control the mining loop.
- `MineBlock` – mine a block at a specified difficulty.
- `SubmitBlock` and `LastBlock` – record and retrieve the most recently mined block.

## Mobile Mining Node

`MobileMiningNode` wraps a `MiningNode` and pauses mining when the device battery falls below a threshold.

**Key methods**
- `UpdateBattery` – record current battery level.
- `Start` – begin mining only when the battery is above the configured threshold.
- `SetThreshold` – adjust the minimum battery level required for mining.

## Staking Node

`StakingNode` tracks stakes for multiple addresses in memory.

**Key methods**
- `Stake` and `Unstake` – lock or release tokens.
- `Balance` – return the staked balance for an address.
- `TotalStaked` – amount staked across all participants.

## Indexing Node

`IndexingNode` offers a simple in-memory key/value index to speed up lookups of ledger data.

**Key methods**
- `Index` – insert or update entries.
- `Query` – fetch a value copy for a key.
- `Remove` – delete an entry.
- `Keys` and `Count` – enumerate the index.

## Content Network Node

`ContentNetworkNode` tracks content items hosted by a node, enabling resource discovery.

**Key methods**
- `Register` and `Unregister` – manage content availability.
- `Content` – fetch metadata for a specific item.
- `List` – return all registered content.

## Energy Efficient Node

`EnergyEfficientNode` records throughput and power usage for validators and can issue sustainability certificates.

**Key methods**
- `RecordUsage` – update processed transaction and energy metrics.
- `AddOffset` – credit carbon offsets.
- `Certify` – produce a `SustainabilityCertificate` summarizing efficiency.
- `ShouldThrottle` – report when efficiency drops below a threshold.

## Environmental Monitoring Node

`EnvironmentalMonitoringNode` evaluates sensor data against registered conditions and triggers actions when thresholds are met.

**Key methods**
- `SetCondition` – register sensor rules.
- `Trigger` – evaluate incoming readings against the stored condition.

## Geospatial Node

`GeospatialNode` collects geolocation history for subjects.

**Key methods**
- `Record` – store a latitude and longitude for a subject.
- `History` – retrieve recorded locations.

## Regulatory Node

`RegulatoryNode` cooperates with a `RegulatoryManager` to evaluate transactions and log flagged entities.

**Key methods**
- `ApproveTransaction` – check a transaction against regulations.
- `FlagEntity` and `Logs` – record and retrieve compliance events.

## Watchtower Node

`WatchtowerNode` monitors network health and detects forks.

**Key methods**
- `Start` / `Stop` – control monitoring routines.
- `ReportFork` – log fork events.
- `Metrics` – return collected system metrics.
- `Firewall` – access the embedded `Firewall` instance for rule management.

## Warfare Node

`WarfareNode` provides military-focused extensions such as signed command execution, logistics/tactical telemetry and commander governance.

**Key methods**
- `SecureCommand` – validate privileged commands.
- `TrackLogistics` and `Logistics` – record and query logistics events.
- `ShareTactical` – placeholder for distributing tactical information.

## Biometric Security Node

`BiometricSecurityNode` secures operations behind biometric authentication.

**Key methods**
- `Enroll` and `Remove` – manage biometric records.
- `Authenticate` – verify data for an address.
- `SecureExecute` – run a function only after successful biometric verification.

---

These building blocks can be combined to create rich node behaviour. Each node is a small, testable component that focuses on one responsibility, making the system easier to extend and maintain.
