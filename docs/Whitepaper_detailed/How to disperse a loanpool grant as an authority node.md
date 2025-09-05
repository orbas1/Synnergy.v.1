# How to Disperse a LoanPool Grant as an Authority Node

## 1. Purpose
Authority nodes are entrusted by **Blackridge Group Ltd** to release treasury funds only after a proposal or application has satisfied governance criteria. This guide walks through the review, approval, and disbursement process so funds are transferred securely and auditable on the Synnergy Network.

## 2. Prerequisites
Before attempting any payout, the operator must:

1. **Hold authority status** – membership is validated against the on-chain registry through `IsAuthorityNode`, ensuring only registered addresses can release funds【F:core/authority_nodes.go†L98-L103】.
2. **Have network access to the LoanPool CLI** and appropriate signing keys.
3. **Verify treasury capacity** – administrators can inspect current balances and counters with `LoanPoolManager.Stats` to ensure the pool is solvent before disbursement【F:core/loanpool_management.go†L31-L44】.

## 3. Pre-Disbursement Audit
Rigorous checks uphold Blackridge Group Ltd compliance standards.

### 3.1 Update Approvals
Process votes and deadlines so proposals reflect their latest status:
```bash
synnergy loanpool tick
```
`Tick` scans active proposals and marks any with at least one vote as approved before their deadline【F:core/loanpool.go†L50-L56】.

For retail applications run:
```bash
synnergy loanpool_apply process
```
which finalises applications with at least one vote【F:core/loanpool_apply.go†L58-L63】.

### 3.2 Inspect Proposals and Applications
List all pending items to review their approval and disbursement flags:
```bash
synnergy loanpool list
synnergy loanpool_apply list
```
The commands surface serialised views built by `ProposalViews` and `ApplicationViews`, exposing vote totals and approval status for audit trails【F:core/loanpool_views.go†L44-L52】【F:core/loanpool_views.go†L89-L96】.

Drill into a specific item:
```bash
synnergy loanpool get [id]
synnergy loanpool_apply get [id]
```
`ProposalInfo` and `ApplicationInfo` return detailed JSON views to confirm `approved: true` and `disbursed: false` before payout【F:core/loanpool_views.go†L35-L42】【F:core/loanpool_views.go†L80-L87】.

### 3.3 Confirm Treasury Health
Retrieve pool metrics to confirm sufficient balance and to document pre-payout state. Administrators can call `LoanPoolManager.Stats` to obtain treasury amount and proposal counters【F:core/loanpool_management.go†L31-L44】.

### 3.4 Validate Authority Credentials
Ensure the executing address remains registered. `IsAuthorityNode` prevents unregistered operators from triggering payouts【F:core/authority_nodes.go†L98-L103】.

## 4. Disbursement Workflow
After the audit checks pass, initiate the payout.

### 4.1 Grants via `loanpool`
```bash
synnergy loanpool disburse [id]
```
The CLI forwards the request to `LoanPool.Disburse`, which validates approval, ensures the proposal has not already been paid, checks treasury sufficiency, deducts the amount, and marks the proposal as disbursed【F:cli/loanpool.go†L51-L58】【F:core/loanpool.go†L60-L74】.

### 4.2 Retail Applications via `loanpool_apply`
```bash
synnergy loanpool_apply disburse [id]
```
This command invokes `LoanPoolApply.Disburse`. The routine repeats approval and treasury checks against the shared pool before marking the application as settled【F:cli/loanpool_apply.go†L53-L60】【F:core/loanpool_apply.go†L67-L81】.

### 4.3 Administrative Safeguards
During incidents or audits, administrators may freeze submissions using:
```go
LoanPoolManager.Pause()
```
Setting `Paused` blocks new proposals via a guard in `SubmitProposal`【F:core/loanpool_management.go†L21-L23】【F:core/loanpool.go†L27-L30】. Resume normal operations with `LoanPoolManager.Resume`【F:core/loanpool_management.go†L26-L28】.

## 5. Post-Disbursement Duties
- **Verify updated pool metrics** using `LoanPoolManager.Stats` to ensure treasury balances and disbursement counters reflect the transaction【F:core/loanpool_management.go†L31-L44】.
- **Archive proposal or application views** for governance records using the relevant `get` command after payout.
- **Monitor remaining treasury** for future grants to maintain liquidity targets set by Blackridge Group Ltd.
- **Log cross-reference IDs** in enterprise compliance systems for traceability.

## 6. Troubleshooting
Typical errors include:
- `proposal not approved or already disbursed` – the proposal is either pending approval or has been paid already.
- `insufficient treasury` – pool balance is below the requested amount.
- `application not approved or already disbursed` – analogous checks for loan applications.
- `loan pool is paused` – submissions are blocked when administrators have paused the pool【F:core/loanpool.go†L27-L30】.

Re-run approval processing (`tick` or `process`), validate authority status, or replenish treasury before retrying.

## 7. Security and Governance Notes
- Only authority nodes validated by Blackridge Group Ltd may invoke disbursement routines; unregistered addresses fail the `IsAuthorityNode` check【F:core/authority_nodes.go†L98-L103】.
- Execute commands over secure channels with hardened key management.
- Use `Pause` during investigations or maintenance to prevent new proposals, and `Resume` once cleared【F:core/loanpool_management.go†L21-L28】.

By following these steps, authority nodes ensure that LoanPool grants and loans are released in a controlled, transparent, and compliant manner consistent with Blackridge Group Ltd governance standards.
