# Synnergy Configuration Guide

This document describes how to configure a Synnergy node using the YAML files in `cmd/config/`. The configuration loader merges `default.yaml` with an optional environment specific file identified by the `SYNN_ENV` variable. Applications should call `config.LoadFromEnv()` which implements this behaviour. All values can be overridden by environment variables thanks to Viper's `AutomaticEnv` integration and the project `.env` file.

## Stage 82 Enhancements

Stage 82 hardens configuration handling across the CLI, virtual machine and web
interfaces. The new `configureLogging` helper validates log levels, formats and
output targets before the runtime starts, preventing misconfigured nodes from
silently discarding telemetry. `registerEnterpriseGasMetadata` loads the
documented gas schedule, enforces the Stage 78 enterprise opcodes and attaches
human-readable metadata so CLI automation, the JavaScript control panel and the
function web all present identical gas pricing. These checks execute during the
bootstrap sequence that powers `synnergy orchestrator bootstrap`, ensuring wallet
sealing, consensus relayer registration and ledger audits operate against a
validated configuration stack.

## Loading Sequence

1. `default.yaml` is loaded first and establishes sane defaults for development.
2. If `SYNN_ENV` is set (for example `prod` or `bootstrap`) a file with the same name is merged in.
3. Environment variables are then applied. Keys use dot notation matching the YAML structure (e.g. `network.max_peers`).
4. The resulting configuration is unmarshaled into the `Config` struct defined in `pkg/config/config.go` and made available as `config.AppConfig`.

Running the CLI without `SYNN_ENV` uses only the default settings which are safe for local testing.

## Configuration Structure

Every YAML file shares the same schema shown below. Each section is described in detail in the following subsections.

```yaml
network:
  id: synnergy-mainnet
  chain_id: 1215
  max_peers: 50
  genesis_file: config/genesis.json
  rpc_enabled: true
  p2p_port: 30303
  listen_addr: "/ip4/0.0.0.0/tcp/4001"
  discovery_tag: synnergy-mesh
  bootstrap_peers: []

consensus:
  type: pos
  block_time_ms: 3000
  validators_required: 3

vm:
  max_gas_per_block: 8000000
  opcode_debug: false

storage:
  db_path: ./data/db
  prune: true

logging:
  level: info
  file: logs/synnergy.log
```

### `network`

| Field | Description |
|-------|-------------|
| `id` | Human friendly identifier for the network. Appears in logs and the genesis file. |
| `chain_id` | Integer identifier used when signing transactions. Nodes on different chain IDs will refuse to connect. |
| `max_peers` | Maximum number of simultaneous peer connections. |
| `genesis_file` | Path to the genesis block JSON. This file defines the initial state, balances and validator set. |
| `rpc_enabled` | Enables the JSON-RPC server when `true`. |
| `p2p_port` | TCP port for libp2p communications. |
| `listen_addr` | Full multiaddress the node listens on. |
| `discovery_tag` | Rendezvous string used during peer discovery. |
| `bootstrap_peers` | List of static peers to dial on startup. |

### `consensus`

| Field | Description |
|-------|-------------|
| `type` | Name of the consensus algorithm (`pos` by default). |
| `block_time_ms` | Target block interval in milliseconds. |
| `validators_required` | Minimum number of validators that must sign a block. |

### `vm`

Configuration for the embedded virtual machine.

| Field | Description |
|-------|-------------|
| `max_gas_per_block` | Maximum total gas that can be consumed by all transactions in a block. |
| `opcode_debug` | When `true`, logs every opcode executed—useful for debugging contracts. |

### `storage`

| Field | Description |
|-------|-------------|
| `db_path` | Directory where the node stores its local database files. |
| `prune` | If enabled, old states are pruned to save disk space. |

### `logging`

| Field | Description |
|-------|-------------|
| `level` | Log verbosity (`debug`, `info`, `warn`, `error`). |
| `file` | Optional file path for log output. If empty, logs go to stdout. |

## Environment Files and Overrides

The repository contains a `.env` file with additional variables such as API endpoints and secrets. `config.Load` reads these values automatically so they can override YAML fields. Typical variables include `API_BIND`, `GOVERNANCE_API_ADDR`, and `JWT_SECRET`. When running in production make sure to provide secure values for any secrets.

## Example Configurations

- **default.yaml** – Used for local development. RPC is enabled and the node listens on all interfaces.
- **prod.yaml** – Starting point for production deployments. Values generally mirror `default.yaml` but can be customised by setting `SYNN_ENV=prod`.
- **bootstrap.yaml** – Template for a dedicated bootstrap node. It sets `discovery_tag` to `synnergy-bootstrap`, typically raises `max_peers`, and omits `bootstrap_peers` so other nodes connect to it first.

## Genesis Block

The genesis file referenced by `network.genesis_file` establishes the very first block in the chain. The sample `genesis.json` included in this repository looks like:

```json
{
  "genesis_time": "2025-07-09T12:00:00Z",
  "chain_id": "synnergy-mainnet",
  "initial_balances": {
    "0xABC...": 1000000000,
    "0xDEF...": 500000000
  },
  "validators": [
    { "address": "0x123...", "stake": 1000000 }
  ]
}
```

At start-up the ledger reads this file and creates the first block with the specified balances and validator set. Adjust the addresses and amounts to suit your deployment. The `chain_id` must match the value under `network.chain_id`.

## Creating a New Configuration

1. Copy `default.yaml` to a new file name such as `staging.yaml`. An example `staging.yaml` is included in this repository for release candidate testing.
2. Modify any fields required for your environment—commonly the `genesis_file` path and `bootstrap_peers`.
3. Set `SYNN_ENV=staging` before running the `synnergy` binary.
4. Optionally create or edit a `.env` file to supply secrets and service endpoints.
5. Use `../scripts/staging_validate.sh` to spin up the staging network and run a sample token transfer.

## Conclusion

The configuration system is intentionally simple: a layered YAML approach backed by environment variables. By editing these files and the genesis block you can tailor a Synnergy deployment for local testing, dedicated bootstrap nodes or a full production network.

## API Version

The configuration loader API is currently versioned as **v0.1.0**. Future changes will follow semantic versioning guidelines.
