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
- `synnergy integration status --format json` – run full-stack diagnostics and retrieve a machine-readable health snapshot including security, scalability, privacy, governance, interoperability and compliance probes
- `synnergy syn3800 create <beneficiary> <name> <amount> --authorizer wallet.json:pass` – register a grant and store the wallet-signed authorizer
- `synnergy syn3800 release <id> <amount> <note> --wallet wallet.json --password pass` – disburse funds with audit logging
- `synnergy syn3900 register <recipient> <program> <amount> --approver wallet.json:pass` – record a benefit programme with an approver signature
- `synnergy syn3900 claim <id> --wallet wallet.json --password pass` – claim benefits using the registered recipient wallet
- `synnergy syn500 grant <addr> --tier 1 --max 5 --window 24h` – grant a service tier with a rolling usage window
- `synnergy syn500 status <addr>` – inspect current tier usage and remaining capacity in JSON form
- `synnergy syn3700 init --name Institutional --symbol IDX --controller wallet.json:pass` – bootstrap the institutional index with a controller wallet
- `synnergy syn3700 add AAA --weight 0.5 --drift 0.1 --wallet wallet.json --password pass` – add a component with drift tolerance enforced by controller signatures
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
