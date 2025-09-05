# Loanpool

## Introduction
The Synnergy Network Loanpool is a branded treasury that allocates on‑chain capital to
community projects, small businesses, and infrastructure ventures. By embedding
compliance, voting, and disbursement logic directly in code, the Loanpool
provides transparent credit issuance and stimulates sustainable economic growth
across the ecosystem.

## Architecture Overview
Loanpool functionality is implemented in a set of core modules:

- **LoanPool** – maintains the pool’s treasury, proposal registry, and pause state
  for new submissions【F:core/loanpool.go†L9-L33】.
- **LoanProposal** – captures the details and voting state of a funding request【F:core/loanpool_proposal.go†L5-L17】.
- **LoanApplication** and **LoanPoolApply** – manage simplified loan applications
  backed by the same treasury【F:core/loanpool_apply.go†L5-L46】.
- **LoanPoolManager** – offers administrative controls such as pausing submissions
  and summarising pool statistics【F:core/loanpool_management.go†L21-L43】.
- **View structures** – serialisable representations of proposals and applications
  for external clients and CLI tooling【F:core/loanpool_views.go†L5-L52】【F:core/loanpool_views.go†L54-L97】.

## Treasury Mechanics
The pool is continuously capitalised through protocol fees. Five percent of every
transaction’s charges is routed into the Loanpool treasury, ensuring a growing
reserve for community credit initiatives【F:core/fees.go†L101-L123】. The canonical
treasury address is published in the VM’s opcode table so contracts can query the
current balance or perform on‑chain audits【F:core/opcode.go†L332-L339】.

## Proposal Lifecycle
The strategic flow of a Loanpool proposal is described in the Synnergy technical
whitepaper:

1. **Submission** – Applicants submit the recipient address, loan type, requested
   amount, collateral terms, and supporting documentation hashes.【F:Synnergy_Network_Future_Of_Blockchan.md†L245-L249】
2. **Review Period** – A minimum review window allows the community to evaluate
   proposals before voting begins.【F:Synnergy_Network_Future_Of_Blockchan.md†L250-L252】
3. **Voting Window** – Token holders vote with square‑root weighted power to avoid
   oligarchic dominance.【F:Synnergy_Network_Future_Of_Blockchan.md†L253-L276】
4. **Authority Ratification** – Designated authority nodes perform compliance and
   risk checks; their votes are binding and publicly documented.【F:Synnergy_Network_Future_Of_Blockchan.md†L255-L297】
5. **Disbursement** – Approved proposals release funds via timelocked contracts,
   locking collateral on‑chain and encoding milestone conditions.【F:Synnergy_Network_Future_Of_Blockchan.md†L258-L260】
6. **Repayment & Reporting** – Automated repayments update borrower credit scores
   and enforce penalties for delinquencies.【F:Synnergy_Network_Future_Of_Blockchan.md†L261-L263】

Within the codebase, proposals are voted on, evaluated, and disbursed through
methods such as `VoteProposal`, `Tick`, and `Disburse`, which ensure quorum
validation and treasury sufficiency before releasing funds【F:core/loanpool.go†L37-L75】.
Creators retain control over active requests: `CancelProposal` removes a proposal
that has not yet been disbursed, while `ExtendProposal` pushes the voting deadline
forward by a specified number of hours【F:core/loanpool.go†L93-L116】. Clients may
query individual items through `GetProposal` or obtain a sorted list with
`ListProposals`, which supports off‑chain indexing and audit trails【F:core/loanpool.go†L77-L91】.

## Loan Applications
For retail‑style borrowing the LoanPoolApply module allows direct applications.
Applications collect applicant identity, amount, term, and purpose, and are
approved once community votes exceed zero. Disbursement checks treasury balance
and application status before transferring funds【F:core/loanpool_apply.go†L33-L81】.
CLI support is available through `synnergy loanpool_apply`, enabling submission,
voting, processing, disbursement, and inspection of applications【F:cli/loanpool_apply.go†L15-L88】.

## Loan and Grant Programs
The Loanpool supports multiple funding instruments tailored for different
stakeholders:

- **Repayable proposals** – capital‑intensive projects are submitted as
  `LoanProposal` records and must reach community quorum and authority
  ratification before disbursement【F:core/loanpool.go†L24-L75】.
- **Micro‑loans** – streamlined retail requests flow through
  `LoanPoolApply`, where a single affirmative vote is sufficient for approval,
  enabling rapid support for individuals and small businesses【F:core/loanpool_apply.go†L5-L46】【F:core/loanpool_apply.go†L60-L81】.
- **Non‑repayable grants** – the SYN3800 grant registry tracks beneficiaries,
  caps, and release notes. Grants are created and disbursed through dedicated
  APIs and opcodes (`Loanpool_CreateGrant`, `Loanpool_ReleaseGrant`, and
  `Loanpool_GetGrant`) and can be administered from the `syn3800` CLI【F:core/syn3800.go†L8-L56】【F:core/opcode.go†L604-L620】【F:cli/syn3800.go†L14-L81】.

### Named Loan Contract Types
Each `LoanProposal` includes a free‑form `Type` field for labelling programmes such as infrastructure or education【F:core/loanpool_proposal.go†L5-L17】. Beyond these labels, the repository defines specialised loan contracts:

- **CollateralizedLoan** – secures borrowings against pledged assets to ensure principal protection before funds are released【F:smart-contracts/solidity/CollateralizedLoan.sol†L4-L6】.
- **UndercollateralizedLoan** – placeholder for credit‑based lending that accepts reduced collateral when borrowers meet off‑chain risk criteria【F:smart-contracts/solidity/UndercollateralizedLoan.sol†L4-L6】.
- **StableRateLoan** – anchors repayments to a fixed interest schedule, offering predictable liabilities regardless of market volatility【F:smart-contracts/solidity/StableRateLoan.sol†L4-L6】.
- **VariableRateLoan** – adjusts interest charges dynamically, allowing terms to track governance‑defined or oracle‑supplied benchmarks【F:smart-contracts/solidity/VariableRateLoan.sol†L4-L6】.
- **FlashLoan** – enables uncollateralised borrowing within a single transaction, requiring repayment before transaction completion to avoid reversal【F:smart-contracts/solidity/FlashLoan.sol†L4-L6】.
- **LoanAuction** – framework for market‑driven rate discovery where lenders bid to fund proposals, improving treasury utilisation【F:smart-contracts/solidity/LoanAuction.sol†L4-L6】.
- **LoanFactory** and **LoanRegistry** – scaffolding contracts for deploying standardised loan instances and cataloguing outstanding positions across the network【F:smart-contracts/solidity/LoanFactory.sol†L4-L6】【F:smart-contracts/solidity/LoanRegistry.sol†L4-L6】.

### Grant Contract Variants
Grant records leverage named categories for programme accounting. The SYN3800 registry stores beneficiary, category `Name`, total allocation, disbursed amount, and annotated release notes【F:core/syn3800.go†L8-L56】. The `Name` field labels grants—tests demonstrate a "research" award as a canonical example【F:core/syn3800_test.go†L5-L8】. Additional smart‑contract stubs outline specialised grant flows:

- **GrantTracker** – foundation for comprehensive on‑chain tracking of grant milestones and spending【F:smart-contracts/solidity/GrantTracker.sol†L4-L6】.
- **GrantMatching** – intended to pair community donations with approved proposals, enabling matched‑fund campaigns【F:smart-contracts/solidity/GrantMatching.sol†L4-L6】.
- **ScholarshipFund** – dedicated vehicle for education‑oriented grants and scholarship disbursements【F:smart-contracts/solidity/ScholarshipFund.sol†L4-L6】.

## State Serialisation and Views
Both proposals and applications expose JSON‑tagged view structs, allowing wallets
and dashboards to consume normalised metadata without direct access to internal
pointers. `LoanProposalView` mirrors the proposal fields and aggregates vote counts,
while `LoanApplicationView` reports applicant data and approval state【F:core/loanpool_views.go†L5-L33】【F:core/loanpool_views.go†L54-L78】.

## Administrative Controls
Authorities or designated operators can manage the Loanpool via the manager
interface. The `Pause` and `Resume` methods flip the pool’s global `Paused` flag to
halt or accept submissions, protecting the treasury during audits or emergency
conditions. `Stats` iterates across stored proposals to report counts of total,
approved, and disbursed items alongside the current treasury balance【F:core/loanpool_management.go†L21-L43】.
These capabilities are exposed to operators through `synnergy loanmgr` CLI
commands for pausing, resuming, and retrieving statistics【F:cli/loanpool_management.go†L13-L41】.

## Authority Oversight and Voting
Authority nodes form a governance layer that validates Loanpool activity.
Candidates register with roles and accumulate one‑vote‑per‑address support from
network participants; the registry exposes methods to register, vote, sample
electorates, and verify membership【F:core/authority_nodes.go†L45-L98】. The VM
provides helper opcodes such as `Loanpool_RandomElectorate` and
`Loanpool_IsAuthority` so contracts can select oversight committees or enforce
authority‑only flows【F:core/opcode.go†L603-L610】.

Authority review of funding requests is coordinated through
`Loanpool_RequestApproval`, `Loanpool_ApproveRequest`, and
`Loanpool_RejectRequest`, allowing compliant disbursements to be ratified on‑chain
and exposing every decision for audit【F:core/opcode.go†L615-L620】. Operationally,
the `synnergy authority` CLI lets operators register nodes, cast votes, and sample
electorates for proposal due diligence【F:cli/authority_nodes.go†L14-L52】.

## Smart Contract Opcodes
The Synnergy virtual machine exposes dedicated opcodes for Loanpool operations,
covering proposal submission, voting, treasury management, and application
processing. Core opcodes include `Loanpool_Submit`, `Loanpool_Vote`, `Loanpool_Tick`,
`Loanpool_CancelProposal`, and application‑level calls such as `LoanApply_Disburse`
and `LoanApply_List`. Additional opcodes reserve space for advanced features like
`Loanpool_Redistribute` to cycle idle capital, `Loanpool_CreateGrant` and
`Loanpool_ReleaseGrant` for grant management, `Loanpool_RequestApproval` and
`Loanpool_ApproveRequest` for authority oversight, and electorate utilities such
as `Loanpool_RandomElectorate` and `Loanpool_IsAuthority`【F:core/opcode.go†L598-L633】.

## Genesis Treasury
The genesis configuration assigns a dedicated wallet to the loan pool under the
`treasury_wallets.loan_pool` entry, ensuring initial funding and traceability for
all disbursements from network launch【F:configs/genesis.json†L4-L22】.

## Command‑Line Integration
The CLI suite offers comprehensive access to Loanpool features:

- `synnergy loanpool` manages proposal submission, voting, disbursement, and
  lifecycle actions such as ticking, listing, cancellation, and deadline
  extensions【F:cli/loanpool.go†L14-L120】.
- `synnergy loanpool_apply` handles application‑level operations from submission
  through disbursement, including voting, batch processing, and view retrieval【F:cli/loanpool_apply.go†L15-L88】.
- `synnergy loanmgr` provides administrative commands for pausing, resuming, and
  viewing statistics【F:cli/loanpool_management.go†L13-L41】.
- `synnergy loanproposal` offers tools for standalone proposal experimentation
  outside the main pool【F:cli/loanpool_proposal.go†L18-L95】.
- `synnergy authority` registers authority nodes, casts governance votes, and
  samples electorates for proposal oversight【F:cli/authority_nodes.go†L14-L52】.
- `syn3800` manages grant records, allowing creation, disbursement, inspection,
  and listing of SYN3800 grants【F:cli/syn3800.go†L14-L81】.
- A helper script simplifies submitting proposals from the command line
  (`cmd/scripts/loanpool_apply.sh`).【F:cmd/scripts/loanpool_apply.sh†L1-L10】

## Testing and Validation
Unit tests exercise the full lifecycle of proposals and applications, covering
submission, voting, approval, disbursement, and view serialisation. The
`TestLoanPoolProposalLifecycle` test ensures that both proposals and applications
transition through approval and disbursement correctly while updating treasury
totals and exposing accurate views【F:core/loanpool_test.go†L5-L42】.

## Risk Management and Bank Integration
Loan contracts integrate credit scoring and collateral valuation, while insurance
tokens can be attached for programmatic claim processing. Bank nodes interfacing
with the Loanpool must maintain KYC records, report performance metrics, and
handle collateral liquidations, all under the oversight of authority nodes that
can freeze operations if risk thresholds are breached【F:Synnergy_Network_Future_Of_Blockchan.md†L312-L327】.

## Conclusion
Through a combination of deterministic smart‑contract logic, multi‑layered
governance, and rich tooling, the Loanpool stands as a flagship component of the
Synnergy Network’s branded financial infrastructure. It channels idle capital to
productive uses while upholding transparency, compliance, and community
participation.

