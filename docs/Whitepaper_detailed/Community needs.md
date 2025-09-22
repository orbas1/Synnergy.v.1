# Community Needs

## Mission
Neto Solaris designs the Synnergy Network to address tangible community priorities. The platform channels on-chain resources toward social welfare, economic empowerment and ecosystem innovation, ensuring that every network participant can contribute to and benefit from shared prosperity.

## Strategic Funding Pools

### LoanPool
The LoanPool aggregates a portion of transaction fees and allocates them to targeted programmes ranging from poverty alleviation to small business support. Managed proposals undergo submission, open voting and timed evaluation before approved disbursements reduce the treasury balance, enabling transparent credit distribution for community initiatives【F:core/loanpool.go†L9-L75】【F:docs/Whitepaper_detailed/guide/loanpool_guide.md†L3-L79】.

Key funding categories include:
- **Poverty Fund** – non‑repayable grants for immediate relief.
- **Unsecured and Secured Loans** – income‑based lending with optional collateral.
- **Business and Personal Grants** – proposal‑driven awards vetted by node and authority votes.
- **Ecosystem Innovation Fund** – resources for research, infrastructure and application development.
- **Healthcare, Education and Environmental Funds** – dedicated pools supporting critical public services and sustainability.
- **Small Business Support** – financing for SMEs to drive local economic growth.

#### Proposal Lifecycle
Proposals are time‑boxed to a 24‑hour voting window and created only when the pool is unpaused【F:core/loanpool.go†L26-L35】. Votes can be cast until the deadline expires, after which the `Tick` routine automatically approves any proposal with recorded votes【F:core/loanpool.go†L37-L56】. Approved items trigger treasury checks during disbursement to prevent overspending【F:core/loanpool.go†L60-L75】.

#### Application Pipeline
For smaller consumer loans, the `LoanPoolApply` service accepts streamlined applications, registers voter endorsements, finalises approvals with a single recorded vote and performs disbursements against the same treasury ledger【F:core/loanpool_apply.go†L5-L82】. This pipeline enables rapid financing while still enforcing community oversight.

#### Administrative Controls
A dedicated manager can pause or resume proposal intake and provides real‑time statistics on treasury levels, proposal counts and disbursement activity for operational dashboards【F:core/loanpool_management.go†L3-L44】. Creators retain unilateral control to cancel or extend their own proposals before funds are released, safeguarding against rushed or malicious submissions【F:core/loanpool.go†L93-L116】.

#### Audit Views and CLI Integration
Serialisable proposal and application views expose creator details, vote totals and disbursement status for external dashboards and audits【F:core/loanpool_views.go†L5-L97】. Comprehensive CLI tooling allows administrators to submit, vote, process, cancel and extend proposals or list the full ledger of requests directly from the command line, supporting automated governance workflows【F:cli/loanpool.go†L15-L120】.

### Charity Pool
The Charity Pool coordinates donations and voting for registered organisations across categories such as hunger relief, child welfare, wildlife preservation, marine support, disaster response and war assistance. Voter eligibility hinges on SYN900 identity tokens, preserving one‑person‑one‑vote integrity while ledger entries record deposits, registrations and ballots for full auditability【F:core/charity.go†L12-L113】.

#### Categories and Registration
Six built‑in categories—ranging from hunger relief to war support—anchor funding priorities at the protocol level【F:core/charity.go†L12-L22】. Registrations capture each charity’s address, name and category; all writes are mutex‑protected to prevent race conditions during concurrent submissions【F:core/charity.go†L43-L100】.

#### Deposits and Fund Handling
Deposits validate ledger configuration, ensure balances cover the contribution and transfer funds into the pooled account, guaranteeing that donations are backed by actual holdings【F:core/charity.go†L76-L88】.

#### Voting and Cycle Operations
Voting invokes an electorate interface that checks SYN900 identity tokens, enforcing one‑person‑one‑vote policies and persisting ballots for later tabulation【F:core/charity.go†L103-L113】. Periodic `Tick` operations handle future payout logic, while both winning selections and historical registrations are queryable for compliance reviews【F:core/charity.go†L116-L136】.

#### CLI Management
A full command‑line suite supports pool registration, voting, manual ticking, donation management and balance audits, enabling enterprise automation of charity workflows【F:cli/charity.go†L72-L224】.

### Fee Distribution and Treasury Routing
Network fees are automatically split between development, charitable reserves, the LoanPool, validator rewards and other stakeholders. Five percent of every fee supports internal development, with a further five percent each directed to internal and external charities; ten percent is reserved for the LoanPool, five percent for passive income programs, fifty‑nine percent for validators and miners, and the remainder for authority nodes, node hosts and the creator wallet【F:core/fees.go†L101-L128】. Deterministic genesis wallet addresses are generated by hashing human‑readable labels, allowing any observer to recompute the addresses that receive these distributions from the first block onward【F:core/genesis_wallets.go†L22-L58】.

Fee policies can impose network‑wide caps and floors, adjusting charges to stay within governance‑approved limits and returning explanatory notes when thresholds are triggered【F:core/fees.go†L130-L158】.

## Governance and Participation
Community members exercise influence through proposal voting, charity selections and ongoing polls. Only verified identity‑token holders can cast ballots in the Charity Pool, ensuring accountable participation and preventing duplicate voting【F:core/charity.go†L103-L113】. Similar mechanisms underpin LoanPool governance, where proposal approval depends on recorded votes and deadlines【F:core/loanpool.go†L37-L56】. Creators retain control over their submissions with pause, cancel and extension capabilities, providing safeguards against misuse or rushed decision‑making【F:core/loanpool_management.go†L21-L29】【F:core/loanpool.go†L93-L116】.

## Transparency and Accountability
All fund movements and governance actions are immutably logged on the ledger. Proposal and application views provide live snapshots of funding decisions, while `GetRegistration` exposes charity metadata for any cycle and `ProposalInfo` returns serialised loan proposals for external verification【F:core/charity.go†L125-L136】【F:core/loanpool_views.go†L35-L52】. Genesis wallet labelling enables independent auditing from deterministic first‑block addresses【F:core/genesis_wallets.go†L22-L58】, and the extensive CLI surface ensures that every governance action can be reproduced or inspected from scripted environments【F:cli/loanpool.go†L15-L120】【F:cli/charity.go†L72-L224】.

## Social Impact and Sustainability
By dedicating treasury streams to welfare, education, healthcare and environmental stewardship, the Synnergy Network embeds social responsibility at protocol level. Targeted funding categories empower communities to respond to immediate needs while also investing in long‑term resilience and innovation.

## Stage 78 Enterprise Enhancements
- **Community diagnostics:** The enterprise orchestrator correlates LoanPool and Charity Pool health with consensus, wallet and node readiness so municipal operators can monitor funding pipelines via `synnergy orchestrator status` or the web dashboards without bespoke scripts.【F:core/enterprise_orchestrator.go†L21-L166】【F:cli/orchestrator.go†L6-L75】
- **Predictable treasury costs:** Stage 78 gas documentation covers orchestrator operations used to elect authority committees, audit community disbursements and seal wallets, keeping civic workflows affordable and auditable across CLI, VM and browser tooling.【F:docs/reference/gas_table_list.md†L420-L424】【F:snvm._opcodes.go†L325-L329】
- **Tested resilience:** Expanded orchestrator test suites subject funding modules to unit, situational, stress, functional and real-world scenarios, ensuring credit issuance and charitable payouts remain reliable during demand spikes or regulatory reviews.【F:core/enterprise_orchestrator_test.go†L5-L75】【F:cli/orchestrator_test.go†L5-L26】

## Conclusion
The Community Needs framework demonstrates how Neto Solaris integrates social impact with decentralized finance. Through structured funding pools, transparent governance and deterministic accounting, the Synnergy Network ensures that community welfare remains a core outcome of network activity.

