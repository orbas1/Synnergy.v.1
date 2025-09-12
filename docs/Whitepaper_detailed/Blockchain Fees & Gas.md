# Blockchain Fees & Gas

## Overview
At **Neto Solaris**, the Synnergy Network employs a comprehensive fee and gas model to allocate resources, discourage spam, and reward participants. Every operation performed on the network consumes gas, and the corresponding fees ensure sustainable network economics and alignment with community stakeholders.

## Gas Accounting
Synnergy prices computation through a dedicated **GasTable** that maps each opcode to a deterministic cost. The table is parsed at runtime from `docs/reference/gas_table_list.md` and cached behind a `sync.Once` gate so all nodes share a consistent schedule. Unlisted operations fall back to a **DefaultGasCost** of `1`, allowing experimental opcodes to execute without prohibitive expense.

### Runtime Controls
- **Query and Introspection** – `GasCost` and `HasOpcode` expose lookups for wallets and explorers, while `GasCostByName` resolves exported function names to their current prices.
- **Dynamic Updates** – Operators can inject or override prices using `RegisterGasCost`; `ResetGasTable` clears the cache so tests or governance actions can reload revised schedules without restarting processes.
- **Audit Support** – `GasTableSnapshot` and its JSON variant provide immutable views of the live schedule for monitoring systems or compliance archives.

### CLI Integration
The `gas list` command prints the active cost of every registered opcode, ensuring developers and auditors can verify pricing directly from the command line.

## Fee Components
Every transaction fee is decomposed into three parts managed by the `FeeBreakdown` structure:
1. **Base Fee** – `CalculateBaseFee` derives a baseline from the median of the latest 1,000 block fees. `AdjustFeeRates` then scales this value to reflect current network load, while the optional AI routine `OptimiseBaseFee` can forecast an optimal starting point from recent statistics.
2. **Variable Fee** – `CalculateVariableFee` multiplies gas units by the current gas price per unit, translating opcode consumption into currency.
3. **Priority Fee** – `CalculatePriorityFee` captures the user’s tip, allowing transactions to signal urgency.

These components combine into a transparent total that wallets and explorers can display or audit in real time.

## Transaction-Specific Calculations
Synnergy supports tailored fee calculators for distinct transaction categories:
- **Transfers** – Fees scale with payload size.
- **Purchases** – Based on the number of contract calls executed.
- **Token Interactions** – Measured by computation units used when interacting with deployed tokens.
- **Contract Operations** – Reflect the complexity factor of new or modified contracts.
- **Wallet Verification** – Adjusted by the security level required.
- **Validated Transfers** – Eligible transfers verified by the network can execute fee-free.
A generic `EstimateFee` helper routes requests to the appropriate calculator, simplifying client implementations.

## Distribution of Fees
`DistributeFees` allocates revenue across nine pools, ensuring every collected unit advances a pillar of the ecosystem:

| Allocation Target            | Percentage |
| ---------------------------- | ---------- |
| Validators and miners        | 64%        |
| Internal development         | 5%         |
| Internal charity             | 5%         |
| External charity             | 5%         |
| Loan pool for ecosystem growth | 5%        |
| Passive-income program       | 5%         |
| Authority nodes              | 5%         |
| Node hosts                   | 5%         |
| Creator wallet               | 1%         |

Distribution contracts can credit these shares directly to ledger accounts via `FeeDistributionContract`. For bespoke arrangements—such as side agreements between validators—`ShareProportional` accepts weightings and reconciles rounding remainders so every unit of value is accounted for.

Deterministic addresses defined by `DefaultGenesisWallets` and seeded through `AllocateToGenesisWallets` receive these allocations at network launch, providing transparent treasuries for development, charity, and infrastructure upkeep.

## Fee Policies and Adjustments
- **Cap and Floor Enforcement** – The `FeePolicy` wrapper invokes `ApplyFeeCapFloor` and returns descriptive notes whenever limits are hit, giving clients a deterministic record of adjustments.
- **Dynamic Rate Adjustment** – `AdjustFeeRates` scales base and variable components with network load, maintaining equilibrium as congestion rises or falls.
- **Proportional Sharing** – `ShareProportional` divides arbitrary fee pools using integer weights, assigning any remainder to the highest-weighted participant.
- **Block Utilisation Rewards** – `AdjustForBlockUtilization` increases or decreases validator payouts by 10% based on actual block fill, incentivising efficient packing.

## Operational Considerations
- **Transaction Structure**: Each transaction includes an explicit fee field, ensuring deterministic ordering and verifiable cost accounting.
- **Reversal Requests**: Users seeking an authority-mediated reversal must reserve funds covering both the transfer amount and a return gas fee, protecting the network from abuse.
- **Evolving Gas Catalogue**: The gas table evolves alongside new features—from cross-chain bridges to biometric security commands—so cost visibility remains current as Synnergy expands.

## Enterprise Tooling and Optimisation
- **CLI Suite** – The `fees estimate` command models charges for transfers, purchases, token interactions, contract deployments and wallet verification while `fees share` computes proportional splits for custom arrangements. These tools mirror on-chain logic so operators can forecast costs before submission.
- **Optimization Nodes** – Dedicated nodes can reorder pending transactions using the `FeeOptimizer`, which sorts by fee density to maximise revenue per byte and maintain throughput under heavy load.
- **AI Insights** – The `OptimiseBaseFee` routine in the AI service digests recent network statistics and suggests base fee targets, enabling autonomous fee markets and predictive scaling.

## Conclusion
Neto Solaris’s fee and gas framework provides a balanced economic engine for the Synnergy Network. By coupling transparent gas accounting with flexible fee policies and equitable distribution, the platform delivers predictable costs, fosters community development, and upholds the long-term viability of the ecosystem.

*© Neto Solaris* 
