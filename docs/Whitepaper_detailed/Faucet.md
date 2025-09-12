# Token Faucet

Neto Solaris provides a controlled token faucet within the Synnergy Network to dispense small amounts of test assets for development, onboarding, and QA activities. This document details the faucet's architecture, operational model, and usage patterns across the stack.

## Design Objectives

- **Developer enablement** – offer a reliable source of tokens for sandbox environments without jeopardising monetary policy.
- **Abuse resistance** – enforce deterministic rate limiting and balance checks to prevent draining the faucet.
- **Configurability** – allow administrators to tune dispense amounts and cooldown periods to match network conditions.

## Core Implementation

The faucet is implemented as a thread‑safe Go service. It maintains an internal balance, a fixed dispense amount, a cooldown duration, and a record of the last request per address【F:faucet.go†L10-L17】. Requests are processed atomically to ensure consistent state updates【F:faucet.go†L27-L37】. Administrators may update operational parameters at runtime, enabling dynamic governance of distribution policies【F:faucet.go†L48-L53】.

### Data Model and Concurrency

The service stores requester timestamps in a map keyed by address to enforce rate limits【F:faucet.go†L11-L16】. Incoming requests acquire a mutex before checking balance and cooldown status, preventing race conditions under heavy load【F:faucet.go†L27-L34】. Error messages surface distinct causes such as an empty balance or premature requests, aiding automated clients in retry logic【F:faucet.go†L30-L35】.

A mirrored implementation exists under the `core` package for integration with higher‑level services. It exposes the same request and configuration semantics, internally sourcing timestamps to simplify call sites【F:core/faucet.go†L28-L41】.

## Command‑Line Interface

Operational control is exposed through the Synnergy CLI. The `faucet` command group, built on the Cobra framework, provides subcommands to initialise the faucet, request tokens, inspect remaining balance, and update configuration parameters【F:cli/faucet.go†L13-L83】. Each subcommand emits human‑readable responses and non‑zero exit codes on failure, enabling automation scripts to handle retries or escalate errors.

### Key Subcommands

| Command | Description |
|---------|-------------|
| `faucet init --balance <n> --amount <n> --cooldown <d>` | Deploys a faucet instance with an initial balance, dispense amount, and cooldown duration. |
| `faucet request <address>` | Attempts to send the configured amount to the target address. Errors are surfaced when the faucet is uninitialised or the cooldown has not elapsed【F:cli/faucet.go†L39-L47】. |
| `faucet balance` | Displays the remaining faucet balance. |
| `faucet config --amount <n> --cooldown <d>` | Updates the dispense amount and cooldown period for subsequent requests, allowing real‑time tuning of output rates【F:cli/faucet.go†L67-L82】. |

Default flag values (balance `1000`, amount `1`, cooldown `1m`) accelerate local onboarding by providing sensible startup parameters【F:cli/faucet.go†L30-L33】.

Administrators may reconfigure live instances using `faucet config`, which immediately updates the in‑memory limits without service restarts【F:cli/faucet.go†L67-L83】. Commands guard against uninitialised use and return contextual messages such as "faucet not initialised" or detailed error strings from the service layer【F:cli/faucet.go†L39-L47】【F:cli/faucet.go†L55-L60】.

### Automation and Scripts

A helper script, `faucet_fund.sh`, simplifies repeated funding requests from shell environments【F:cmd/scripts/faucet_fund.sh†L1-L7】. It validates that a destination address is supplied and then proxies the request through the CLI. An accompanying deployment scaffold (`scripts/deploy_faucet_contract.sh`) outlines the planned automation for rolling out faucet contracts in regulated environments and exposes a `--help` flag for future extensibility【F:scripts/deploy_faucet_contract.sh†L1-L16】.

## Smart Contract Implementations

For on‑chain scenarios, Neto Solaris supplies a reference Solidity contract `TokenFaucet.sol`. It exposes immutable token binding, owner identification, a configurable drip amount, and a 24‑hour wait period between requests【F:smart-contracts/solidity/TokenFaucet.sol†L10-L27】. Anyone may deposit additional liquidity, while only the owner may modify the drip rate or withdraw balances【F:smart-contracts/solidity/TokenFaucet.sol†L31-L45】.

An equivalent model is provided in Rust for WebAssembly targets, enabling cross‑platform faucet deployments with consistent semantics. The Rust struct records the owner, per‑call drip amount, and wait interval, tracking last request timestamps in a `HashMap` to enforce cooldowns【F:smart-contracts/rust/src/token_faucet.rs†L1-L23】. The compiled artefact (`smart-contracts/token_faucet.wasm`) allows the same logic to execute within Synnergy's WASM runtime.

### Virtual Machine and Gas Integration

Faucet operations are first‑class citizens within the Synnergy Virtual Machine. Opcode mappings register creation, requests, balance queries, and configuration updates, ensuring transactions can invoke faucet functions deterministically【F:snvm._opcodes.go†L596-L644】. Both the `core` and top‑level faucet APIs receive dedicated opcode numbers (`0x000245`–`0x000275`), giving tooling a stable interface for contract invocation【F:snvm._opcodes.go†L596-L644】. Corresponding gas prices are tracked so dashboards can forecast the cost of deploying faucet templates alongside other smart‑contract modules【F:gas_table.go†L18-L29】.

## Testing Coverage

Automated tests validate faucet interactions through the CLI by deploying the token faucet template and verifying its registration within the contract registry【F:tests/contracts/faucet_test.go†L43-L60】. The harness captures both stdout and stderr to assert that the deployment address is surfaced and subsequently listed, establishing end‑to‑end coverage from command invocation to registry persistence【F:tests/contracts/faucet_test.go†L14-L40】. These tests also ensure gas schedules are loaded prior to execution, aligning with the broader Synnergy virtual machine framework. Additional unit tests for the in‑memory Go service are scaffolded and will expand coverage as enterprise requirements mature【F:faucet_test.go†L1-L6】.

## Operational Workflow

1. **Initialisation** – create the faucet with a defined balance, dispense amount, and cooldown.
2. **Funding** – deposit test tokens if using the on‑chain implementation.
3. **Request Cycle** – participants invoke `faucet request` to receive the drip amount. The system records timestamps per address to enforce waiting periods.
4. **Monitoring** – administrators periodically inspect `faucet balance` and adjust parameters through `faucet config` as required.

## Security Considerations

- **Rate Limiting** – cooldown tracking mitigates automated draining attempts.
- **Balance Safeguards** – requests are denied when the faucet lacks sufficient funds, preserving deterministic behaviour.
- **Owner Controls** – smart contract variants restrict administrative actions (parameter changes, withdrawals) to the contract owner, while permitting third‑party deposits to replenish reserves【F:smart-contracts/solidity/TokenFaucet.sol†L31-L45】.
- **Clear Failure Modes** – the Go service returns discrete error strings such as "faucet empty" or "cooldown period not met" so API gateways can differentiate abuse from legitimate exhaustion【F:faucet.go†L30-L35】.
- **Mutable Rate Limits** – live configuration updates via the CLI allow operators to throttle traffic during incidents without redeploying services【F:cli/faucet.go†L67-L83】.

## Enterprise Deployment Considerations

- **Persistence and Scaling** – the reference implementation maintains state in memory; production deployments should persist balances and request logs in a replicated datastore to survive restarts and enable horizontal scaling.
- **Observability** – integration with the Synnergy telemetry stack allows operators to export metrics on request rates and error counts, supporting proactive capacity planning.
- **Access Mediation** – public endpoints should layer CAPTCHA or signature schemes on top of the default cooldown logic to resist sybil attacks and scripted abuse.
- **Governance and Audit** – immutable request logs and administrator actions should feed central audit trails so compliance teams can reconstruct allocation histories and enforce policy boundaries.
- **Disaster Recovery** – backups of faucet state and configuration should be replicated across regions to ensure continuity during data‑centre failures.
- **Operational Boundaries** – when multiple teams share infrastructure, deploy isolated faucet instances or namespace controls to prevent cross‑tenant interference.

## Future Enhancements

Road‑map items include integrating CAPTCHA or signature verification for public endpoints, metric export for observability, and optional identity checks through Synnergy's `syn900` modules.

## Conclusion

The Neto Solaris faucet provides a configurable and secure mechanism for distributing test tokens across the Synnergy Network. By combining lightweight Go services, CLI tooling, and smart‑contract implementations, the system ensures developers and stakeholders can experiment safely while maintaining governance over token issuance.

