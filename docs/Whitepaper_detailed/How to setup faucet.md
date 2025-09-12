# How to Setup the Faucet

Neto Solaris provides a built-in faucet to distribute small amounts of SYN tokens to developers and testers. The faucet enforces rate limits to prevent abuse while allowing quick access to funds for experimentation. This guide explains how to initialise, operate, and extend the faucet using the Synnergy Network toolset.

## Architecture Overview

The faucet maintains an internal balance and tracks the timestamp of each request. A thread-safe mutex guards shared state while a map records the last time each address drew funds. When a user requests tokens, the faucet verifies that sufficient balance remains and the caller has waited longer than the configured cooldown interval before dispensing funds【F:core/faucet.go†L9-L41】. Administrators can query the remaining balance and update the dispense amount or cooldown without redeploying the service【F:core/faucet.go†L44-L56】.

### Virtual Machine Opcodes

Faucet primitives are surfaced as deterministic opcodes in the Synnergy Virtual Machine. The opcode table assigns unique identifiers to creation, request, balance and configuration operations, enabling contracts or higher-level tooling to invoke faucet logic directly from the VM【F:snvm._opcodes.go†L596-L644】.

### Cross-Language Support

In addition to the Go implementation, the faucet model is available in Rust and compiled to WebAssembly for lightweight deployments. The Rust module serialises owner, drip amount and wait-time fields and enforces per-account cooldowns using a hash map of last-request timestamps【F:smart-contracts/rust/src/token_faucet.rs†L1-L24】. A Solidity contract template offers comparable functionality for ERC-20 tokens on EVM-compatible chains【F:smart-contracts/solidity/TokenFaucet.sol†L9-L45】.

### CLI Integration

The `synnergy` command-line interface exposes subcommands to manage the faucet. Initialisation sets the starting balance, the per-request amount, and the cooldown duration. Subsequent subcommands allow token requests, balance inspection, and configuration updates【F:cli/faucet.go†L14-L83】. For repetitive development tasks, the repository bundles a helper script that wraps the request command and provides basic argument validation【F:cmd/scripts/faucet_fund.sh†L1-L7】.

### Programmatic Usage

Services may embed the faucet directly using the Go APIs. A faucet instance tracks balance, dispense policy and last request time in a mutex-protected structure:

```go
f := core.NewFaucet(100, 10, time.Minute)
amt, err := f.Request("addr")
```

Requests fail with clear error messages when funds are exhausted or a caller triggers the cooldown【F:core/faucet.go†L29-L37】. Operators can adjust policy at runtime without redeploying by invoking `Configure` on the public library variant【F:faucet.go†L48-L53】.

### Error Handling

Typical errors returned by the faucet include:

- `cooldown active` – the address has requested funds too recently【F:core/faucet.go†L33-L34】.
- `insufficient faucet balance` – available tokens are below the dispense amount【F:core/faucet.go†L36-L37】.
- `faucet empty` or `cooldown period not met` when using the package-level faucet【F:faucet.go†L30-L34】.

Applications should surface these errors to end users and log them for auditability.

## Prerequisites

- **Synnergy CLI** – build the binary from `cmd/synnergy` or obtain a release build.
- **Node Access** – connectivity to a Synnergy Network node capable of accepting faucet transactions.
- **Funding Account** – an address that will supply the initial faucet balance.

## Initialising the Faucet

Create a faucet instance with an initial balance and policy settings:

```bash
synnergy faucet init --balance 1000 --amount 5 --cooldown 2m
```

- `--balance` – total tokens available for distribution.
- `--amount` – tokens dispensed per request.
- `--cooldown` – minimum interval between requests from the same address.

On success the CLI prints `faucet initialised`, indicating the structure is ready to serve requests.

## Requesting Tokens

Developers request tokens by supplying a destination address:

```bash
synnergy faucet request <address>
```

If the faucet is uninitialised or the cooldown has not elapsed, the command returns an error. Successful requests display the number of tokens dispensed.

## Checking Faucet Balance

Monitor remaining funds to ensure the faucet stays stocked:

```bash
synnergy faucet balance
```

The command prints the current balance, enabling operators to top up the faucet before depletion.

## Updating Configuration

Adjust the dispense amount or cooldown without restarting the service:

```bash
synnergy faucet config --amount 10 --cooldown 5m
```

This command calls `UpdateConfig`, replacing the previous policy values with the supplied parameters【F:core/faucet.go†L51-L56】【F:cli/faucet.go†L67-L83】.

## Deploying the Token Faucet Contract

For on-chain distribution, Synnergy ships a reusable `TokenFaucet` smart contract. The contract holds an ERC-20 token balance and restricts each address to one request per wait period【F:smart-contracts/solidity/TokenFaucet.sol†L9-L28】. Owners can deposit additional tokens, adjust the drip amount, or withdraw excess funds as needed【F:smart-contracts/solidity/TokenFaucet.sol†L31-L41】.

Deploy the template via the CLI:

```bash
synnergy contracts deploy-template --name token_faucet
```

Record the returned address and transfer tokens to it using the `deposit` function. Clients then call `requestTokens` to receive the configured drip amount.

### Opcode and Gas Costs

Faucet operations are implemented as first-class opcodes in the Synnergy Virtual Machine. Deploying a new faucet consumes 500 gas, while requests, balance queries, and configuration updates have zero runtime cost, enabling inexpensive developer onboarding【F:docs/Whitepaper_detailed/guide/opcode_and_gas_guide.md†L2699-L2712】. Gas table entries register predictable pricing for contract templates, including the faucet, so costs remain stable across releases【F:gas_table.go†L23-L31】.

## Automation and Testing

Continuous integration pipelines deploy and verify the faucet contract using the same CLI commands documented above. The contract test harness invokes the `deploy-template` subcommand, captures the returned address and asserts that the deployment is registered in the contract listing【F:tests/contracts/faucet_test.go†L43-L61】. Unit tests also validate cooldown behaviour and configuration updates in the core library【F:core/faucet_test.go†L8-L23】. A placeholder script exists for automating contract deployment in shell environments, offering a starting point for bespoke CI workflows【F:scripts/deploy_faucet_contract.sh†L1-L17】.

## Best Practices

- **Maintain Sufficient Balance** – monitor usage and replenish the faucet to avoid service interruptions.
- **Tune Cooldowns** – longer intervals reduce abuse on public testnets; shorter intervals facilitate rapid prototyping in isolated environments.
- **Secure the Funding Account** – only authorised personnel should control the account that tops up or withdraws from the faucet.
- **Audit Smart Contracts** – when using the on-chain faucet, review the contract code and access controls before deployment.
- **Log Requests** – capture request metadata and outcomes to support auditing and anomaly detection.

## Enterprise Deployment Considerations

Large organisations often host faucets for multiple internal teams or external partners. To sustain heavy traffic:

- **Horizontal Scaling** – run multiple faucet instances behind a load balancer and share a backing account to prevent bottlenecks. The mutex-protected design ensures each instance processes requests safely even under concurrency【F:core/faucet.go†L9-L41】.
- **Observability** – integrate request and balance metrics with existing monitoring suites to track utilisation and forecast top-ups.
- **Access Control** – restrict who can adjust drip parameters or withdraw funds by leveraging owner-only functions in the contract template【F:smart-contracts/solidity/TokenFaucet.sol†L33-L41】.
- **Disaster Recovery** – back up faucet configuration and funding accounts, and script redeployment to minimise downtime.
- **Scriptable Workflows** – pair the CLI with automation tools or use the provided shell scripts to seed test accounts during continuous delivery pipelines【F:cmd/scripts/faucet_fund.sh†L1-L7】.

## Conclusion

The Synnergy faucet, engineered by Neto Solaris, accelerates development by providing rate-limited test tokens through both CLI utilities and smart-contract templates. By following the steps above, operators can reliably configure faucets for internal labs or public testnets, ensuring seamless onboarding and controlled token distribution across the Synnergy Network.

