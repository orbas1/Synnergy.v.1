# Charity Pool Guide

The charity module coordinates donations and community voting for registered
charitable organisations. Funds are deposited into a shared pool and, on a
regular cycle, distributed to the charities that receive the most support from
token holders.

## Stage 82 Enterprise Integration

Stage 82 connects the charity workflow to the enterprise bootstrap sequence so
philanthropic campaigns inherit the same guarantees as the core ledger. The new
`synnergy orchestrator bootstrap` command spins up the heavy virtual machine,
seals the orchestrator wallet, registers consensus relayers and performs a
ledger audit before any charity transactions are processed. The JavaScript
control panel now surfaces the orchestrator's wallet seal status, consensus
relayer count, authority role distribution and the most recent gas synchronisation
timestamp so operators can verify compliance before routing donations. These
diagnostics are available programmatically through `/api/orchestrator`, allowing
charity dashboards to pause disbursements if the VM, consensus mesh or wallet
ever fall out of alignment.

Security and privacy controls were also expanded. The runtime now registers
charity opcodes with descriptive gas metadata which feeds both the CLI and web
interfaces, ensuring predictable fees for registration, voting and donation
transactions. VM execution hooks log failed opcode executions with context so
incident responders can isolate faulty automation without restarting the node.
All CLI commands inherit the hardened logging configuration introduced in Stage
82: log destinations and formats are validated at startup, the wallet subsystem
is sealed against missing key material and the consensus relayer whitelist is
persisted, ensuring that charity votes and deposits traverse an authenticated
infrastructure path.

## Features

- **Deposits** – Any address can transfer tokens into the pool account
  `charity_pool`.
- **Registration** – Charities register an address, human‑readable name and a
  category describing their mission.
- **Voting** – Holders of SYN900 identity tokens may vote for a charity once per
  cycle. Votes are persisted to the ledger for auditability.
- **Payouts** – At the end of a cycle the pool can distribute a portion of its
  balance to the top voted charities. The current implementation stores votes
  and registrations; payout logic is left to future work.

## Categories

Charities self‑identify using one of the predefined categories in
`core/charity.go`:

| ID | Category |
|----|----------|
| `1` | HungerRelief |
| `2` | ChildrenHelp |
| `3` | WildlifeHelp |
| `4` | SeaSupport |
| `5` | DisasterSupport |
| `6` | WarSupport |

Additional categories can be added by extending the `CharityCategory` enum.

## CLI Usage

The `synnergy` CLI exposes helper commands under `charity_pool` for interacting
with the module. These commands use an in‑memory ledger so they can be run
without a full node.

### Register a Charity

```bash
synnergy charity_pool register <address> <category> <name>
```

Registers a charity with the given address and human‑readable name. The
`<category>` argument is one of the numeric identifiers listed above.

### Vote for a Charity

```bash
synnergy charity_pool vote <voter_addr> <charity_addr>
```

Records a vote from `voter_addr` to `charity_addr`. Only one vote per cycle is
counted for each voter.

### Deposit to the Pool

```bash
synnergy charity_pool deposit <from> <amount>
```

Transfers tokens from `from` into the shared pool account. Deposits accumulate
until a payout is triggered.

### Query Winners

```bash
synnergy charity_pool winners <cycle>
```

Returns the list of winning charities for a given cycle. The current prototype
returns an empty list until distribution logic is implemented.

## Integrating the Pool

Applications may instantiate the pool in code using:

```go
led := ... // implements core.StateRW
electorate := ... // implements IsIDTokenHolder
cp := core.NewCharityPool(logger, led, electorate, genesisTime)
```

The `StateRW` interface provides ledger access. `electorate` is used to validate
whether a voter holds an identity token. `genesisTime` defines the start of the
first voting cycle.

## Future Work

The current module is intentionally lightweight. Planned enhancements include:

- Automatic calculation of payout cycles and distribution of pool funds.
- On‑chain governance to add or remove categories.
- Richer query APIs for analytics and historical reporting.

Contributions and design proposals are welcome via pull requests.
