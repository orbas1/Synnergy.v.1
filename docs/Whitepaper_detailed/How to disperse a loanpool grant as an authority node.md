# How to Disperse a LoanPool Grant as an Authority Node

## 1. Purpose
Authority nodes are entrusted by **Neto Solaris** to release treasury funds only after a proposal or application has satisfied governance criteria. This guide walks through the review and disbursement process so funds are transferred securely and auditable on the Synnergy Network.

## 2. Architecture Overview
The LoanPool stack is composed of tightly scoped components that together enforce governance:

- **LoanPool** – core ledger that tracks proposals, approvals and disbursements【F:core/loanpool.go†L9-L15】.
- **LoanPoolManager** – administrative facade exposing treasury statistics and the ability to pause or resume submissions【F:core/loanpool_management.go†L21-L29】.
- **LoanPoolApply** – simplified application layer sharing the same treasury for retail requests【F:core/loanpool_apply.go†L17-L23】.
- **AuthorityNodeRegistry** – validates node membership, restricting sensitive operations to registered addresses【F:core/authority_nodes.go†L98-L103】.

## 3. Prerequisites
Before attempting any payout, the operator must:

1. **Hold authority status** – membership is validated against the on‑chain registry through `IsAuthorityNode`, ensuring only registered addresses can release funds【F:core/authority_nodes.go†L98-L103】.
2. **Have network access to the LoanPool CLI**, hardened credentials and an auditable signing environment.
3. **Verify treasury capacity** – administrators can inspect current balances and counts with `LoanPoolManager.Stats` before disbursement【F:core/loanpool_management.go†L31-L44】.
4. **Confirm pool availability** – use `loanmgr pause`/`resume` to control proposal intake during maintenance windows【F:cli/loanpool_management.go†L19-L29】.

## 4. Review Proposal Status
Disbursement is only permitted for proposals that have been approved and not yet paid. Run the proposal processing step to update approval flags:

```bash
synnergy loanpool tick
```

The `Tick` routine scans active proposals and marks any with at least one vote as approved prior to their deadline【F:core/loanpool.go†L50-L56】.

List proposals and inspect the specific request:

```bash
synnergy loanpool list
synnergy loanpool get [id]
```

Confirm that `approved` is `true` and `disbursed` is `false` in the returned JSON view before continuing.

## 5. Disbursement Workflow
Once approval is confirmed:

1. **Execute the disbursement command**:

   ```bash
   synnergy loanpool disburse [id]
   ```

   The CLI delegates to the `Disburse` method which validates approval, checks that the proposal has not previously been paid, confirms sufficient treasury, and then deducts the amount from the pool【F:cli/loanpool.go†L51-L58】【F:core/loanpool.go†L60-L74】.  
   Validate the payout by re-running `loanmgr stats` to see updated treasury figures【F:cli/loanpool_management.go†L31-L37】.

2. **Record the transaction** in internal ledgers and compliance systems for future audits.

### Disbursing Loan Applications
Retail applications managed via `loanpool_apply` follow an analogous flow:

```bash
synnergy loanpool_apply disburse [id]
```

This command triggers the application disbursement routine which performs the same approval and treasury checks before marking the application as paid【F:cli/loanpool_apply.go†L53-L60】【F:core/loanpool_apply.go†L67-L81】.

## 6. Administrative Safeguards
Enterprise operators should employ the following controls:

- **Pause or resume submissions** to manage incident response or system upgrades via `loanmgr pause` and `loanmgr resume`【F:cli/loanpool_management.go†L19-L29】.
- **Cancel or extend proposals** when creators request changes or need more time, using `loanpool cancel` or `loanpool extend`【F:cli/loanpool.go†L97-L120】【F:core/loanpool.go†L93-L117】.
- **List outstanding work** with `loanpool list` to maintain visibility across proposals【F:cli/loanpool.go†L87-L95】.

## 7. Post-Disbursement & Audit

- **Verify updated pool metrics** using `LoanPoolManager.Stats` to ensure treasury balances and disbursement counters reflect the transaction【F:core/loanpool_management.go†L31-L44】.
- **Archive proposal views** for governance records using `synnergy loanpool get [id]` after the payout; the view serialises all vote counts and flags for retention【F:core/loanpool_views.go†L5-L35】.
- **Feed audit trails** into internal compliance platforms and monitoring dashboards for continuous oversight.
- **Monitor remaining treasury** for future grants to maintain liquidity targets set by Neto Solaris

## 8. Troubleshooting

Common errors returned by the disbursement routines include:

- `proposal not approved or already disbursed` – the proposal is either pending approval or has been paid already.
- `insufficient treasury` – pool balance is below the requested amount.
- `loan pool is paused` – new proposals are blocked until `loanmgr resume` is executed【F:core/loanpool.go†L26-L33】【F:core/loanpool_management.go†L21-L28】.
- `application not approved or already disbursed` – analogous checks for loan applications.

Re-run `tick` to process approvals, verify authority status, or replenish treasury before retrying.

## 9. Security and Governance Notes

- Only authority nodes validated by Neto Solaris may invoke disbursement routines. Attempting to disburse from an unregistered address will fail the `IsAuthorityNode` check【F:core/authority_nodes.go†L98-L103】.
- Execute all commands over secure channels with hardware-backed keys, and replicate logs to the enterprise monitoring stack.
- Use administrative pause and resume controls to isolate incidents or coordinate multi-region upgrades【F:cli/loanpool_management.go†L19-L29】.

By following these steps, authority nodes ensure that LoanPool grants and loans are released in a controlled, transparent, and compliant manner consistent with Neto Solaris governance standards.

