# How to Apply for a Grant or Loan from the LoanPool

## 1. Overview
The LoanPool is Neto Solaris's on-chain credit facility for the Synnergy Network. It allows our community to propose, review and disburse grants or loans from a shared treasury in a transparent and auditable manner. Two complementary flows exist:

- **Proposals** – Authority nodes or authorised users submit proposals to direct treasury funds to a recipient. These can represent grants or structured loans.
- **Applications** – Individual users submit personal loan requests that are voted on before disbursement.

Both flows rely on the same underlying treasury and governance controls, ensuring funds are allocated fairly and within available balances.

## 2. Eligibility and Preparation
Before applying, ensure that you:

1. Hold a verified **SYN900 identity token**, registered through the IdentityService, and have completed all KYC/AML checks.
2. Possess a Synnergy wallet with any required collateral and have enabled multi-factor authentication.
3. Have gathered supporting documentation describing the purpose of the grant or loan, repayment expectations, and projected impact.
4. Understand whether your request is a **grant** (no repayment) or a **loan** (with repayment terms). This distinction is provided in the `type` field when submitting a proposal.
5. Ensure your identity verification logs remain up to date for audit and regulatory review.

## 3. Submitting a Proposal (Grants or Directed Loans)
Authority nodes and approved contributors submit proposals using the `loanpool` CLI module:

```bash
synnergy loanpool submit [creator] [recipient] [type] [amount] [description]
```

- **creator** – wallet address submitting the proposal.
- **recipient** – address that will receive funds if approved.
- **type** – `grant`, `loan`, or other categories accepted by your governance policy.
- **amount** – denomination in smallest unit of the treasury currency.
- **description** – human‑readable justification for the request.
- Proposals remain open for a 24‑hour voting window by default; use `loanpool extend [creator] [id] [hrs]` to add time or `loanpool cancel [creator] [id]` to withdraw the request.

A helper script (`cmd/scripts/loanpool_apply.sh`) wraps this command and provides sensible defaults. After submission, participants can cast votes:

```bash
synnergy loanpool vote [voter] [id]
```

Each proposal accrues votes until the deadline, after which the pool can be ticked to evaluate approvals:

```bash
synnergy loanpool tick
```

The `tick` operation scans all pending proposals and marks them approved when at least one vote exists and the deadline has not yet passed.

Approved proposals are disbursed from treasury if sufficient balance remains:

```bash
synnergy loanpool disburse [id]
```

Disbursements draw from the shared treasury and will error if insufficient funds remain.

## 4. Submitting a Personal Loan Application
Individuals request financing directly through the `loanpool_apply` module, which manages loan applications on top of the same treasury:

```bash
synnergy loanpool_apply submit [applicant] [amount] [termMonths] [purpose]
```

- **applicant** – verified user ID.
- **amount** – requested funds.
- **termMonths** – repayment duration in months.
- **purpose** – reason for the loan.

Community members vote with:

```bash
synnergy loanpool_apply vote [voter] [id]
```

An administrator or scheduled process finalises approvals:

```bash
synnergy loanpool_apply process
```

During processing, any application with at least one recorded vote is marked as approved.

Once approved, funds are released to the applicant:

```bash
synnergy loanpool_apply disburse [id]
```

## 5. Loan and Grant Programme Types
The LoanPool supports branded funding instruments tailored for distinct use cases:

- **Repayable proposals** – capital-intensive projects are submitted as `LoanProposal` records and must reach community quorum before disbursement.
- **Micro-loans** – streamlined retail requests flow through `LoanPoolApply`, where a single affirmative vote is sufficient for approval【F:core/loanpool_apply.go†L33-L64】.
- **Non-repayable grants** – the SYN3800 registry records each grant with a beneficiary and named category, enabling programme-level accounting【F:core/syn3800.go†L8-L37】. Tests illustrate a canonical `research` grant example【F:core/syn3800_test.go†L5-L8】.

Named smart-contract templates provide additional loan flavours:

- **CollateralizedLoan** – secures borrowing against pledged assets before funds release【F:smart-contracts/solidity/CollateralizedLoan.sol†L4-L5】.
- **UndercollateralizedLoan** – allows credit-based lending with reduced collateral requirements【F:smart-contracts/solidity/UndercollateralizedLoan.sol†L4-L5】.
- **StableRateLoan** – fixes interest charges for predictable repayment schedules【F:smart-contracts/solidity/StableRateLoan.sol†L4-L5】.
- **VariableRateLoan** – adjusts interest dynamically to governance or oracle benchmarks【F:smart-contracts/solidity/VariableRateLoan.sol†L4-L5】.
- **FlashLoan** – enables uncollateralised borrowing within a single transaction cycle【F:smart-contracts/solidity/FlashLoan.sol†L4-L5】.
- **LoanAuction** – framework for market-driven rate discovery among competing lenders【F:smart-contracts/solidity/LoanAuction.sol†L4-L5】.
- **LoanFactory** and **LoanRegistry** – scaffolds for deploying standardised loans and cataloguing active positions【F:smart-contracts/solidity/LoanFactory.sol†L4-L5】【F:smart-contracts/solidity/LoanRegistry.sol†L4-L5】.

Grant programmes are likewise extensible:

- **GrantTracker** – foundation for tracking grant milestones and spending【F:smart-contracts/solidity/GrantTracker.sol†L3-L5】.
- **GrantMatching** – pairs community donations with approved proposals for matched-fund campaigns【F:smart-contracts/solidity/GrantMatching.sol†L4-L5】.
- **ScholarshipFund** – dedicated channel for education-oriented grants and scholarships【F:smart-contracts/solidity/ScholarshipFund.sol†L4-L5】.

## 6. Authority Oversight and Distribution
Disbursement authority is restricted to vetted governance participants:

- The **AuthorityNodeRegistry** registers nodes and tracks votes, exposing membership checks through `IsAuthorityNode`【F:core/authority_nodes.go†L34-L55】【F:core/authority_nodes.go†L98-L103】.
- The virtual machine provides a `Loanpool_IsAuthority` opcode so contracts and CLIs can verify that the caller is a registered authority before releasing funds【F:core/opcode.go†L603-L610】.
- Authority nodes are managed through the `synnergy authority` CLI, which supports registration, voting and electorate sampling for due diligence【F:cli/authority_nodes.go†L15-L48】.

Only nodes passing these checks may invoke `loanpool disburse` or `loanpool_apply disburse`, ensuring that loan and grant distribution remains under explicit Neto Solaris oversight.

## 7. Tracking and Management
Both proposals and applications can be inspected and listed for auditability:

```bash
synnergy loanpool get [id]
synnergy loanpool list
synnergy loanpool_apply get [id]
synnergy loanpool_apply list
```

Lifecycle adjustments for proposals are also available:

```bash
synnergy loanpool cancel [creator] [id]
synnergy loanpool extend [creator] [id] [hrs]
```

Administrative operators at Neto Solaris can pause new submissions or review treasury statistics through the `loanmgr` commands:

```bash
synnergy loanmgr pause
synnergy loanmgr resume
synnergy loanmgr stats
```

These controls interact with the underlying `LoanPool` manager, which reports `Treasury`, `ProposalCount`, `ApprovedCount` and `DisbursedCount` metrics.

## 8. Enterprise Integration and Automation
Structured outputs make the LoanPool suitable for enterprise monitoring and automation:

- `loanpool` and `loanpool_apply` `get`/`list` commands return JSON `LoanProposalView` and `LoanApplicationView` records, enabling dashboards and compliance archives.
- `loanmgr stats` surfaces aggregate treasury data for capacity planning.
- Machine‑readable responses allow orchestration tools to trigger disbursements or workflow steps programmatically.

## 9. Compliance and Security
Neto Solaris enforces strict security around LoanPool operations:

- **Identity verification** – All participants must authenticate via SYN900 tokens with KYC/AML compliance, with verification methods logged for audit.
- **Multi‑factor authentication** – Sensitive actions require multiple verification factors, including optional biometric checks.
- **AI‑driven monitoring** – Machine‑learning models scan for voting anomalies and fraudulent submissions, alerting operators to suspicious behaviour.
- **Immutable records** – Every proposal, vote and disbursement is recorded on-chain for full audit trails.

## 10. Best Practices
- Provide concise and complete descriptions to speed up review.
- Monitor deadlines to ensure proposals receive sufficient votes before expiry.
- For loans, maintain repayment discipline to remain eligible for future financing.
- Use the `get` and `list` commands regularly to track status and maintain transparency with stakeholders.
- Review `loanmgr stats` periodically and cancel or extend proposals proactively when circumstances evolve.

## Stage 78 Enterprise Enhancements
- **LoanPool diagnostics:** The enterprise orchestrator surfaces proposal counts, consensus relayers and wallet status so grant officers can confirm availability via `synnergy orchestrator status --json` or the function web before announcing new funding rounds.【F:core/enterprise_orchestrator.go†L21-L166】【F:cli/orchestrator.go†L6-L75】
- **Transparent fees:** Stage 78 gas documentation covers orchestrator operations that underpin loan approvals, authority elections and wallet sealing, keeping applicant-facing costs stable across CLI, VM and browser experiences.【F:docs/reference/gas_table_list.md†L420-L424】
- **Reliability testing:** Unit, situational, stress, functional and real-world orchestrator suites validate LoanPool integrations under heavy submission volumes and regulatory audits, protecting applicants from downtime during grant windows.【F:core/enterprise_orchestrator_test.go†L5-L75】【F:cli/orchestrator_test.go†L5-L26】

## 11. Support
For additional assistance or to escalate issues, contact the Neto Solaris support team through your authority node representative or the official Synnergy Network communication channels.

By following these steps and aligning with our compliance standards, applicants can seamlessly access grants or loans from the LoanPool while contributing to the health and growth of the Synnergy ecosystem.
