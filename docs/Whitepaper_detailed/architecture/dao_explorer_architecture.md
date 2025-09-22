# DAO Explorer Architecture

## Overview
The DAO Explorer gives stakeholders a streamlined view into decentralized governance activity. It wraps the `dao` CLI commands with a lightweight web interface so users can create organisations, stake tokens and vote on proposals without managing raw JSON outputs.

## Key Modules
- `cli/dao.go` – top-level command wiring together proposal, staking and token subcommands.
- `cli/dao_proposal.go` – create, list and execute governance proposals.
- `cli/dao_staking.go` – manage staking balances that grant voting power.
- `cli/dao_quadratic_voting.go` – apply quadratic weighting to cast votes.
- `core/dao_proposal.go` – in-memory store for proposals and results.

## Workflow
1. **Organisation discovery** – the explorer lists DAOs by shelling out to `synnergy dao list`.
2. **Proposal creation** – users draft proposals which `dao_proposal` submits to the core manager.
3. **Staking** – `dao_staking` records token deposits and exposes balances.
4. **Voting** – votes are cast through `dao_quadratic_voting` to mitigate whales.
5. **Execution** – successful proposals trigger follow‑up actions such as parameter changes or fund dispersal.

## Security Considerations
- The explorer operates in read‑only mode for most queries, keeping private keys in the CLI layer.
- Proposal and vote submissions require signatures from staking accounts.
- All transactions are logged with proposal IDs to maintain an auditable trail.

## CLI Integration
- `synnergy dao` – entry point for DAO operations.
- `synnergy dao-proposal` – manage proposals.
- `synnergy dao-stake` – deposit or withdraw staking amounts.

## Enterprise Diagnostics
- `synnergy integration status` exercises DAO staking automatically by crediting a validator wallet and observing the block that finalises the transaction. This proves the explorer’s backing services remain hydrated without a manual smoke test.
- The Enterprise Integration Health widget in the web UI relays the same diagnostics stream so governance operators can spot consensus or wallet regressions at a glance.
