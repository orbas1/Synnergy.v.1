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

## Help
Run `synnergy --help` or `synnergy <command> --help` for more details.

## Global Flags
The root command provides options that apply to every sub-command:

- `--config` – path to a CLI configuration file
- `--log-level` – log verbosity (`info` or `debug`)
- `--stage73-state` – optional path to the Stage 73 snapshot JSON. When supplied the `syn3700`, `syn3800`, `syn3900`, `syn4200_token`, `syn4700` and `syn500` commands persist their state for reuse by the web console, integration tests and the Stage73Orchestrator. The function web defaults to `.synnergy/stage73_state.json`, ensures the parent directory exists and serialises CLI executions so concurrent API calls cannot corrupt the snapshot.

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

### Stage 73 Runtime Workflows – Grants, Benefits and Index Operations

Stage 73 commands now require authenticated wallets and produce machine-readable telemetry so that the web console and automated scripts can operate deterministically. When the CLI is launched with `--stage73-state <path>` the commands read and write a shared snapshot managed by `core.Stage73Store`. The Stage73Orchestrator consumes the same snapshot to expose JSON to the VM and adjust consensus weights based on grant and benefit activity. The `/stage73` browser console reads the identical snapshot and digest through `/api/stage73`, so issuing a command from the UI produces the same persisted state and telemetry as running it locally.

1. **SYN3700 index controller setup**
   ```sh
   synnergy syn3700 init --name "Institutional" --symbol IDX --controller /path/wallet.json:password
   synnergy syn3700 add AAA --weight 0.5 --drift 0.1 --wallet /path/wallet.json --password password
   synnergy syn3700 snapshot
   synnergy syn3700 rebalance --wallet /path/wallet.json --password password
   ```
   All component writes require a controller wallet. Snapshot, controllers, status and audit commands emit JSON that can be consumed directly by dashboards.

2. **SYN3800 grant orchestration**
   ```sh
   synnergy syn3800 create bob research 100 --authorizer /path/ops.json:password
   synnergy syn3800 release 1 40 phase1 --wallet /path/ops.json --password password
   synnergy syn3800 authorize 1 --wallet /path/secondary.json --password pass2
   synnergy syn3800 audit 1
   synnergy syn3800 status
   ```
   Releases and authorization events are persisted with actor metadata and surface in the audit log for compliance review.

3. **SYN3900 benefit approval pipeline**
   ```sh
   synnergy syn3900 register 0xRecipient housing 200 --approver /path/authority.json:password
   synnergy syn3900 claim 1 --wallet /path/beneficiary.json --password password
   synnergy syn3900 approve 1 --wallet /path/authority.json --password password
   synnergy syn3900 list
   synnergy syn3900 status
   ```
   Claims validate that the signing wallet matches the registered recipient and approvals require previously whitelisted approvers.

4. **SYN500 utility limits and telemetry**
   ```sh
   synnergy syn500 create --name Loyalty --symbol LOY --owner treasury --dec 2 --supply 1000
   synnergy syn500 grant bob --tier 1 --max 2 --window 1h
   synnergy syn500 use bob
   synnergy syn500 status bob
   synnergy syn500 telemetry
   ```
   Usage windows automatically reset once the configured interval elapses, and telemetry exposes how many grants remain active before limits are reached.

Each command prints a gas cost prefix followed by JSON output, making it straightforward to route responses into the web UI function web and external orchestration pipelines.
