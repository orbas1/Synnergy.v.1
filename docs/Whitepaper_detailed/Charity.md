# Charity

## Overview
Blackridge Group Ltd. embeds philanthropy into the core of the Synnergy Network. The charity subsystem channels on-chain revenue and community contributions toward transparent, programmable aid. Every donation, vote, and disbursement is auditable on the ledger, ensuring that social impact accompanies technological innovation.

## Funding Allocation
A fixed percentage of every transaction fee is earmarked for charitable causes. The `DistributeFees` policy routes five percent of fees to Blackridge's internal initiatives and another five percent to an external community pool【F:core/fees.go†L101-L127】. This predictable allocation guarantees continual funding regardless of market activity.

## Genesis Configuration and Treasury Accounts
The initial ledger defines dedicated addresses for the internal and external pools, ensuring separation of duties and traceability from day one【F:configs/genesis.json†L4-L21】. Both accounts start with a zero balance, relying entirely on the automated fee stream and community donations for funding. Addresses are deterministically derived from hashed labels and the fee distribution routine maps allocations to these wallets, enabling independent auditing of genesis funds【F:core/genesis_wallets.go†L8-L59】.

## Charity Pool Architecture
Two dedicated ledger accounts segregate internal and external funds: `internal_charity` and `charity_pool`【F:core/charity.go†L66-L69】. The `CharityPool` smart component coordinates registrations, deposits, voting, and future payouts while guarding state with a mutex for concurrent safety【F:core/charity.go†L57-L60】:

- **Deposits:** Transfers into `charity_pool` are validated to ensure sufficient balance and a non-zero amount before crediting the pool【F:core/charity.go†L76-L87】.
- **Registrations:** Organizations submit an address, name, and mission category, which are persisted under `charity:reg:*` keys for auditability【F:core/charity.go†L90-L100】.
- **Voting:** Eligible SYN900 identity holders can cast a vote for a registered charity each funding cycle; votes are stored using the `charity:vote:*` prefix【F:core/charity.go†L103-L113】.
- **Record Retrieval:** Registration data can be queried per cycle for transparency and compliance tracking【F:core/charity.go†L125-L136】.
- **Maintenance:** Manual `tick` operations reserve space for future automated payouts and cycle management【F:cli/charity.go†L97-L113】.
- **Winners:** The CLI exposes a `winners` command, preparing the infrastructure for distributing pooled funds to top‑voted charities【F:cli/charity.go†L145-L170】.

The pool currently tracks votes and registrations while leaving automatic disbursement logic to future iterations. Categories span hunger relief, child welfare, wildlife and marine support, disaster response, and war relief, enabling targeted reporting and governance【F:core/charity.go†L15-L22】. All ledger interactions flow through a narrow `StateRW` interface, allowing alternative backends and prefix‑based queries for auditors without exposing unrelated state【F:core/state_rw.go†L11-L19】.

## Community Workflow
1. **Register** – A charity submits its wallet, name, and category to join the pool.
2. **Vote** – SYN900 token holders endorse preferred charities during each cycle; non‑holders are rejected at the consensus layer【F:core/charity.go†L103-L110】【F:core/charity_test.go†L22-L29】.
3. **Donate** – Anyone may contribute additional tokens to strengthen the external pool or Blackridge’s internal fund.
4. **Disburse** – At cycle end, top‑voted charities become eligible for payouts from accumulated fees and donations.

## Command‑Line Interface
The network CLI exposes comprehensive controls for charity operations:

- `synnergy charity_pool register [addr] [category] [name]` — register a charity.
- `synnergy charity_pool vote [voterAddr] [charityAddr]` — cast a vote.
- `synnergy charity_pool tick [timestamp]` — advance the pool’s cron tasks.
- `synnergy charity_pool winners [cycle]` — list winning charities for a cycle.
- `synnergy charity_mgmt donate [from] [amt]` — deposit funds into the pool.
- `synnergy charity_mgmt withdraw [to] [amt]` — withdraw internal charity funds.
- `synnergy charity_mgmt balances` — display pooled and internal balances.

Each command optionally emits JSON for integration with dashboards and automation tools【F:cli/charity.go†L68-L222】. Unit tests confirm the JSON flag works for both registration lookups and balance reports, enabling reliable machine consumption【F:cli/charity_test.go†L8-L32】.

## Ledger Integration and Auditing
The `AllocateToGenesisWallets` routine applies the fee policy to the deterministically derived addresses, producing a verifiable map of allocations for compliance review【F:core/genesis_wallets.go†L44-L59】. The pool’s reliance on the `StateRW` abstraction allows operators to plug in enterprise ledger implementations while retaining prefix-scoped inspection and balance checks.

## Security and Quality Assurance
Mutex protections around pool state and read‑write locks within the SYN4200 token registry safeguard concurrent operations under heavy load【F:core/charity.go†L57-L60】【F:core/syn4200_token.go†L16-L30】. Dedicated tests validate deposits, registrations and voter eligibility, ensuring that invalid participants cannot influence outcomes【F:core/charity_test.go†L34-L52】.

## SYN4200 Charity Tracking Token
The SYN4200 token family records donations to specific campaigns and exposes progress metrics. Campaigns capture a symbol, purpose, goal, and individual contributions, allowing sponsors to audit fundraising in real time【F:core/syn4200_token.go†L7-L42】. The registry guards access with a read‑write mutex so high‑volume donation flows remain thread‑safe【F:core/syn4200_token.go†L16-L30】. CLI helpers let donors contribute and query progress directly from scripts or enterprise backends【F:cli/syn4200_token.go†L13-L49】.

## Opcode and Virtual Machine Hooks
At the opcode layer the network reserves dedicated slots for charity functions, from pool management to token progress queries【F:core/opcode.go†L292-L301】【F:core/opcode.go†L988-L993】. These opcodes stabilise gas metering and give the virtual machine consistent entry points for monitoring or pausing charitable activity during audits.

## Enterprise Integration and Reporting
Operational commands support machine‑readable output, scripted donations and automated reconciliation. The management CLI exposes deposit, withdrawal and balance queries for treasury teams【F:cli/charity.go†L178-L220】, while the SYN4200 token CLI records campaign donations and reports progress for dashboards【F:cli/syn4200_token.go†L13-L49】.

## Transparency and Future Directions
All charity interactions are immutably stored, aligning Blackridge Group Ltd.'s charitable commitments with verifiable blockchain records. Planned enhancements include automated payout workflows, expanded campaign analytics, and on-chain governance to fine‑tune funding strategies.

