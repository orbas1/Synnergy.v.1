# Token Guide

This guide explains how tokens are organised and implemented in the Synnergy Network.
It summarises the core token abstractions, the registry used to manage them and the
built-in token standards provided by the project.

## Package layout

Token code lives under `internal/tokens`.  Key files are:

- `base.go` – defines the `Token` interface and a reusable `BaseToken` implementation.
- `index.go` – provides a registry to create, lookup and list token instances.
- Additional files such as `syn10.go` or `syn4200_token.go` implement individual
  token standards.

The command line wrappers under `cli/` expose similar functionality for interacting
with tokens via CLI commands. For sensitive assets these commands can be wrapped
by the `BiometricSecurityNode`, requiring signed biometric proofs before
transfers or mint operations are accepted.

## Base token

Every token implements the `Token` interface which exposes common behaviour:
`ID`, `Name`, `Symbol`, `Decimals`, `TotalSupply`, `BalanceOf`, `Transfer`,
`TransferFrom`, `Mint`, `Burn`, `Approve` and `Allowance`.
The `BaseToken` struct in `base.go` provides a straightforward implementation of
this interface and handles supply accounting and allowance management.

To define a new fungible token, compose your struct with `*BaseToken` and add any
custom state or methods required by the standard.

## Token registry

The registry in `index.go` assigns unique identifiers and keeps a map of active
tokens.  It can return metadata for a single token or list all registered tokens.
Projects typically create a registry at start up and register each token as it is
initialised.

```go
reg := tokens.NewRegistry()
id  := reg.NextID()
my  := tokens.NewSYN20Token(id, "My Token", "MYT", 18)
reg.Register(my)
```

## Built-in token standards

The repository includes several reference token implementations that demonstrate
how specialised assets can be modelled on top of the base abstractions:

| Token file | Purpose |
|------------|---------|
| `syn10.go` | Central bank digital currency with issuer and exchange-rate controls. |
| `syn12.go` | Tokenised treasury bill instrument including maturity and discount fields. |
| `syn20.go` | Fungible token with pause and address freeze capabilities. |
| `syn70.go` | Registry for in‑game assets, supporting attributes and achievements. |
| `syn200.go` | Carbon credit projects with issuance, retirement and verification tracking. |
| `syn223_token.go` | Secure transfer token that enforces whitelist and blacklist rules. |
| `syn300_token.go` | Governance token supporting delegation and on‑chain proposals. |
| `syn845.go` | Debt token registry for recording loans and repayments. |
| `syn1000.go` & `syn1000_index.go` | Reserve‑backed stablecoin and index for managing multiple instances. |
| `syn1100.go` | Healthcare record storage with access control lists. |
| `syn2369.go` | Virtual item registry for metaverse assets. |
| `syn2500_token.go` | DAO membership registry with voting power metadata. |
| `syn2600.go` | Investor tokens that record share ownership and return distributions. |
| `syn2800.go` | Life‑insurance policies with premium and claim management. |
| `syn2900.go` | General insurance policies and claim handling. |
| `syn3400.go` | Foreign‑exchange pair registry with rate updates. |
| `syn3500_token.go` | Fiat‑pegged currency token with mint and redeem operations. |
| `syn3700_token.go` | Index token that aggregates multiple assets by weight. |
| `syn4200_token.go` | Charity campaign token used for tracking donations. |
| `syn4700.go` | Legal document token recording parties, signatures and dispute status. |

These examples can be used as templates when designing new token types.

## Native coin

The native asset **Synthron** is defined in `core/coin.go`. It exposes helper
functions for block rewards, supply tracking and staking economics. The `coin`
CLI command surfaces these utilities and provides JSON output for web
dashboards. Minting beyond the capped supply is prevented by the
`CentralBankingNode` which enforces checks against `RemainingSupply`.

## AI model access tokens

Stage 2 introduces AI-enhanced contracts and audit modules. Projects can mint
tokens that meter access to deployed models or audit logs, leveraging the
standard registry and `BaseToken` abstractions.

## Institutional and governance tokens

Stage 3 adds authority and banking node modules. Tokens representing voting
power or institutional membership can be combined with the `AuthorityNodeRegistry`
and `BankInstitutionalNode` to restrict transfers or confer special privileges.
The CLI outputs JSON, allowing web interfaces to track tokenised governance
rights alongside registered financial institutions.

## Creating a new token

1. Decide whether the token is fungible.  For fungible tokens embed `*BaseToken`.
2. Obtain a unique ID from the registry and instantiate your token.
3. Register the token so that it can be retrieved by ID or symbol.
4. Extend with any domain‑specific methods or data structures.

## Testing

Unit tests for the token package can be executed with:

```bash
go test ./internal/tokens
```

Running the full repository test suite will also exercise all token
implementations.
