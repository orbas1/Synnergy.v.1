# Synnergy: Enterprise Blockchain Framework

Synnergy is a modular, high-performance blockchain written in Go and built for enterprise production use.  It exposes pluggable node roles, cross-chain interoperability and AI-assisted contract tooling so organisations can prototype, pilot and run distributed networks with the same codebase.

## Table of Contents
- [Key Features](#key-features)
- [Architecture & Code Map](#architecture--code-map)
- [Bootstrap Sequence](#bootstrap-sequence)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Build from Source](#build-from-source)
  - [Configuration](#configuration)
  - [Run a Local Node](#run-a-local-node)
  - [Multi-node Devnet](#multi-node-devnet)
- [CLI Modules](#cli-modules)
- [Production Deployment](#production-deployment)
  - [Docker Compose](#docker-compose)
  - [Kubernetes (Helm)](#kubernetes-helm)
  - [Terraform & Ansible](#terraform--ansible)
- [High Availability & Scaling](#high-availability--scaling)
- [Monitoring & Observability](#monitoring--observability)
- [Security & Compliance](#security--compliance)
- [Testing](#testing)
- [Documentation](#documentation)
- [Contributing](#contributing)
- [License](#license)

## Key Features
- **Pluggable node roles** – mining, staking, authority, regulatory, watchtower, warfare and more via constructors such as `core.NewMiningNode` and `core.NewRegulatoryNode`.
- **Cross-chain interoperability** – bridges, connection managers and transaction relays (`core.NewBridgeRegistry`, `core.NewChainConnectionManager`, `core.NewCrossChainTxManager`).
- **AI modules** – contract management, inference analysis, anomaly detection and secure storage (`core.NewAIEnhancedContract`, `core.NewAIDriftMonitor`).
- **Gas accounting** – deterministic costs loaded via `synnergy.LoadGasTable()` and registered with `synnergy.RegisterGasCost()`.
- **Role-based security** – biometric authentication and security node CLI (`core.NewBiometricService`, `synnergy bioauth`, `synnergy bsn`), zero‑trust data channels and PKI tooling.
- **Extensible CLI** – built with [Cobra](https://github.com/spf13/cobra) and backed by `cli.Execute()`.
- **Validated block utilities** – Stage 40 adds sub-block creation and block assembly commands with strict argument checking.
- **Central bank controls** – `synnergy centralbank` manages monetary policy and CBDC issuance with structured JSON output.
- **Monetary policy utilities** – `synnergy coin` provides validated reward and supply queries with optional JSON.
- **Charity pool management** – `synnergy charity_pool` and `charity_mgmt` support registration, donations and balance queries with JSON responses.
- **Data distribution monitor** – CLI and GUI for network telemetry.
- **Node operations dashboard** – TypeScript CLI and GUI for real-time node health monitoring.
- **Security operations center** – monitors and aggregates security events with a CLI-accessible GUI.
- **Smart contract marketplace** – publish and trade vetted contracts through CLI or GUI with deterministic gas costs.
- **Storage marketplace & analytics dashboards** – TypeScript modules exposing CLI and GUI entrypoints for decentralised storage and system metrics.
- **Infrastructure-as-code** – Dockerfiles, Helm charts, Terraform and Ansible playbooks for reproducible environments.
- **Strong encryption and signatures** – all transactions and messages are secured using well‑vetted cryptography with digital signatures.
- **Permissioned privacy** – fine‑grained access controls enable private channels and selective data disclosure.
- **Customisable governance** – token‑weighted voting and DAO modules allow on‑chain policy definition.
- **Horizontal scalability** – sharding-ready architecture and multi‑node orchestration for high throughput networks.
- **Modular web interfaces** – NFT marketplace and node operations dashboard shipped with tested TypeScript frontends.
- **Regulatory compliance hooks** – KYC/AML modules and audit logging for standards alignment.

## Architecture & Code Map
```
cmd/          Command-line entry points (`cmd/synnergy`, `cmd/watchtower`, …)
cli/          Modular CLI commands (network, wallet, contracts, mining, …)
core/         Blockchain runtime (consensus, networking, VM, node roles)
configs/      YAML configuration templates
internal/     Shared utilities such as `config.Load` and token helpers
GUI/          Web and desktop front-ends
deploy/       Docker, Helm, Terraform and Ansible manifests
scripts/      Operational scripts for setup, testing and automation
docs/         MkDocs sources and reference guides
pkg/          Reusable libraries and experimental modules
```

## Bootstrap Sequence
`cmd/synnergy/main.go` orchestrates start-up:
1. Load environment variables with `gotenv.Load()`.
2. Resolve configuration path from `SYN_CONFIG` or `config.DefaultConfigPath`.
3. Parse YAML via `config.Load` and configure logging.
4. Initialise tracer provider (`otel.SetTracerProvider`).
5. Warm caches by calling `synnergy.LoadGasTable()` and `synnergy.RegisterGasCost()` for core operations like `MineBlock`, `OpenConnection` and `MintNFT`.
6. Pre-load modules used by the CLI:
   - `core.NewNetwork` for pub‑sub networking
   - `core.NewContractRegistry` backed by `core.NewSimpleVM`
   - DAO managers (`core.NewDAOManager`, `core.NewProposalManager`, …)
   - Wallet, watchtower and warfare nodes (`core.NewWallet`, `core.NewWatchtowerNode`, `core.NewWarfareNode`)
   - Token constructors in `internal/tokens` (`tokens.NewSYN223Token`, etc.)
7. Finally, invoke `cli.Execute()` to dispatch Cobra commands.

## Getting Started

### Prerequisites
- Go **1.24+**
- Make (optional but recommended)
- Docker & Node.js for container builds and GUI projects

### Build from Source
```bash
go build ./cmd/synnergy    # or: make build
```
The `synnergy` binary is written to the repository root.

### Configuration
Synnergy loads a YAML configuration file using `config.Load`. The path is resolved from the `SYN_CONFIG` environment variable and defaults to `configs/dev.yaml`. Example:
```yaml
environment: development
log_level: debug
server:
  host: "127.0.0.1"
  port: 8080
database:
  url: "https://dev-db.example.com"
```
For production builds, compile with `go build -tags prod ./cmd/synnergy` so `config.DefaultConfigPath` points to `configs/prod.yaml`.

### Run a Local Node
```bash
export SYN_CONFIG=configs/dev.yaml
./synnergy network start           # start networking stack
./synnergy network peers           # list peers
./synnergy network stop            # stop services
./synnergy wallet new --out wallet.json --password pass
./synnergy system_health snapshot  # inspect runtime metrics
```

### Multi-node Devnet
Launch a disposable development network with multiple nodes:
```bash
scripts/devnet_start.sh 3   # spawns 3 nodes
```

## CLI Modules
Run `./synnergy --help` for the full command tree. Common modules include:

| Command | Description |
| ------- | ----------- |
| `network start|stop|peers|broadcast|subscribe` | Manage the P2P layer backed by `core.NewNetwork` |
| `wallet new` | Generate encrypted wallets via `core.NewWallet` |
| `mining start|status|stop|attempt` | Operate a mining node (`core.NewMiningNode`) |
| `staking_node start|status|stop` | Control the staking service |
| `contracts compile|deploy|invoke|list|info` | WASM smart contract lifecycle through `core.NewContractRegistry` |
| `system_health snapshot|log` | Emit metrics and structured logs |
| `data monitor status` | Report network data distribution metrics |
| `audit log|list` | Record and query audit events via `core.NewAuditManager` |
| `audit_node start|log|list` | Operate a bootstrap audit node for network-wide logs |
| `authority register|vote|list` | Manage the authority node registry (`core.NewAuthorityNodeRegistry`) |
| `authority_apply submit|vote|finalize|list` | Handle authority node applications (`core.NewAuthorityApplicationManager`) |
| `bankinst register|list|is` | Manage institutional bank participants (`core.NewBankInstitutionalNode`) |
| `banknodes types` | Display supported bank node categories |
| `basenode start|stop|running|peers|dial` | Control a base network node (`core.NewBaseNode`) |
| `basetoken init|mint|balance` | Interact with a basic token (`tokens.NewBaseToken`) |
| `dex liquidity <pair>` | Query on-chain liquidity pool reserves |

Additional modules cover DAO governance, cross-chain bridges, regulatory nodes, watchtowers and more.

## Production Deployment

### Docker Compose
```bash
docker compose -f docker/docker-compose.yml up --build
```
The compose file builds the `cmd/synnergy` Docker image and launches an example node.

### Kubernetes (Helm)
```bash
helm install synnergy deploy/helm/synnergy
```
Override settings with `--set` or by editing `deploy/helm/synnergy/values.yaml` to specify image tags, replicas and resource limits.

### Terraform & Ansible
```bash
cd deploy/terraform
terraform init
terraform apply -var 'ami_id=ami-123456'

ansible-playbook -i <inventory> deploy/ansible/playbook.yml
```
These templates provision cloud infrastructure and configure nodes with hardened defaults.

## High Availability & Scaling
Scripts under `scripts/` provide building blocks for resilient clusters:
```bash
scripts/high_availability_setup.sh    # configure active-active replicas
scripts/ha_failover_test.sh           # simulate node failure
scripts/active_active_sync.sh         # cross-region state sync
```
For geo-distributed deployments use Kubernetes with multiple replicas and load balancers to achieve automated failover.

## Monitoring & Observability
Synnergy emits structured logs via Logrus and initializes an OpenTelemetry tracer. Metrics can be exported or inspected directly:
```bash
./synnergy system_health snapshot
scripts/metrics_export.sh     # stream metrics to Prometheus
scripts/metrics_alert_dispatch.sh  # send alert to webhook
```
Integrate results with Prometheus, Grafana or the ELK stack.

## Security & Compliance
```bash
make security                 # run staticcheck, gosec and govulncheck
scripts/pki_setup.sh          # generate certificate authority and node certs
scripts/wallet_hardware_integration.sh  # configure hardware wallets
```
See [SECURITY.md](SECURITY.md) for policies, and use `scripts/aml_kyc_process.sh` or `scripts/compliance_audit.sh` for regulatory workflows.

## Testing
```bash
go test ./...
make test                 # convenience wrapper
scripts/run_tests.sh          # execute full suite including integration tests
npm test --prefix GUI/storage-marketplace            # run storage marketplace tests
npm test --prefix GUI/system-analytics-dashboard     # run dashboard tests
```

## Documentation
```bash
make docs        # build static site into site/
make docs-serve  # serve documentation locally
```
Reference documentation is generated from the `docs/` directory using MkDocs.

## Contributing
Development follows the staged workflow outlined in [AGENTS.md](AGENTS.md). Format code with `go fmt`, run `go vet`, `go build` and `go test` on touched packages before opening a pull request. Review [CHANGELOG.md](CHANGELOG.md) and [PRODUCTION_STAGES.md](PRODUCTION_STAGES.md) for roadmap and release process.

## License
Synnergy is provided for research and educational purposes. Third-party dependencies retain their original licenses.

