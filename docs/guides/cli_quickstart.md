# CLI Quickstart

The Synnergy CLI allows interaction with the blockchain for common tasks.

## Installation
Ensure the project is built:
```bash
make build
```

## Common Commands
- `synnergy wallet create` – generate a new wallet
- `synnergy tx send` – send a transaction
- `synnergy node status` – display node synchronization status
- `synnergy mining mine-until <data> <prefix> --timeout <sec>` – hash input until the prefix is found or a timeout elapses
- `synnergy peer count` – show the number of known peers
- `synnergy charity_pool --json registration <addr>` – view charity registration info as JSON
- `synnergy charity_mgmt donate <from> <amount>` – donate tokens to the charity pool
- `synnergy coin --json info` – inspect monetary parameters
- `synnergy cross_chain_bridge deposit <bridge> <from> <to> <amount> --json` – lock assets for bridging with structured output
- `synnergy authority vote <voter> <candidate> --pub <hex> --sig <hex>` – cast a signed vote for an authority node
- `synnergy authority_apply vote <voter> <id> <approve> --pub <hex> --sig <hex>` – cast a signed vote on an authority application
- `synnergy bankinst register <name> --pub <hex> --sig <hex>` – enrol a bank institution with a signed request
- `synnergy bank_index add <id> <type>` – record a bank node in the index
- `synnergy basenode dial <addr> --pub <hex> --sig <hex>` – connect to a peer with signature validation
- `synnergy bioauth enroll <addr> <data> <pubHex>` – register a biometric template with an Ed25519 key

## Stage 85 – SYN3800 grant lifecycle

Stage 85 hardens the grant distribution CLI with wallet-gated flows, persistent
state and enterprise telemetry coordinated by the Grant Orchestrator.

- `synnergy syn3800 create <beneficiary> <name> <amount> --wallet <path> --password <pw> [--authorizer <wallet:path>]` – create an orchestrated grant with signed provenance and optional initial authorisers.
- `synnergy syn3800 authorize <id> --wallet <path> --password <pw>` – enrol an
  additional wallet using encrypted credentials.
- `synnergy syn3800 release <id> <amount> [note] --wallet <path> --password <pw>` – disburse funds with signature validation and
  automatic status transitions.
- `synnergy syn3800 audit <id>` / `synnergy syn3800 status` – stream JSON event
  logs and lifecycle totals for observability tooling and the `/grants` web
  console.

## Help
Run `synnergy --help` or `synnergy <command> --help` for more details.

## Global Flags
The root command provides options that apply to every sub-command:

- `--config` – path to a CLI configuration file
- `--log-level` – log verbosity (`info` or `debug`)

## Regulatory Operations
- `synnergy regulator add <id> <jurisdiction> <description> <max>` – register a transaction rule
- `synnergy regnode approve <from> <amount> --wallet <file> --password <pw>` – sign with a wallet and validate a transaction against regulations
- `synnergy regnode logs <addr>` – view recorded flags for an address
- `synnergy regnode audit <addr>` – check whether an address has been flagged and view its logs
- `synnergy regnode flag <addr> <reason>` – manually flag an address for review (reason required)

## Replication
- `synnergy replication start --json` – launch the replication subsystem
- `synnergy replication status --json` – check whether replication is running

## Rollups
- `synnergy rollups submit <tx...> --json` – create a new rollup batch
- `synnergy rollupmgr pause --json` – pause the rollup aggregator

## Sharding
- `synnergy sharding leader set <id> <addr> --json` – assign a leader to a shard
- `synnergy sharding map --json` – list shard-to-leader mappings
