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
- `synnergy system_health snapshot --json` – Stage 75 adds VM metrics, sandbox counts and cross-chain transfer totals to the health payload
- `synnergy warfare commander issue <id>` – mint an Ed25519 commander credential and return the signing key material
- `synnergy warfare command <cmd> --commander <id> --private <hex>` – execute a signed command envelope; omit flags to use the root commander
- `synnergy warfare events --since <seq>` – stream command, logistics and tactical events for dashboards or automation
- `synnergy zero-trust open <id> <hexkey> --owner <name> --meta scope=confidential` – create a channel with metadata and optional retention
- `synnergy zero-trust authorize <id> <participant> <pubkey>` – add a participant public key capable of pushing encrypted payloads
- `synnergy zero-trust events <id> --since <seq>` – fetch channel lifecycle events (open, message, rotate, close) for monitoring pipelines
- `synnergy zero-trust rotate <id> <hexkey>` – rotate the symmetric encryption key while retaining authorised participants
- `synnergy gas list --json` – enumerate opcode pricing with Stage 91 metadata including categories and descriptions for dashboards
- `synnergy highavailability report --json` – Stage 77 signed resilience report covering active node, backups, ledger height and compliance warnings

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

## Stage 100 Automation Helpers

- `scripts/devnet_start.sh` provisions a wallet-backed devnet, starts the
  simple VM, bootstraps networking and adjusts consensus weights in one step.
- `scripts/docker_build.sh` and `scripts/docker_compose_up.sh` provide
  reproducible container builds with retry, logging and optional registry
  pushes.
- `scripts/e2e_network_tests.sh` runs a CLI-driven smoke test that asserts swarm
  membership counts, weight balancing and broadcast health.
