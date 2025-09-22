# Server Setup and Testing Guide

This guide explains how to deploy the Synnergy blockchain on a fresh Linux server and run its full test suite. The steps use the repository root `/workspace/Synnergy.v.1` but can be adapted for any path.

## Stage 82 Readiness Checks

Stage 82 streamlines server bring-up by folding the enterprise bootstrap flow
into the recommended setup. After building the CLI, run `synnergy orchestrator
bootstrap` to initialise the VM, synchronise consensus relayers, seal the wallet
and perform a ledger audit. The command returns JSON diagnostics that mirror the
JavaScript control panel—wallet seal state, consensus relayer count, authority
role distribution and gas synchronisation timestamp—so automation can block
deployments if any prerequisite fails. Logging destinations, formats and levels
are validated via `configureLogging`, while `registerEnterpriseGasMetadata`
ensures all nodes share the same gas schedule before integration tests execute.

## Prerequisites

Prepare a server with the following software:

- Linux distribution with `bash`
- Git
- Go 1.20 or newer
- Optional: Docker and Docker Compose for containerised runs

Update system packages and install prerequisites (Debian/Ubuntu):

```bash
sudo apt update
sudo apt install -y git curl build-essential docker.io docker-compose
```

## Clone the Repository

```bash
git clone https://github.com/your-org/Synnergy.v.1.git
cd Synnergy.v.1
```

## Build the Command Line Interface

The repository includes helper scripts for installation and builds:

```bash
./setup_synn.sh        # install Go and compile the synnergy binary
# or for the full environment with optional tooling
./Synnergy.env.sh
```

To build manually:

```bash
cd synnergy-network
go mod tidy
GOFLAGS="-trimpath" go build -o synnergy ./cmd/synnergy
```

## Initialise a Local Node

```bash
cd synnergy-network
./synnergy ledger init --path ./ledger.db
./synnergy network start &
```

Open a new shell to interact with the node via the CLI. Useful scripts for multi-node setups are available in `scripts/` such as `devnet_start.sh` and `testnet_start.sh`.

## Run the Test Suite

Execute all Go tests to verify the environment:

```bash
go test ./...
```

Running tests may take several minutes and requires network access for some packages.

## Optional: Docker Deployment

A containerised node can be built with the included Dockerfile:

```bash
docker build -t synnergy .
docker run --rm -p 30303:30303 synnergy
```

For multi-service setups refer to `docker-compose.yml`.

## Next Steps

- Explore additional guides under `docs/guides/`
- Use `scripts/testnet_start.sh` to launch a configurable test network
- Consult `README.md` for more development details

