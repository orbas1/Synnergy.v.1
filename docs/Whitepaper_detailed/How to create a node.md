# How to Create a Node

## Introduction
Neto Solaris designs the Synnergy Network to be a modular, secure, and sustainable blockchain platform. Nodes are the operational backbone of this network, processing transactions, safeguarding data, and enabling advanced features like regulatory oversight and energy-aware consensus. This guide explains how to create and deploy a node that aligns with Synnergy’s architecture and governance principles.

## Node Archetypes
Synnergy provides a broad catalog of node implementations so operators can tailor infrastructure to their mission. The main families are outlined below.

### Core Infrastructure
- **Base Node** – provides foundational peer management and lifecycle control, exposing start/stop hooks and seed peer dialing【F:core/base_node.go†L10-L76】.
- **Full Node** – stores the complete chain in archive or pruned modes and validates blocks accordingly【F:core/full_node.go†L5-L38】.
- **Light Node** – maintains only block headers for resource‑constrained devices while still verifying chain continuity【F:core/light_node.go†L5-L31】.
- **Validator Node** – couples a ledger with validator and quorum management to orchestrate staking‑based consensus【F:core/validator_node.go†L5-L48】.
- **Gateway Node** – bridges external data sources through custom endpoints configurable at runtime. Endpoint registration is thread-safe and gated on the node's running state【F:core/gateway_node.go†L9-L68】.
- **Consensus‑Specific Node** – locks the consensus engine to Proof‑of‑Work, Proof‑of‑Stake, or Proof‑of‑History modes as required【F:core/consensus_specific_node.go†L1-L30】.
- **Elected Authority Node** – grants authority privileges for a fixed term, automatically expiring at the configured end date【F:core/elected_authority_node.go†L5-L19】.

### Operational and Data Nodes
- **Mining Node** – lightweight PoW miner that submits candidate blocks asynchronously, suitable for simulations or testnets【F:mining_node.go†L12-L76】.
- **Mobile Mining Node** – extends the mining node with battery‑aware throttling to protect mobile hardware【F:mobile_mining_node.go†L8-L44】.
- **Staking Node** – tracks delegated stakes and exposes locking and unlocking primitives without a full ledger dependency【F:staking_node.go†L5-L47】.
- **Indexing Node** – builds in‑memory key/value indexes for accelerated queries and analytics workloads【F:indexing_node.go†L5-L33】.
- **Content Node** – advertises available content items so peers can discover hosted assets across the network【F:content_node.go†L5-L33】.
- **Geospatial Node** – records geospatial coordinates and histories for tracked subjects【F:geospatial_node.go†L8-L46】.
- **Environmental Monitoring Node** – evaluates sensor feeds against programmable conditions to trigger on‑chain actions【F:environmental_monitoring_node.go†L9-L64】.

### Security, Oversight and Governance
- **Watchtower Node** – monitors peers, enforces firewall rules and streams system health metrics for anomaly detection【F:watchtower_node.go†L13-L66】.
- **Energy‑Efficient Node** – aggregates power consumption, issues sustainability certificates and throttles when efficiency drops【F:energy_efficient_node.go†L16-L80】.
- **Regulatory Node** – checks transactions against registered rules and logs violations for audit trails【F:regulatory_node.go†L8-L43】.
- **Custodial Node** – safekeeps user assets and releases funds back to the ledger with balance checks【F:core/custodial_node.go†L8-L48】.
- **Central Banking Node** – manages CBDC token supply under a declared monetary policy while respecting the fixed SYN cap【F:core/central_banking_node.go†L9-L37】.
- **Government Authority Node** – enforces governance without minting or policy powers, reflecting separation of duties【F:core/government_authority_node.go†L5-L27】.
- **Bank Institutional Node** – registers participating institutions and exposes their presence to the ledger【F:core/bank_institutional_node.go†L1-L43】.
- **Biometric Security Node** – gates privileged operations behind biometric verification routines【F:biometric_security_node.go†L8-L46】.
- **Audit Node** – couples bootstrap services with an audit manager to capture on‑chain events for compliance reports【F:core/audit_node.go†L8-L37】.
- **Forensic Node** – retains lightweight transaction and network traces for later investigation with capped buffers that prune old records to protect memory【F:core/forensic_node.go†L5-L56】.
- **Historical Node** – archives block summaries by height and hash to service long‑term retrieval requests【F:core/historical_node.go†L1-L55】.
- **Warfare Node** – tracks logistics and tactical updates for military assets, safeguarding command execution【F:warfare_node.go†L11-L48】.
- **Holographic Node** – distributes holographic frames for redundancy, enabling advanced data resilience experiments【F:internal/nodes/holographic_node.go†L9-L33】.

## Prerequisites
- **Operating system:** Linux or macOS with administrative access.
- **Dependencies:**
  - [Go](https://go.dev/) 1.20+ for building the node software.
  - Git for source management.
  - Optional: Docker and Kubernetes for containerised deployments.
- **Network:** Ability to open TCP ports 30303 (P2P) and 8080 (RPC) or adjust firewalls accordingly.

## Environment Setup
1. **Install Go and Git.** Ensure `go` is present in your `PATH`.
2. **Clone the repository.**
   ```bash
   git clone https://github.com/blackridge-group/Synnergy.v.1.git
   cd Synnergy.v.1
   ```
3. **Build the Synnergy binary.**
   ```bash
   go build -o bin/synnergy ./cmd/synnergy
   ```

## Creating a Node in Go
Instantiate the node type that suits your requirements. The following example creates a mining node with a nominal hash rate and starts it:

```go
package main
import (
    "synnergy"
)
func main() {
    miner := synnergy.NewMiningNode("miner-01", 100.0)
    miner.Start()
    defer miner.Stop()
}
```
The same pattern applies to other node types, using constructors such as `NewStakingNode`, `NewIndexingNode`, or `NewWatchtowerNode`.

## Network Registration via CLI
Use the `synnergy` command‑line interface to register and manage nodes. The authority workflow below illustrates submitting an application and activating the node【F:docs/Whitepaper_detailed/How to use the CLI.md†L7-L18】:
```bash
synnergy authority_apply submit node1 validator "candidate node"
synnergy authority_apply list --json
synnergy authority register node1 validator
synnergy authority is node1
```
Financial institutions can also register via banking commands before linking their nodes to the network【F:docs/Whitepaper_detailed/How to use the CLI.md†L14-L18】.

## Configuration and Deployment
For containerised environments, customise the provided Kubernetes manifest. It exposes both peer‑to‑peer and RPC interfaces and mounts a configuration file into the container【F:deploy/k8s/node.yaml†L1-L54】. Apply the manifest with `kubectl apply -f deploy/k8s/node.yaml` after adjusting resource limits and environment variables to your infrastructure.

## Security Hardening
Each node embeds a configurable firewall for blocking wallet addresses, token identifiers, or peer IPs, supporting real‑time rule updates without restarting【F:firewall.go†L5-L44】.

## Monitoring and Maintenance
System health metrics—including CPU usage, memory footprint, peer counts, and last block height—are collected through the `SystemHealthLogger` and can be exported to external monitoring stacks【F:system_health_logging.go†L11-L41】.

## High Availability and Resilience
The FailoverManager monitors node heartbeats and automatically promotes the most recent backup when the primary becomes unresponsive, minimising downtime for enterprise deployments【F:high_availability.go†L8-L69】.

## Audit, Forensics, and Archival
Audit, Forensic, and Historical nodes extend operational assurance: AuditNode records compliance events, ForensicNode captures transaction and network traces, and HistoricalNode serves archived block data for long‑term analytics【F:core/audit_node.go†L8-L37】【F:core/forensic_node.go†L5-L40】【F:core/historical_node.go†L1-L55】.

## Sustainability and Compliance
Energy‑efficient nodes can issue sustainability certificates and determine when throttling is necessary to maintain efficiency targets【F:energy_efficient_node.go†L52-L75】. For jurisdictions with oversight requirements, regulatory nodes evaluate transactions and log compliance violations for auditability【F:regulatory_node.go†L25-L43】.

## Conclusion
By following this guide, organisations can deploy nodes that integrate seamlessly with the Synnergy Network while meeting operational, environmental, and regulatory expectations. For enterprise support and advanced tooling, contact Neto Solaris

