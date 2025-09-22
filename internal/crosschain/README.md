# Cross-Chain Services

The cross-chain package coordinates bridge operations, settlement proofs and
state synchronisation with external ledgers. It enables the Synnergy CLI and web
interface to execute asset transfers, governance actions and oracle updates
across heterogeneous ecosystems without compromising gas determinism or replay
protection.

## Capabilities

* **Bridge Orchestration** – Manages validator attestations, signature
  aggregation and encrypted payload distribution when tokens are wrapped or
  redeemed across partner networks.
* **Gas Harmonisation** – Consumes the shared gas tables defined in
  `internal/config.VMConfig` so that outbound transactions stay within the
  limits enforced by the Synnergy virtual machine.
* **Consensus Awareness** – Hooks into `internal/config.ConsensusConfig` to
  delay settlement until finality thresholds are met, preventing forks from
  leaking into bridged chains.
* **Wallet and CLI Integration** – Provides deterministic APIs for the CLI
  (`synnergy bridge send`) and the JavaScript dashboard to display proof
  progress, audit log references and real-time throughput metrics.
* **Security Controls** – Applies policy from `internal/security` to enforce
  permissioned access, double-spend detection and encrypted payload routing
  between authority nodes.

## Workflow Overview

1. A user submits a cross-chain command through the CLI or web UI.
2. The request is authenticated via the governance audit log and queued with a
   deterministic replay-protected identifier.
3. Bridge relayers collect signatures from authority nodes using the wallet
   module’s multi-signature configuration.
4. Finalised proofs are broadcast to the destination chain and the resulting
   receipts are streamed back to telemetry for operator visibility.

This module ties together the consensus, VM, wallet and telemetry layers to
deliver fault-tolerant interoperability suitable for regulated enterprise
deployments.
