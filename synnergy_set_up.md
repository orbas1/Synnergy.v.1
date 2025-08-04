# Synnergy CLI Enterprise Setup Guide

This guide explains how to prepare a development or testing environment for the
`synnergy` command line interface (CLI) and associated services.  It covers
system requirements, building from source, configuration management and
starting a small multi-service network suitable for enterprise evaluation.

## 1. Prerequisites

Before compiling the CLI ensure the following software is installed:

- **Go 1.20 or newer** – verify with `go version`
- **Git** – for fetching the repository
- **GNU build tools** – `gcc`, `make`, and related packages
- **Docker** *(optional)* – used to run supporting services such as
  databases or reverse proxies

The provided `Synnergy.env.sh` script installs Go and common utilities on
Debian based systems.  Execute the script as root or with `sudo` if you do not
already have a Go toolchain installed.

```bash
./Synnergy.env.sh
```

## 2. Cloning the Repository

```bash
git clone <repo-url>
cd synnergy-network
```

All commands below assume the working directory is `synnergy-network`.

## 3. Building the CLI

1. Download Go module dependencies:
   ```bash
   go mod tidy
   ```
2. Compile the main program:
   ```bash
   go build -o synnergy ./cmd/synnergy
   ```
   The resulting executable `synnergy` is placed in the current directory.
   Optionally move it somewhere in your `PATH` for easier access.

You can also build the CLI using the helper script found in
`cmd/scripts/build_cli.sh` which performs the same steps.

## 4. Configuration Files

Runtime settings are defined in the `cmd/config` directory.  The CLI loads
`default.yaml` by default and merges any environment specific file when the
`SYNN_ENV` variable is set.  For example, `SYNN_ENV=prod` loads `prod.yaml` in
addition to `default.yaml`.

Configuration options include network parameters, consensus settings, virtual
machine limits and logging preferences.  See `cmd/config/config_guide.md` for a
full description of each field.

## 5. Environment Variables

Several commands rely on environment variables in addition to the YAML
configuration files.  The most common variables are:

- `LEDGER_PATH` – path to the ledger database
- `P2P_LISTEN_ADDR` – multiaddr for the libp2p node
- `P2P_BOOTSTRAP` – comma-separated list of bootstrap peers
- `KEYSTORE_PATH` – directory containing node and validator keys
- `SECURITY_API_ADDR` – address of the security daemon

Create a `.env` file with these values or export them in your shell before
running commands.  The `Synnergy.env.sh` script automatically loads a
`synnergy-network/.env` file if present.

## 6. Starting a Local Network

The simplest way to launch all core services is via the example script
`cmd/scripts/start_synnergy_network.sh`.  It compiles the CLI and starts the
network, consensus, replication and virtual machine daemons.

```bash
cd cmd/scripts
./start_synnergy_network.sh
```

To start components manually:

1. **Initialise the ledger**
   ```bash
   ./synnergy ledger init --path ./ledger.db
   ```
2. **Start networking and consensus services**
   ```bash
   ./synnergy network start
   ./synnergy consensus start
   ```
3. **Create a wallet and fund it**
   ```bash
   ./synnergy wallet create --out wallet.json
   ./synnergy coin mint $(jq -r .address wallet.json) 1000
   ```

At this point you may explore command groups such as `contracts` for deploying
WebAssembly code, `tokens` for managing ERC‑20 style assets or `transactions`
for building and broadcasting raw transactions.  See `cmd/cli/cli_guide.md` for
a description of every available command.

## 7. Running Unit Tests

The repository includes a suite of unit tests under `tests/`.  Execute them with

```bash
go test ./...
```

Some tests require auxiliary services or specific environment variables.  Review
the comments in each test file for details.

## 8. Production Considerations

For larger deployments consider the following practices:

- Use `bootstrap.yaml` as a template for dedicated bootstrap nodes
- Configure persistent storage volumes for the ledger and keystore directories
- Run security-related services on isolated hosts with restricted firewall rules
- Enable detailed logging and centralise logs using your preferred observability
  stack
- Maintain separate configuration files for development, staging and production
  environments to avoid accidental cross-contamination

## 9. Additional Resources

- **Command Reference:** `cmd/cli/cli_guide.md`
- **Scripting Examples:** `cmd/scripts/script_guide.md`
- **Core Modules Overview:** `core/module_guide.md`
- **Whitepaper:** `WHITEPAPER.md`

These documents provide further background on the Synnergy architecture and the
individual command groups available within the CLI.

---

Once the CLI is built and the network services are running you are ready to
experiment with the modular blockchain components that make up Synnergy.
