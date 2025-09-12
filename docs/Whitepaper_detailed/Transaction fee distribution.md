# Transaction Fee Distribution

Neto Solaris's Synnergy Network employs a transparent and programmable fee model that funds ecosystem growth while incentivizing network actors. This document details how fees are computed, how they are adjusted for enterprise policy, and how every collected unit is allocated across the platform.

## Fee Components

Synnergy transactions expose three fee components, allowing dynamic pricing and user-driven prioritisation:

- **Base Fee** – Derived from the median fee of the last 1,000 blocks and adjusted by current network load. This smooths volatility and anchors fees to recent economic activity.
- **Variable Fee** – Scales with the resources consumed by a transaction. Depending on the transaction type, this may reflect data size, contract calls, computation units or security level.
- **Priority Fee (Tip)** – Optional user-specified amount that elevates a transaction’s priority during congestion.

These components are combined into a `FeeBreakdown` structure, enabling downstream services to reason about each portion independently.

## Fee Calculation Mechanics

1. **Base Fee Sampling**
   - The `CalculateBaseFee` function sorts recent block fees, selects the median, then applies an adjustment factor to reflect network conditions【F:core/fees.go†L28-L46】.
2. **Variable Fee Determination**
   - `CalculateVariableFee` multiplies gas units by gas price per unit, accommodating resource-heavy operations【F:core/fees.go†L48-L50】.
3. **Priority Fee Inclusion**
   - `CalculatePriorityFee` directly returns the user’s tip, allowing wallets to surface the explicit prioritisation amount【F:core/fees.go†L53-L54】.
4. **Transaction-Type Abstraction**
   - `EstimateFee` orchestrates the above components across five transaction categories (transfer, purchase, token interaction, contract, wallet verification) ensuring consistent logic across the ecosystem【F:core/fees.go†L172-L189】.
5. **Network Load Adjustments**
   - `AdjustFeeRates` scales both base and variable fees by a load factor, encouraging throughput during busy periods while remaining predictable under light usage【F:core/fees.go†L160-L169】.
6. **Policy Enforcement**
   - `ApplyFeeCapFloor` constrains fees to governance‑defined limits, and `FeePolicy` surfaces explanatory notes when caps or floors trigger, enabling boards to communicate pricing boundaries to wallets and exchanges【F:core/fees.go†L130-L158】.
7. **Feeless Transfer Validation**
   - `FeeForValidatedTransfer` returns zero fees for pre-approved flows, enabling charitable or promotional transfers to bypass charges while defaulting to standard fees when validation fails【F:core/fees.go†L91-L99】.

## Transaction Prioritisation

Optimisation nodes order mem‑pool transactions by fee density (fee per byte). This deterministic selection ensures blocks maximise revenue while adhering to size constraints【F:internal/nodes/optimization_nodes/index.go†L13-L35】.

## Gas and Opcode Pricing

Variable fees reference a comprehensive gas table that prices each virtual machine opcode. `GasTable` loads costs from `gas_table_list.md`, caching results and falling back to a default for unknown operations【F:gas_table.go†L18-L37】【F:gas_table.go†L49-L88】. Tooling retrieves prices through `GasCost` and can register overrides at runtime for governance-driven adjustments or testing scenarios【F:gas_table.go†L98-L106】【F:gas_table.go†L115-L124】.

## Fee Distribution Framework

Once a transaction is included, the ledger deducts `amount + fee` from the sender and credits the recipient with the amount, isolating the fee for distribution【F:core/ledger.go†L148-L169】.

`DistributeFees` allocates the collected fee across nine strategic pools, each supporting a pillar of the Synnergy ecosystem【F:core/fees.go†L101-L127】:

| Allocation Target        | Percentage |
| ------------------------ | ---------- |
| Internal Development     | 5%         |
| Internal Charity         | 5%         |
| External Charity         | 5%         |
| Loan Pool                | 5%         |
| Passive Income Program   | 5%         |
| Validators & Miners      | 64%        |
| Authority Node Incentives| 5%         |
| Node Host Rewards        | 5%         |
| Creator Wallet           | 1%         |

The validator/miner pool can be further apportioned using `ShareProportional`, which ensures deterministic rounding by assigning any remainder to the highest-weighted participant【F:core/fees.go†L194-L221】.

### Block Utilisation Modifier

To promote efficient block usage, the validator/miner pool is adjusted based on realised block capacity. Utilisation above 90% yields a 10% bonus, while utilisation below 50% incurs a 10% reduction【F:core/fees.go†L240-L255】.

## Genesis and Treasury Addresses

`DefaultGenesisWallets` defines deterministic addresses for each allocation category using hashed labels, ensuring transparency from network launch【F:core/genesis_wallets.go†L8-L41】. The `AllocateToGenesisWallets` helper applies the standard distribution to these wallets, forming the basis of Synnergy’s treasury and philanthropic reserves【F:core/genesis_wallets.go†L44-L58】.

## Distribution Execution

Smart‑contract style distribution is simulated via `FeeDistributionContract`, which credits ledger accounts with their calculated share【F:core/fees.go†L223-L238】. The same logic powers CLI utilities for genesis initialisation, ongoing fee splits and auditability.

## Dispute Resolution and Refunds

In rare cases where transactions must be reversed, authority nodes coordinate a structured refund process. `RequestReversal` freezes the recipient's balance for the transfer amount plus the return fee, preventing double‑spend attempts while the dispute is evaluated【F:core/transaction_control.go†L60-L71】. Once sufficient approvals are gathered within the 30‑day reversal window, `FinalizeReversal` releases the frozen funds and executes a compensating transaction that carries the original fee back to the sender【F:core/transaction_control.go†L79-L105】. If approvals fail or the window expires, `RejectReversal` unfreezes the funds and the request is logged for audit trails【F:core/transaction_control.go†L107-L114】.

## Validator Penalties and Risk Management

To preserve network integrity, the `StakePenaltyManager` tracks validator stake balances and accumulates penalty points for misbehaviour or downtime. Stake can be adjusted up or down and infractions are recorded with reasons and timestamps, enabling governance bodies to modulate future fee shares or slash rewards when necessary【F:stake_penalty.go†L8-L58】.

## Service-Level and Market Fees

Beyond protocol fees, application modules may charge their own service fees. The on‑chain `LiquidityPool` for token swaps expresses fees in basis points (`FeeBps`) and deducts them from swap inputs before updating reserves, allowing DEX operators to feed revenue into the standard distribution channels or custom incentives【F:core/liquidity_pools.go†L23-L28】【F:core/liquidity_pools.go†L95-L106】.

## Audit and Compliance Logging

Enterprise deployments require transparent reporting. The `AuditManager` records per-address events with metadata and timestamps, allowing downstream systems to reconstruct fee movements【F:core/audit_management.go†L9-L49】. Building on this ledger, the `ComplianceService` validates KYC commitments, records fraud signals, maintains risk scores and exposes immutable audit trails for each address. It can also monitor transaction amounts for anomalies and verify zero‑knowledge proofs for privacy-preserving compliance checks【F:core/compliance.go†L55-L99】【F:core/compliance.go†L111-L133】.

## CLI Tooling

Neto Solaris provides command‑line tools to interact with the fee system:

- `tx fee` – Estimates fees and shows the resulting distribution, enforcing optional caps and floors via `--cap` and `--floor` flags【F:cli/transaction.go†L79-L117】.
- `fees estimate` – Quick fee forecasts incorporating live network load【F:cli/fees.go†L17-L55】.
- `fees feedback` – Sends user feedback on fee estimates for continuous tuning【F:cli/fees.go†L57-L66】.
- `fees share` – Computes proportional splits between validators and miners for customised weighting【F:cli/fees.go†L74-L83】.
- `genesis allocate` – Applies distribution logic to genesis wallets, displaying per‑address allocations【F:cli/genesis.go†L36-L50】.

## Self-Sustaining Economic Model

Synnergy’s fee allocation ensures the chain finances its own evolution without perpetual external funding. Each block channels a defined percentage of revenue into engineering, infrastructure and community treasuries through `DistributeFees`, creating a closed loop where usage directly bankrolls future upgrades and outreach【F:core/fees.go†L101-L127】. Deterministic wallet assignments codified in `DefaultGenesisWallets` and `AllocateToGenesisWallets` guarantee these reserves are transparently routed to long-term development and operations from genesis onward【F:core/genesis_wallets.go†L28-L58】.

Validator rewards dynamically scale with network demand. `AdjustForBlockUtilization` increases the validator/miner pool when blocks are near capacity and tapers payouts during periods of underuse, aligning security expenditure with real throughput and keeping resource costs predictable for enterprises【F:core/fees.go†L240-L255】. Governance bodies can further stabilise revenue by enforcing fee caps or floors, shielding mission-critical applications from volatility while maintaining sufficient income for network stewardship【F:core/fees.go†L130-L158】.

## Universal Basic Income Forecast

Five percent of every fee funds Synnergy’s Passive Income program, a dedicated pool designed to evolve into a universal basic income for all verified participants【F:core/fees.go†L101-L127】. The `IDRegistry` records identity-backed wallets eligible for distributions, preventing Sybil attacks and enabling targeted payouts【F:idwallet_registration.go†L8-L27】. When the treasury releases funds, `ShareProportional` divides the pool across registered addresses with deterministic rounding, allowing administrators to weight allocations or deliver equal shares at scale【F:core/fees.go†L194-L221】.

Assuming an average network fee of 0.02 SYN and 1,000,000 daily transactions, the Passive Income pool would accrue 1,000 SYN per day. With 100,000 registered citizens, each could receive roughly 0.01 SYN daily. As adoption and transaction volume rise, the pool compounds, moving closer to a sustainable baseline income for all participants. Enterprises can forecast longer-term payouts by modelling projected transaction growth against the fixed five‑percent contribution and adjusting eligibility thresholds to balance inclusivity with economic prudence.

## Summary

Through a modular fee model, comprehensive gas accounting, and a clearly defined distribution schema, the Synnergy Network channels transaction revenue into development, community programmes, validator incentives and long‑term sustainability. Neto Solaris ensures every fee reinforces the ecosystem, balancing profitability with social impact, regulatory compliance and operational resilience.

