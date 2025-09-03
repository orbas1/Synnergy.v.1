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
- Stage 12 adds a hex-encoded wallet implementation and specialised warfare and watchtower nodes. The CLI can generate wallets, track logistics for military assets and monitor network health through these modules.
- Stage 13 introduces zero trust data channels with authenticated encryption and regulatory nodes that automatically flag non-compliant transactions.
- Stage 14 consolidates node lifecycle management under an internal `nodes` package with a reusable interface and reference implementations for light, watchtower and logistics nodes.
- Stage 15 expands internal node variants with in-memory forensic, geospatial, historical and elected authority nodes, enabling richer diagnostics and data services across the network.
- Stage 16 introduces a concurrency-safe token registry and base token with micro-benchmarks to track transfer throughput.
- Stage 17 delivers standard token contracts including CBDC, pausable utility and gaming asset tokens. Each implementation is thread-safe and accessible via dedicated CLI modules.
- Stage 18 expands the token library with investor share registries, life and general insurance policies, forex pairs, fiat‑pegged currencies, index funds, charity campaigns and legal document tokens, all validated and manageable through the CLI.
- Stage 19 adds a reserve-backed stablecoin (`SYN1000`) with an index manager and high-precision, thread-safe reserve accounting accessible through dedicated CLI commands.
- Stage 20 introduces dividend, convertible, governance, capped supply, vesting,
  loyalty and multi-chain token standards with accompanying CLI and VM
  integration.
- Stage 21 streamlines core CLI operations for network, node and access
  management, adding structured output and error propagation for peer
  discovery, staking and address utilities.
- Stage 22 refines AI contract and audit CLI modules, providing consistent error
  handling and JSON-formatted output for integration with external tooling.
- Stage 23 adds gas-aware consensus and DAO governance commands. CLI operations
  such as block mining and DAO creation now emit their expected gas cost for
  better planning and integration with wallets and GUIs.
- Stage 24 expands cross-chain bridges and Plasma management. CLI commands now
  surface gas usage for inter-chain transfers and support JSON output for web
  dashboards.
- Stage 25 adds comprehensive node management. Full, light, mining, mobile,
  optimisation, staking, watchtower and warfare nodes expose JSON emitting CLI
  operations for integration with GUIs and automation.
- Stage 26 enhances operational utilities. Gas table management now allows
  runtime opcode price adjustments and JSON snapshots so dashboards and
  governance tools can consume pricing data directly from the CLI.
- Stage 27 adds developer and testnet automation scripts, covering network bootstrapping, contract deployment, linting and test execution for streamlined workflows.
- Stage 28 introduces release packaging, documentation generation, CI setup and ledger backup scripts for reproducible builds and disaster recovery.
- Stage 29 adds deployable smart contract templates for token faucets, storage markets, DAO governance, NFT minting and AI model marketplaces. These templates are accessible via CLI with gas-priced opcodes.
- Stage 30 introduces utility smart contract modules including escrow payments, cross-chain bridges, multisig wallets and regulatory compliance contracts. Each template ships precompiled as WASM and can be deployed through the CLI and VM.
- Stage 31 debuts a TypeScript GUI wallet that interacts with the CLI. Wallets are generated locally, encrypted with scrypt-derived AES-GCM keys and can query ledger balances through JSON emitting commands.
- Stage 32 adds a CLI-backed Explorer GUI that displays chain height and block data.
- Stage 33 introduces an AI Marketplace GUI that deploys AI-enhanced contracts through the CLI.
- The virtual machine supports smart contracts compiled from WebAssembly, Go, JavaScript, Solidity, Rust, Python and Yul, ensuring opcode compatibility across ecosystems.

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
