# How to Become an Authority Node

Becoming an authority node in the Synnergy Network is a formal on‑chain
process designed to preserve the integrity, security and continuity of the
ecosystem.  This guide from **Neto Solaris** explains the
prerequisites, application workflow and ongoing expectations for
participants who wish to obtain and maintain authority status.

## 1. Role of Authority Nodes

Authority nodes are trusted entities that participate in governance
activities, manage critical services and help steer protocol evolution.  In
code, an authority node is defined by its address, role and recorded votes
from the community, all stored in an indexed registry that enables rapid
lookups and election sampling【F:core/authority_nodes.go†L11-L80】.

At scale, this registry leverages a concurrent in‑memory index that stores
addresses, roles and vote maps, allowing constant‑time queries even in large
networks【F:internal/nodes/authority_nodes/index.go†L12-L56】. Operators interact
with the registry through the `authority` CLI, which exposes commands for
registration, voting, electorate sampling and deregistration【F:cli/authority_nodes.go†L15-L112】.

### Specialised Variants

The platform recognises specialised authority implementations, such as:

* **Elected Authority Nodes** – nodes granted authority for a fixed term,
  after which their status expires unless re‑elected【F:core/elected_authority_node.go†L5-L17】.
* **Government Authority Nodes** – regulator‑operated nodes that cannot
  mint coins or alter monetary policy, preserving the network's fixed supply
  guarantees【F:core/government_authority_node.go†L1-L24】.

## 2. Prerequisites

Prospective operators should prepare the following before applying:

1. **Wallet & Address** – generate and secure a Synnergy wallet address with
   sufficient funds for staking and operational costs.
2. **Infrastructure** – provision reliable compute, storage and network
   resources.  The `network.yaml` configuration illustrates typical port and
   peer settings used in production deployments【F:configs/network.yaml†L1-L38】.
3. **Synnergy CLI** – install the `synnergy` command‑line interface or build
   it from source using the project’s tooling.
4. **Identity & KYC** – register operator information through the built‑in
   identity service and record a hashed KYC commitment so that governance
   actions can be audited【F:identity_verification.go†L9-L58】【F:compliance.go†L12-L70】.
5. **Regulatory Policies** – review Neto Solaris’s operational standards
   and confirm that any jurisdiction‑specific regulations are encoded in the
   network’s policy manager before launch【F:regulatory_management.go†L8-L38】.

## 3. Application Lifecycle

Authority status is granted through an application and voting process
administered on‑chain via the `AuthorityApplicationManager`【F:core/authority_apply.go†L9-L110】.

### 3.1 Submit an Application

```bash
synnergy authority_apply submit <candidate_address> <role> "<description>"
```

This command creates an application with a defined time‑to‑live.  Each
submission receives a unique identifier that is used for subsequent voting
and finalisation.

### 3.2 Gather Votes

Community members and existing authorities review the proposal and vote:

```bash
synnergy authority_apply vote <voter_address> <application_id> <true|false>
```

Votes can be amended or withdrawn during the application window.  Multiple
approval votes are required to advance a candidate.

### 3.3 Finalise the Application

When the voting period concludes, finalise the outcome:

```bash
synnergy authority_apply finalize <application_id>
```

If approvals exceed rejections, the candidate is registered automatically in
the authority node registry.  Expired or finalised applications are cleaned
up by periodic `tick` operations. The `authority` CLI can then be used to
inspect membership, sample electorates or deregister nodes as required
【F:cli/authority_nodes.go†L20-L112】.

### 3.4 Convenience Scripts

For automated environments, the repository provides helper scripts:

* `cmd/scripts/authority_apply.sh` – submits an application with default
  parameters for rapid testing【F:cmd/scripts/authority_apply.sh†L1-L8】.
* `scripts/authority_node_setup.sh` – registers an already approved node in
  the registry using the CLI【F:scripts/authority_node_setup.sh†L1-L12】.

## 4. Post‑Approval Responsibilities

Once registered, authority nodes are expected to:

* Participate actively in governance by reviewing and voting on future
  authority applications and protocol proposals.
* Maintain high availability and security standards, including timely patch
  application and continuous monitoring.
* Operate additional services relevant to their role (e.g., regulatory,
  auditing or cross‑chain oversight).
* Adhere to financial controls; for example, government authority nodes are
  restricted from minting or altering monetary policy by design.

The registry and accompanying CLI commands allow operators to inspect their
status, sample electorates or deregister when necessary.

## 5. Network Configuration

Authority nodes must be listed in the network configuration so peers can
identify them.  Each entry specifies the node’s identifier, type and public
address as seen in the sample configuration file【F:configs/network.yaml†L42-L52】.
Genesis files may also allocate initial balances or stakes to authority
addresses to bootstrap participation【F:configs/genesis.json†L5-L32】.

## 6. Security and Compliance

Authority nodes are bound by strict auditability requirements:

* **Identity Service** – stores operator metadata and verification logs so
  that actions can be traced back to a real‑world entity【F:identity_verification.go†L9-L58】.
* **Compliance Service** – records hashed KYC commitments, risk scores and
  audit entries, enabling transparent oversight of node behaviour【F:compliance.go†L12-L70】.
* **Regulatory Manager** – encodes jurisdictional policies that transactions
  must satisfy; regulatory nodes can flag or reject non‑compliant activity
  based on these rules【F:regulatory_management.go†L8-L38】【F:regulatory_node.go†L8-L33】.

## 7. Operational Tooling

Enterprise operators often automate routine tasks:

* **System Health Logging** – the `SystemHealthLogger` gathers runtime metrics
  for dashboards and alerting【F:system_health_logging.go†L11-L47】.
* **Failover Management** – `FailoverManager` promotes standby nodes when the
  primary becomes unresponsive, supporting seamless upgrades and recovery
  processes【F:high_availability.go†L8-L69】.
* **Dashboards** – the web control panel exposes an `/authority` page that
  lists registered nodes via the CLI, while the Authority Node Index GUI
  offers a TypeScript foundation for rich enterprise interfaces【F:web/README.md†L33-L35】【F:web/pages/authority.js†L1-L30】【F:GUI/authority-node-index/README.md†L1-L33】.

## 8. Best Practices

* **Security Hardening** – isolate keys, enable firewalls and monitor logs.
* **Observability** – use the system health logging utilities to capture
  metrics and detect anomalies【F:system_health_logging.go†L11-L47】.
* **High Availability** – configure failover managers to keep authority
  services online during maintenance or outages【F:high_availability.go†L8-L69】.
* **Regular Upgrades** – track releases from Neto Solaris and apply
  upgrades in a staged manner to minimise downtime.
* **Community Engagement** – remain responsive to governance discussions and
  contribute to ecosystem documentation and tooling.

## 9. Further Reading

* [How to Vote for an Authority Node](./How%20to%20vote%20for%20authority%20node.md)
* [How to Disperse a Loanpool Grant as an Authority Node](./How%20to%20disperse%20a%20loanpool%20grant%20as%20an%20authority%20node.md)
* [Authority Node Index GUI](../architecture/node_operations_dashboard_architecture.md)

---

For questions or support, contact **Neto Solaris** through the
project’s official communication channels.  Operating an authority node is a
privileged role; diligence and transparency are essential to sustain the
Synnergy Network’s vision.

