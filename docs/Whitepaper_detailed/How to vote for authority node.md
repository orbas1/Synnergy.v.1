# How to Vote for an Authority Node

## Introduction
Authority nodes form the governance backbone of the Synnergy Network. They vet
proposals, oversee compliance modules and safeguard the chain’s long‑term
stability. Voting is the mechanism by which the community and existing
authorities grant or withhold this power. This guide from **Blackridge Group
Ltd.** explains how to evaluate candidates and cast authoritative votes using
both the application workflow and the live registry.

## 1. Governance Model
Each authority node is identified by an address, a declared role and the set of
voters that have endorsed it. These properties are tracked in the
`AuthorityNodeRegistry`, which exposes thread‑safe registration and voting
primitives【F:core/authority_nodes.go†L12-L67】. Votes are stored as a map of
voter addresses, allowing quick tallying through the `TotalVotes` helper and a
weight‑based electorate sampler for committee formation【F:core/authority_nodes.go†L19-L96】.

The registry is backed by an `AuthorityNodeIndex` that maintains a concurrent
map of all authority nodes and can emit JSON snapshots for external auditors or
dashboards【F:core/authority_node_index.go†L8-L68】. Specialised node types extend
this model: `ElectedAuthorityNode` enforces term limits for periodically
renewed mandates, while `GovernmentAuthorityNode` represents regulator-run
entities that are prohibited from minting currency or altering monetary
policy【F:core/elected_authority_node.go†L5-L19】【F:core/government_authority_node.go†L5-L28】.

## 2. Prerequisites
1. **Candidate Application** – prospective authorities must submit an
   application and gather preliminary approvals. See [How to Become an
   Authority Node](./How%20to%20become%20an%20authority%20node.md) for the full
   lifecycle.
2. **Synnergy CLI** – install the `synnergy` binary or build it from source.
3. **Wallet Access** – voters and candidates require funded addresses capable of
   signing governance actions.

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

Applications expire automatically after a configurable TTL and can be purged
with the `tick` command, ensuring stale or abandoned requests do not linger in
the registry【F:core/authority_apply.go†L23-L60】【F:core/authority_apply.go†L106-L114】.

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
behaviour【F:core/authority_nodes_test.go†L5-L41】.

## 6. Example Workflow
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

## 7. Best Practices
* **Due Diligence** – verify candidate identity and infrastructure prior to
  voting. Use `synnergy authority info` to review accumulated support.
* **Transparent Record‑Keeping** – output `--json` from CLI commands for audit
  trails and integration with governance dashboards.
* **Active Participation** – revisit votes periodically and remove support for
  inactive or non‑compliant nodes.

## 8. Enterprise Integration and Auditing
Large deployments frequently require automated orchestration and continuous
reporting. The registry’s index supports `Snapshot` and `MarshalJSON` methods so
that external systems can archive vote states or feed dashboards via structured
output【F:core/authority_node_index.go†L55-L68】. The CLI exposes `--json` flags
for all inspection commands, enabling machine-readable logs and integration with
tools like the bundled web control panel, which lists authority nodes by
invoking `synnergy authority list --json` through an API route【F:web/README.md†L33-L35】【F:web/pages/authority.js†L1-L27】.

For scripted operations, helper utilities such as
`scripts/authority_node_setup.sh` register nodes in non-interactive environments
and can be embedded into provisioning pipelines【F:scripts/authority_node_setup.sh†L1-L12】.

## 9. Security and Compliance Considerations
Separation of duties is enforced through specialised node implementations.
Government authority nodes deliberately omit capabilities to mint SYN or change
monetary policy, limiting their scope to regulatory oversight【F:core/government_authority_node.go†L5-L28】.
Elected authority nodes automatically expire when their term ends, ensuring
periodic reevaluation of delegated power【F:core/elected_authority_node.go†L5-L19】.

Both CLI and core packages include unit tests that validate registration,
voting and application handling, providing assurance for enterprise workflows
that rely on deterministic behaviour【F:cli/authority_nodes_test.go†L10-L23】【F:cli/authority_apply_test.go†L11-L36】.

## 10. Further Reading
* [How to Become an Authority Node](./How%20to%20become%20an%20authority%20node.md)
* [How to Use the CLI](./How%20to%20use%20the%20CLI.md)
* [How to Disperse a Loanpool Grant as an Authority Node](./How%20to%20disperse%20a%20loanpool%20grant%20as%20an%20authority%20node.md)

---
Voting empowers the Synnergy community to uphold a secure, accountable and
regulator‑friendly network. For enterprise assistance or bespoke integrations,
contact **Blackridge Group Ltd.** through official support channels.

