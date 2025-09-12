# Block Rewards Dispersions and Halving

*Prepared by Neto Solaris*

## 1. Emission Overview
Synnergy Network's native currency, the Synthron coin (SYN), follows a capped issuance model. Key monetary constants are baked into the core protocol: a maximum supply of 500 million SYN, a 5 million SYN genesis allocation, an initial block reward of 1,252 SYN, and a halving interval of 200,000 blocks【F:core/coin.go†L6-L16】. With a three‑second target block time, each halving epoch spans roughly 6.9 days【F:configs/network.yaml†L12-L19】.

## 2. Halving Mechanism
The halving schedule governs issuance using the relation `Reward(height) = 1252 / 2^n`, where `n` equals the number of completed 200,000‑block intervals. Once the shift has eliminated all significant bits, rewards fall to zero, enforcing the fixed supply. The `CirculatingSupply` and `RemainingSupply` helpers iterate through historical rewards to model supply at any block height and clamp output once the cap is reached【F:core/coin.go†L19-L49】.

### 2.1 Emission Forecast
The deterministic issuance curve rapidly approaches the cap, as summarised below.

| Epoch | Approx. day | Reward per block (SYN) | Coins minted in epoch (M SYN) | Cumulative supply (M SYN) |
|------:|------------:|-----------------------:|------------------------------:|--------------------------:|
| 1 | 6.9 | 1,252.00 | 250.40 | 255.40 |
| 2 | 13.9 | 626.00 | 125.20 | 380.60 |
| 3 | 20.8 | 313.00 | 62.60 | 443.20 |
| 4 | 27.8 | 156.50 | 31.30 | 474.50 |
| 5 | 34.7 | 78.25 | 15.65 | 490.15 |
| 6 | 41.7 | 39.13 | 7.83 | 497.98 |
| 7 | 48.6 | 19.56 | 2.03 | 500.00 |

Beyond the seventh epoch the block reward underflows to zero and the circulating supply remains at the 500 million SYN cap.

## 3. Reward Dispersion Model
Minted rewards per block are distributed exclusively to consensus participants. Synnergy’s consensus engine seeds weights of 30 % to Proof‑of‑Stake validators, 30 % to Proof‑of‑History schedulers and 40 % to Proof‑of‑Work miners【F:core/consensus.go†L41-L43】【F:configs/network.yaml†L12-L19】. The dispersion is expressed as:

```
A_pos      = 0.30 * A_block  # PoS validators
A_poh      = 0.30 * A_block  # PoH schedulers
A_pow      = 0.40 * A_block  # PoW miners
```

`A_pos` and `A_poh` are further subdivided in proportion to validator stake and time‑slot attendance, while rewards are subject to a cooldown period to discourage short‑term hopping【F:Synnergy_Network_Future_Of_Blockchan.md†L405-L417】. Treasuries and ecosystem funds are financed solely from transaction fees, ensuring block rewards remain an incentive for security.

### 3.1 Dynamic Weight Rebalancing
Weights are not static. `AdjustWeights` recalibrates the distribution based on observed demand and stake concentration, clamping each mechanism between 7.5 % and 100 % before renormalising to unity【F:core/consensus.go†L61-L74】. This allows Neto Solaris operators to emphasise security primitives during periods of elevated load or threat.

### 3.2 Runtime Reward Controls
Operational teams may disable classes of validators or suspend PoW compensation without redeploying the network. `SetAvailability` toggles each consensus path, and `SetPoWRewards` halts payouts while leaving PoW participation for utility purposes【F:core/consensus.go†L130-L142】. The CLI exposes these levers via `synnergy consensus availability [pow] [pos] [poh]` and `synnergy consensus powrewards [enabled]`, the latter emitting audit logs when rewards are toggled【F:cli/consensus.go†L164-L176】.

## 4. Supply Tracking and Operational Tooling
Real‑time monitoring of emission progress is available through the CLI. `synnergy coin reward [height]` reveals the reward at a particular height, and `synnergy coin supply [height]` reports the circulating and remaining supply using the core helpers mentioned above【F:cli/coin.go†L31-L57】.

## 5. Economic and Security Implications
The aggressive halving cadence front‑loads distribution, quickly rewarding early participants while converging on the fixed supply to instill long‑term scarcity. Splitting rewards across Proof‑of‑Stake, Proof‑of‑History, and Proof‑of‑Work actors diversifies incentives and aligns security with resource contribution. As the reward diminishes, transaction fees and staking yields become the dominant revenue streams, promoting sustainable network operation.
