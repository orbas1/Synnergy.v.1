# Synnergy


Synnergy is an enterprise production blockchain written in Go. This repository contains the command line applications, core packages, GUI front‑ends and example smart contracts used to simulate a full network. The code is primarily intended for research and learning. For the vision and background see [`synnergy-network/WHITEPAPER.md`](synnergy-network/WHITEPAPER.md).
## Documentation

All guides and architecture decision records are located under the `docs/` directory and are built with MkDocs. Run `mkdocs serve` for a live server or `mkdocs build` to generate the static site. For deploying and exercising the full test suite on a server, see [`docs/guides/server_setup_guide.md`](docs/guides/server_setup_guide.md).


## Features
- Pluggable node types for mining, staking, authority, regulatory, watchtower and other roles.
- Extensive CLI built with [Cobra](https://github.com/spf13/cobra) located under `cli/`.
- AI modules for model management, inference analysis and anomaly detection.
- AI-enhanced contract registry with on-chain audit logging.
- Cross‑chain bridge and protocol support.
- Role‑based access control with validated address utilities.
- Biometric security backed by ECDSA signatures for sensitive node operations.
- YAML based configuration for development, test and production environments.
- Web interfaces for wallets, explorers and marketplaces under `GUI/`.
- JSON emitting CLI commands for authority and institutional banking nodes to ease integration with web dashboards.
- Built-in ledger synchronization and snapshot compression utilities with
  central bank and charity pool modules enforcing capped supply and donation
  tracking.
- Structured JSON logging with pluggable backends for compliance, connection and consensus modules.
- Stage 7 adds a unified errors package and OpenTelemetry tracing for consensus and contract management components, improving diagnostics across the CLI and services.
- Stage 8 introduces a contract registry and cross‑chain transaction managers with full CLI access and gas‑priced opcodes for deterministic execution.
- Stage 9 adds DAO governance, custodial nodes and cross-consensus network tooling with gas-priced opcodes and CLI support.
- Stage 11 introduces a context-aware virtual machine and sandbox manager, enabling contract execution with timeouts and full lifecycle control through the CLI. Inactive sandboxes can be purged automatically via `synnergy sandbox purge` to reclaim resources.

## Repository layout
```
cmd/          Command line entry points (e.g. `cmd/synnergy`)
cli/          CLI command implementations
core/         Blockchain runtime modules (consensus, networking, VM, …)
configs/      Default configuration files (`dev.yaml`, `test.yaml`, `prod.yaml`)
internal/     Shared utilities and configuration loading
GUI/          Web front‑end projects
docs/         Guides and MkDocs documentation sources
scripts/      Helper scripts and automation
pkg/          Reusable packages and libraries
```
Additional extensions live under `node_ext/` and `internal/`.

## Getting started
### Prerequisites
- Go 1.24 or newer
- (optional) Docker for containerised builds
- (optional) Node.js for GUI projects

### Build the CLI
```
go build ./cmd/synnergy
# or
make build
```
The resulting binary named `synnergy` is written to the repository root.

### Run
Select a configuration file from `configs/` or provide your own via the `SYN_CONFIG` environment variable.
```
export SYN_CONFIG=configs/dev.yaml
./synnergy --help
./synnergy network start
```
Helper scripts in `scripts/` can launch multi‑node devnets or testnets for experimentation.

## Testing and security checks
Run the unit tests and static analysis tools before submitting changes:
```
go test ./...
make security   # runs staticcheck, gosec and govulncheck
```

## Documentation
Project guides and architecture notes live under `docs/` and are built with [MkDocs](https://www.mkdocs.org/):
```
make docs        # build the static site into site/
make docs-serve  # serve the documentation locally
```

## Contributing
Development follows the staged workflow described in [AGENTS.md](AGENTS.md).  Keep pull requests focused, format code with `go fmt`, verify with `go vet`, and run `go build` and `go test` on the packages you touch.

## License
Synnergy is provided for research and educational purposes.  Third‑party dependencies retain their original licenses.
