# Applying to the Synnergy Charity Pool

_Official guidance prepared by **Neto Solaris**_

## Overview
The Synnergy Charity Pool is an on‑chain treasury that channels a portion of network fees and voluntary donations toward registered philanthropic initiatives. Charities register with the pool, receive votes from the community, and, at the close of each funding cycle, top‑voted organisations become eligible for disbursements. All interactions—registrations, deposits, votes, and payouts—are immutably recorded on the ledger for full transparency.

## Funding Sources and Fee Allocation
The pool is capitalised through two mechanisms:

- **Protocol Fees** – Five percent of every transaction fee is routed into the external charity pool while a separate five percent supports Neto Solaris’s internal charitable programmes【F:core/fees.go†L115-L123】. These allocations are enforced by the network’s `DistributeFees` policy, ensuring predictable funding regardless of market activity.
- **Voluntary Donations** – Community members or corporate sponsors may augment the pool by depositing additional tokens to the `charity_pool` account, increasing the resources available for distribution.

All funds destined for external causes accumulate under the ledger address `charity_pool`, distinct from `internal_charity` which holds Neto Solaris‑directed initiatives【F:core/charity.go†L66-L69】.

### Transaction Fee Distribution Model
To provide full transparency to network participants, the Synnergy protocol publishes the complete breakdown of how each transaction fee is allocated:

| Allocation Target | Percentage |
|-------------------|-----------:|
| Internal Development | 5% |
| **Internal Charity (Neto Solaris programmes)** | **5%** |
| **External Charity Pool** | **5%** |
| Loan Pool | 5% |
| Passive Income Dividend | 5% |
| Validators & Miners | 64% |
| Authority Nodes | 5% |
| Node Hosts | 5% |
| Creator Wallet | 1% |

These values are hard‑coded in the `DistributeFees` routine, guaranteeing that charity allocations remain proportionate regardless of fee volume【F:core/fees.go†L115-L127】.

## Eligibility and Preparation
Before submitting an application, prospective charities should ensure the following:

- **Blockchain Address** – Each charity must control a Synnergy‑compatible wallet capable of receiving token transfers.
- **Clear Mission Category** – Applicants select one of the predefined mission classes such as _HungerRelief_, _ChildrenHelp_, _WildlifeHelp_, _SeaSupport_, _DisasterSupport_, or _WarSupport_ when registering【F:core/charity.go†L12-L22】.
- **Organisation Details** – A human‑readable name and concise description of activities help voters understand the cause.
- **Compliance** – Charities should be legally registered within their jurisdictions and prepared to provide any documentation requested by Neto Solaris during due diligence.

## Step‑by‑Step Application Process
1. **Prepare a Wallet**  
   Generate or designate a secure wallet address that will act as the charity’s on‑chain identity. This address must remain under the organisation’s control throughout participation.

2. **Gather Registration Data**  
   Compile the charity’s name, chosen mission category, and wallet address. The category determines how the cause is classified in the pool and is stored alongside the registration metadata for auditing purposes【F:core/charity.go†L43-L50】【F:core/charity.go†L90-L100】.

3. **Submit the Registration**  
   Use the Synnergy CLI or an integrated application to transmit the registration:
   ```bash
   synnergy charity_pool register <charity_address> <category_id> "<organisation_name>"
   ```
   - `category_id` corresponds to one of the numeric identifiers in the mission table above.  
   - Upon success, the ledger records the registration under the key `charity:reg:<address>`.

4. **Confirm Ledger Entry**
   After submission, query the registration to verify it was persisted:
   ```bash
   synnergy charity_pool --json registration <charity_address>
   ```
   The response returns the stored metadata, including the selected category and initial vote count.

5. **Renew for Subsequent Cycles**
   Registrations are tied to funding cycles. To remain eligible for future distributions, charities should resubmit or update their registration before each new cycle begins. The `Cycle` field in the registration record facilitates historical auditing and ensures votes apply to the correct period【F:core/charity.go†L43-L50】.

## Community Voting and Evaluation
- **Eligible Voters** – Only holders of SYN900 identity tokens may vote, ensuring that each community member can support a single charity per cycle【F:core/charity.go†L103-L110】.
- **Casting Votes**  
  ```bash
  synnergy charity_pool vote <voter_address> <charity_address>
  ```
  Votes are written to the ledger for auditability and will be tallied when funding decisions are made.
- **Cycle Timing** – The pool operates in discrete cycles beginning from network genesis. Administrators may trigger a maintenance tick to process end‑of‑cycle logic.

## Capacity, Cycle Timing, and Deadlines
The external charity allocation derives from cumulative protocol fees and voluntary contributions, so the amount available each cycle varies with network activity. Unused funds roll over to subsequent cycles, preserving long‑term capacity.

Each cycle includes a **registration window**, a **voting window**, and a **settlement window**:

- **Registration Window** – Charities must submit or renew their registration before this window closes. Late submissions are automatically deferred to the following cycle.
- **Voting Window** – Verified community members cast their votes. Real‑time totals can be audited from the ledger to help organisations gauge support.
- **Settlement Window** – After voting ends, administrators tally results and prepare disbursements. Funds not awarded remain in the pool for future rounds.

Neto Solaris publishes the official calendar, including cycle start and end dates, on our governance portal. Organisations planning major campaigns should file registrations at least one full cycle in advance to guarantee eligibility.

## Funding, Donations, and Disbursement
- **Donations** – Anyone can contribute to the pool by depositing tokens, which increases the pool balance available for future distribution【F:core/charity.go†L76-L87】.
- **Payouts** – At the conclusion of each cycle, top‑voted charities may receive allocations from the pool. The reference implementation currently stores votes and registrations while leaving automated payout logic for future development【F:core/charity.go†L116-L123】.
- **Withdrawal of Internal Funds** – Neto Solaris maintains an internal charity account for corporate initiatives; withdrawals from this account are managed separately from community‑directed funds【F:core/charity.go†L66-L69】.

### Auditing and Balance Queries
Administrators and auditors can verify pool holdings at any time using the management CLI:

```bash
synnergy charity_mgmt balances --json
```

The command returns the current token balances for the community pool and the internal charity account, enabling transparent reconciliation【F:cli/charity.go†L206-L218】.

## Ongoing Responsibilities
- **Transparency** – Successful applicants are expected to publish periodic impact reports so token holders can verify how funds were used.
- **Key Management** – Safeguard private keys controlling the registered address. Compromised addresses cannot be recovered by the network.
- **Governance Participation** – Charities are encouraged to engage with the Synnergy governance process, providing feedback on pool operations and category definitions.

## Enterprise Integration and Compliance
Larger organisations can integrate directly with the Synnergy Network through our APIs and enterprise SDKs. Ledger interactions such as registration, voting, donation tracking, and balance queries can be automated within existing back‑office systems, enabling near real‑time reconciliation and audit trails. Event webhooks and batch export endpoints are provided for enterprise resource‑planning platforms that require periodic snapshots.

Neto Solaris conducts due diligence on all applicants and may request supplemental documentation to satisfy regulatory obligations or anti‑money‑laundering controls. Dedicated compliance contacts are available to help multinational charities align their internal policies with on‑chain transparency requirements.

## Support
For assistance with registration or technical integration, contact the Neto Solaris support team via our official channels. Development proposals and enhancements to the charity module are welcomed through our public repository.

---
*Neto Solaris remains committed to leveraging the Synnergy Network for transparent, community‑driven philanthropy.*
