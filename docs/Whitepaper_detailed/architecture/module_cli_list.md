# Stage 81 – Enterprise Module and CLI Alignment

Stage 81 hardens the Synnergy function web by aligning every mission-critical CLI
surface with its backing core module, gas schedule and documentation. The upgrade
introduces a deterministic module catalogue, surfaces gas metadata through both
the CLI and web control plane, and refreshes architectural guidance so operators
can deploy permissioned networks with auditable governance, privacy and
interoperability guarantees.

## Enterprise Command Matrix

| Domain | CLI Command | Responsibilities | Key Opcodes (Gas) |
| --- | --- | --- | --- |
| Consensus Control | `synnergy consensus` | Mine blocks, adjust hybrid PoW/PoS/PoH weights, compute failover thresholds, toggle validator availability and PoW rewards. | `MineBlock` (50), `AdjustWeights` (1), `TransitionThreshold` (1), `SetAvailability` (4), `SetPoWRewards` (4) |
| Virtual Machine | `synnergy simplevm` | Provision deterministic VM instances, enforce execution timeouts, meter gas consumption and decode WASM outputs for on-chain auditability. | `VMCreate` (15), `VMStart` (10), `VMStop` (6), `VMStatus` (4), `VMExec` (45) |
| Wallet & Identity | `synnergy wallet` / `synnergy idwallet` | Generate encrypted operator wallets, register identity wallets for permissioned access and validate digital signatures. | `WalletNew` (20), `VerifySignature` (1), `IDWalletRegister` (1) |
| Validator Nodes | `synnergy node` | Inspect validator health, stake validators, apply slashing, rehabilitate nodes, inject transactions and mine blocks with deterministic gas metering. | `NodeInfo` (4), `NodeStake` (25), `NodeSlash` (30), `NodeRehab` (12), `NodeAddTx` (8), `NodeMempool` (3), `NodeMine` (55) |
| Authority Governance | `synnergy authority` | Register authority nodes, capture governance votes, rotate terms and maintain quorum controlled leadership. | `AuthorityApplyVote` (5), `RenewAuthorityTerm` (7), `UpdateMemberRole` (5) |
| Module Catalogue | `synnergy modules list` | Enumerate enterprise modules, validate opcode documentation and expose gas metadata for CLI, VM and JavaScript orchestration. | `ModuleCatalogueList` (3), `ModuleCatalogueInspect` (5) |

The catalogue powers the JavaScript UI and the CLI simultaneously—operators can
issue `synnergy modules list --json` to stream structured data directly into
automation pipelines or monitoring dashboards.

## Workflow Integration

### CLI & Core Modules
- **Fault tolerance** – The catalogue verifies that every module exposes the
  opcodes required for consensus, VM, wallet, node and authority workflows.
  Missing documentation is surfaced immediately so upgrades cannot regress the
  runtime.
- **Security & privacy** – Wallet creation, authority elections and identity
  registration enforce digital signature validation and permissioned access while
  recording gas costs for audit trails.
- **Scalability** – Validator, consensus and VM commands expose deterministic gas
  pricing to guarantee predictable performance during stress, situational and
  real-world tests.

### Virtual Machine & Consensus
- VM lifecycle operations now share gas metadata with consensus mining to ensure
  that cross-module orchestration never exceeds configured block gas limits.
- Consensus commands emit the same opcode data consumed by the VM so governance
  policies remain enforceable even when validator counts surge.

### Wallet, Node and Authority Bridges
- Wallet generation integrates with staking and authority CLI flows so
  enterprise deployments can enforce multi-party governance from the command
  line.
- Authority node registration is now referenced directly by the catalogue so the
  governance whitepaper, CLI help and gas table stay synchronized.

## Function Web and JavaScript UI

A dedicated API endpoint (`/api/modules`) executes `synnergy modules list
--json`, normalises the output and feeds the React dashboard. The homepage now
renders the module matrix alongside orchestrator telemetry, providing a unified
view into gas coverage, opcode health and participant readiness without leaving
the browser.

## Testing and Validation

New Go unit tests exercise the module catalogue, JSON rendering and consensus
coverage, while existing stress, situational and functional suites reuse the
shared metadata. The gas table and opcode reference have been extended for every
Stage 81 opcode, ensuring that consensus, VM, wallet, node and authority
integrations remain cryptographically secure, interoperable and compliant with
enterprise policy.
