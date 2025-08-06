# Synnergy

Synnergy is a modular blockchain network written in Go.  This repository contains the command line interface, core runtime packages, node implementations, documentation site and example tooling used to experiment with the protocol.  The code base showcases advanced concepts such as AI‑assisted contracts, cross‑chain operations and specialised node roles.

## Features
- Pluggable node types for mining, staking, authority, regulatory, watchtower and other roles.
- Extensive CLI built with [Cobra](https://github.com/spf13/cobra) located under `cli/`.
- AI modules for model management, inference analysis and anomaly detection.
- Cross‑chain bridge and protocol support.
- YAML based configuration for development, test and production environments.
- Web interfaces for wallets, explorers and marketplaces under `GUI/`.

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
