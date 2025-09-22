# How to Vote for an Authority Node

## Introduction
Authority nodes form the governance backbone of the Synnergy Network. They vet
proposals, oversee compliance modules and safeguard the chain’s long‑term
stability. Voting is the mechanism by which the community and existing
authorities grant or withhold this power. This guide from **Neto Solaris** explains how to evaluate candidates and cast authoritative votes using
both the application workflow and the live registry.

## 1. Governance Model
Each authority node is identified by an address, a declared role and the set of
voters that have endorsed it. These properties are tracked in the
`AuthorityNodeRegistry`, which exposes thread‑safe registration and voting
primitives【F:core/authority_nodes.go†L12-L67】. Votes are stored as a map of
voter addresses, allowing quick tallying through the `TotalVotes` helper and a
weight‑based electorate sampler for committee formation【F:core/authority_nodes.go†L19-L96】.

## 2. Prerequisites
1. **Candidate Application** – prospective authorities must submit an
   application and gather preliminary approvals. See [How to Become an
   Authority Node](./How%20to%20become%20an%20authority%20node.md) for the full
   lifecycle.
2. **Synnergy CLI** – install the `synnergy` binary or build it from source.
3. **Wallet Access** – voters and candidates require funded addresses capable of
   signing governance actions.

Stage 79 bootstrap automation lets governance teams enrol authority roles and ledger replication ahead of ballots. Execute `synnergy orchestrator bootstrap --authority council=governor --replicate` to call `core.EnterpriseOrchestrator.BootstrapNetwork`, which signs the bootstrap, registers the authority mapping, enables replication and reports diagnostics so voters can confirm quorum and infrastructure readiness before casting votes.【F:cli/orchestrator.go†L58-L117】【F:core/enterprise_orchestrator.go†L71-L209】【F:core/enterprise_orchestrator_test.go†L73-L178】 Startup also synchronises Stage 79 gas costs with documentation, and the control panel offers the same workflow to keep browser operators and CLI automation aligned on pricing and state.【F:cmd/synnergy/main.go†L63-L106】【F:web/pages/index.js†L1-L214】【F:web/pages/api/bootstrap.js†L1-L45】

## 3. Voting on Applications
During the application phase, votes determine whether a candidate is admitted to
the registry.

```bash
synnergy authority_apply list --json      # view pending applications
synnergy authority_apply vote <voter> <appID> <true|false>
```

The `authority_apply` module records approval and rejection sets for each
submission and validates inputs before updating the tally【F:cli/authority_apply.go†L21-L41】.
Finalisation promotes approved candidates to registered authority nodes:

```bash
synnergy authority_apply finalize <appID>
```

Applications carry a time‑to‑live and are garbage‑collected if not finalized.
The manager timestamps each submission and prunes expired or already processed
entries via `Tick`, ensuring stale requests cannot linger in the queue【F:core/authority_apply.go†L23-L30】【F:core/authority_apply.go†L106-L114】.

## 4. Voting for Registered Nodes
Once a node is registered, community members may express ongoing support or
participate in committee selection using the registry commands.

```bash
synnergy authority vote <voter> <candidate>      # cast a vote
synnergy authority info <candidate>              # inspect vote count
synnergy authority elect 3                      # sample top 3 electorate
```

These commands delegate to the registry’s `Vote`, `Info` and `Electorate`
methods to update vote maps, retrieve node metadata and assemble
vote‑weighted committees【F:cli/authority_nodes.go†L33-L52】【F:cli/authority_nodes.go†L57-L88】.

## 5. Managing Votes
Voters may withdraw support by removing their ballot, and operators can verify
membership status or list all authorities for audit purposes.

```bash
synnergy authority vote <voter> <candidate>      # cast or update a vote
synnergy authority list --json                   # inspect registry
synnergy authority deregister <candidate>        # remove node (admin only)
```

While the CLI focuses on casting votes, the registry also exposes a
`RemoveVote` method for custom tooling to withdraw support when necessary.
Registry queries such as `IsAuthorityNode` and `List` enable health checks and
automation scripts【F:core/authority_nodes.go†L69-L129】. Unit tests exercise
vote casting, electorate sampling and vote removal to guarantee deterministic
behaviour【F:core/authority_nodes_test.go†L5-L36】.

## 6. Audit and Automation
The registry is backed by an `AuthorityNodeIndex` that supports snapshotting
for consistent reads and deterministic JSON output for downstream services.
Operators can safely iterate over a `Snapshot` or export the index through
`MarshalJSON` for integration with monitoring systems and dashboards【F:core/authority_node_index.go†L55-L67】.

CLI commands mirror this machine‑readable philosophy. Both `authority` and
`authority_apply` expose `--json` flags on `info`, `list` and `get` commands
to facilitate scripted audits and GUI consumption【F:cli/authority_nodes.go†L54-L92】【F:cli/authority_apply.go†L54-L93】.

## 7. Advanced Governance Roles
Beyond simple registrations, the codebase provides specialised authority node
constructs. `ElectedAuthorityNode` embeds a term limit and includes `IsActive`
to verify whether its mandate has expired【F:core/elected_authority_node.go†L5-L19】.
`GovernmentAuthorityNode` models regulator‑operated validators and hard‑codes
restrictions against minting SYN coins or altering monetary policy, enforcing
separation of duties in regulated environments【F:core/government_authority_node.go†L5-L27】.

## 8. Example Workflow
1. **Submit and Review Application**
   - Candidate runs `synnergy authority_apply submit <addr> <role> "description"`.
   - Voters inspect with `synnergy authority_apply list`.
2. **Cast Votes**
   - Supporters use `synnergy authority_apply vote <voter> <id> true`.
   - Opponents supply `false` to reject.
3. **Finalize and Register**
   - When quorum is met, run `synnergy authority_apply finalize <id>`.
   - Successful candidates appear in `synnergy authority list`.
4. **Ongoing Support**
   - Cast registry votes via `synnergy authority vote`.
   - Use `synnergy authority elect <n>` to draw committee members for special
     tasks.

## 9. Best Practices
* **Due Diligence** – verify candidate identity and infrastructure prior to
  voting. Use `synnergy authority info` to review accumulated support.
* **Transparent Record‑Keeping** – output `--json` from CLI commands for audit
  trails and integration with governance dashboards.
* **Active Participation** – revisit votes periodically and remove support for
  inactive or non‑compliant nodes.

## 10. Further Reading
* [How to Become an Authority Node](./How%20to%20become%20an%20authority%20node.md)
* [How to Use the CLI](./How%20to%20use%20the%20CLI.md)
* [How to Disperse a Loanpool Grant as an Authority Node](./How%20to%20disperse%20a%20loanpool%20grant%20as%20an%20authority%20node.md)

---
Voting empowers the Synnergy community to uphold a secure, accountable and
regulator‑friendly network. For enterprise assistance or bespoke integrations,
contact **Neto Solaris** through official support channels.

