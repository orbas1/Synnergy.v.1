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
- [Stage 78 Enterprise Diagnostics](#stage-78-enterprise-diagnostics)
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
- **Cross-chain interoperability** – bridges, protocol registries, connection managers and transaction relays (`core.NewBridgeRegistry`, `core.NewProtocolRegistry`, `core.NewChainConnectionManager`, `core.NewCrossChainTxManager`). Stage 42 finalises these modules with JSON emitting CLIs for deposits, claims and contract mappings.
- **AI modules** – contract management, inference analysis, anomaly detection and secure storage (`core.NewAIEnhancedContract`, `core.NewAIDriftMonitor`).
- **Gas accounting** – deterministic costs loaded via `synnergy.LoadGasTable()` now include rich metadata surfaced through `synnergy.GasCatalogue()` so CLI and web dashboards can describe categories, default limits and intent. Runtime overrides continue to work through the `SYN_GAS_OVERRIDES` environment variable, while Stage 91 adds `synnergy.RegisterGasMetadata()` to validate names, attach descriptions and keep the VM and consensus engines aligned on pricing.
- **Stage 81 module catalogue** – `synnergy modules list` streams consensus, VM, wallet, node and authority readiness with gas prices for both CLI automation and the JavaScript control panel.
- **Kademlia DHT tools** – CLI support for storing, retrieving and computing XOR distance between keys with gas-aware execution.
- **Controlled proof-of-work** – mining nodes expose a `MineUntil` helper and `synnergy mining mine-until` command so hashing stops when a context is cancelled, a prefix target is reached, or a timeout fires.
- **Opcode-aware contracts** – sample Solidity bridges, liquidity, multisig, oracle and token contracts invoke SNVM opcodes with deterministic gas costs.
- **On-chain governance token** – the SYN300 module supports delegated voting, proposal lifecycle management and CLI tooling with documented gas and opcode mappings.
- **Regulatory node logging** – `synnergy regnode approve` surfaces rejection reasons and address logs for auditability.
- **Regulatory audits** – `synnergy regnode audit` reports whether an address has been flagged and returns recorded reasons.
- **Wallet-signed approvals** – regulatory nodes verify transactions against registered wallet public keys, and the CLI can load and sign with a wallet file before submitting for approval.
- **Web regulatory console** – the `web/pages/regnode.js` interface exposes approval, flagging and log retrieval through a browser UI.
- **Strict flagging** – regulatory flags require explicit non-empty reasons for improved audit integrity.
- **Regulator-optional consensus** – if no regulatory node is configured, sub-block validation bypasses compliance checks so transactions continue to flow.
- **Role-based security** – biometric authentication and security node CLI (`core.NewBiometricService`, `synnergy bioauth`, `synnergy bsn`), zero‑trust data channels and PKI tooling. Authority node voting, base-node peering and biometric flows all require Ed25519 signatures for verifiable governance.
- **Institutional banking** – bank nodes require Ed25519-signed requests to register or remove participating institutions.
- **Bank node index** – track banking nodes via a thread-safe index exposed through `bank_index` CLI commands.
- **Extensible CLI** – built with [Cobra](https://github.com/spf13/cobra) and backed by `cli.Execute()`.
- **Resource-managed CLI** – connection pool commands can release individual peers and `contractopcodes` reports gas costs for contract operations.
- **Content node pricing** – gas table and opcode registry expose costs for registering nodes, uploading content, retrieving items and listing hosts so storage workflows remain predictable across the CLI and web UI.
- **Content registry & secrets tooling** – `synnergy content_node` manages hosted content while the standalone `secrets-manager` binary validates stored keys.
- **Enterprise orchestrator** – Stage 78 introduces `core.NewEnterpriseOrchestrator` and the `synnergy orchestrator` CLI to unify VM readiness, consensus relayers, wallet bootstrap, authority registry state and gas documentation with telemetry for CLI and web dashboards.
- **DAO governance** – `synnergy dao` manages decentralised autonomous organisations with optional JSON output, ECDSA signature verification, admin-controlled member role updates via `dao-members update`, and elected authority node term renewals.
- **Resilient node primitives** – forensic nodes prune over-capacity logs, full node modes are mutex-protected, gateway endpoints require a running node, and failover managers can remove stale peers.
- **Thread-safe mempools and plasma bridge safeguards** – node mempools are mutex-protected for concurrent submissions and Plasma bridge operations surface explicit paused errors.
- **Peer management utilities** – `synnergy peer count` reports known peers with gas-aware output for network monitoring.
- **Validated block utilities** – Stage 40 adds sub-block creation and block assembly commands with strict argument checking.
- **Consensus tooling** – Stage 41 adds validated commands for adaptive weighting, difficulty control and service management.
- **Adaptive consensus management** – Stage 63 averages recent demand and stake
  metrics to stabilise weight shifts across PoW, PoS and PoH.
- **Central bank controls** – `synnergy centralbank` manages monetary policy and CBDC issuance with structured JSON output.
- **Monetary policy utilities** – `synnergy coin` provides validated reward and supply queries with optional JSON.
- **Compliance management** – `synnergy compliance` and `compliance_management` emit JSON results for KYC validation, fraud scoring and address policy status.
- **Charity pool management** – `synnergy charity_pool` and `charity_mgmt` support registration, donations and balance queries with JSON responses.
- **Data distribution monitor** – CLI and GUI for network telemetry.
- **Node operations dashboard** – TypeScript CLI and GUI for real-time node health monitoring.
- **Security operations center** – monitors and aggregates security events with a CLI-accessible GUI.
- **Smart contract marketplace** – publish and trade vetted contracts through CLI or GUI with deterministic gas costs.
- **Storage marketplace & analytics dashboards** – TypeScript modules exposing CLI and GUI entrypoints for decentralised storage and system metrics.
- **Infrastructure-as-code** – Dockerfiles, Helm charts, Terraform and Ansible playbooks for reproducible environments.
- **Strong encryption and signatures** – all transactions and messages are secured using well‑vetted cryptography with digital signatures.
- **Stage 75 enterprise runtime** – deterministic VM gas enforcement with execution traces, sandbox telemetry for fault-tolerant VM pools, wallets with deterministic seeds and shared-secret support, and event-driven cross-chain registries/bridges that stream lifecycle updates to CLI and web dashboards.
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
5. Warm caches by calling `synnergy.LoadGasTable()`, synchronising the Stage 78 enterprise schedule with `synnergy.EnsureGasSchedule()` and registering contextual metadata with `synnergy.RegisterGasMetadata()` for core operations like `MineBlock`, `OpenConnection`, `MintNFT` and the enterprise orchestrator opcodes.
6. Pre-load modules used by the CLI:
   - `core.NewNetwork` for pub‑sub networking
   - `core.NewContractRegistry` backed by `core.NewSimpleVM`
   - DAO managers (`core.NewDAOManager`, `core.NewProposalManager`, …)
   - Wallet, watchtower and warfare nodes (`core.NewWallet`, `core.NewWatchtowerNode`, `core.NewWarfareNode`)
   - Token constructors in `internal/tokens` (`tokens.NewSYN223Token`, etc.)
7. Finally, invoke `cli.Execute()` to dispatch Cobra commands.

Enterprise deployments rely on the defensive modules initialised above:

- **Warfare node telemetry** – `core.NewWarfareNode` now issues per-commander key pairs, enforces signed command envelopes with
  replay protection, and streams logistics/tactical events to both CLI (`synnergy warfare events`) and the web console.
- **Watchtower observability** – `core.NewWatchtowerNode` emits start/stop/fork alerts alongside periodic metric snapshots via
  `Watchtower.SubscribeEvents`, enabling dashboards and the `synnergy monitoring` suite to visualise consensus health in real
  time.
- **Zero-trust data channels** – `core.NewZeroTrustEngine` orchestrates encrypted channels with participant governance,
  key rotation and retention controls. CLI helpers (`synnergy zero-trust authorize`, `rotate`, `events`) expose the
  same lifecycle flows used by automation scripts and the browser UI.

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
Stage 100 ships a hardened bootstrapper that provisions wallets, the
virtual-machine sandbox and the adaptive consensus engine before joining a
swarm of nodes. The script logs every CLI invocation and persists wallets to
`var/devnet/wallets` by default.

```bash
# Spin up a 4-node swarm, initialise a heavy VM profile and bias consensus weights.
scripts/devnet_start.sh --count 4 --vm-mode heavy --demand 0.65 --stake 0.55

# Dry-run mode prints the exact CLI calls without executing them.
scripts/devnet_start.sh --count 2 --dry-run
```
Logs are written to `scripts/logs/devnet_start.log`. Override the target
directory or wallet password via `--wallet-dir` and `--wallet-password`.

## CLI Modules
Run `./synnergy --help` for the full command tree. Common modules include:

| Command | Description |
| ------- | ----------- |
| `network start|stop|peers|broadcast|subscribe` | Manage the P2P layer backed by `core.NewNetwork` |
| `wallet new` | Generate encrypted wallets via `core.NewWallet` |
| `modules list [--json]` | Inspect Stage 81 module catalogue, opcode coverage and gas pricing for CLI/VM/web orchestration |
| `mining start|status|stop|attempt` | Operate a mining node (`core.NewMiningNode`) |
| `staking_node start|status|stop` | Control the staking service |
| `contracts compile|deploy|invoke|list|info` | WASM smart contract lifecycle through `core.NewContractRegistry` |
| `system_health snapshot|log` | Emit metrics and structured logs |

| `data feed create|apply|snapshot|get|delete` | Persist structured data feeds with JSON manifests |
| `data resource import|info|list|usage|put|prune` | Manage binary resources with retention metadata and manifest-driven sync |
=======
| `orchestrator status|sync [--json]` | Aggregate VM, consensus, wallet and authority diagnostics via `core.NewEnterpriseOrchestrator` |
| `data monitor status` | Report network data distribution metrics |
| `audit log|list` | Record and query audit events via `core.NewAuditManager` |
| `audit_node start|log|list` | Operate a bootstrap audit node for network-wide logs |
| `authority register|vote|list` | Manage the authority node registry (`core.NewAuthorityNodeRegistry`) |
| `authority_apply submit|vote|finalize|list` | Handle authority node applications (`core.NewAuthorityApplicationManager`) |
| `bankinst register|remove|list|is --pub --sig` | Manage institutional bank participants with signed requests (`core.NewBankInstitutionalNode`) |
| `banknodes types` | Display supported bank node categories |
| `basenode start|stop|running|peers|dial` | Control a base network node (`core.NewBaseNode`) |
| `basetoken init|mint|balance` | Interact with a basic token (`tokens.NewBaseToken`) |
| `dex liquidity <pair>` | Query on-chain liquidity pool reserves |
| `syn500 create|grant|use` | Manage service-tier utility tokens |
| `syn3800 create|release|get|list` | Manage programmatic grants via `core.GrantRegistry` |
| `syn3900 register|claim|get` | Track government benefits (`core.BenefitRegistry`) |
| `syn4200_token donate|progress` | Record charity donations and view campaign totals |
| `syn4700 create|sign|status|info|dispute` | Administer legal-document tokens |
| `dao-members add|remove|role|list` | Manage DAO membership with JSON output and ECDSA verification |

Additional modules cover DAO governance, cross-chain bridges, regulatory nodes, watchtowers and more.

Helper scripts under `cmd/scripts` honour a `SYN_CLI` environment variable to locate the compiled binary and enable `set -euo pipefail` for safer automation.


### Data Governance Automation

Stage 100 introduces a hardened data-governance toolchain that couples new CLI primitives with operational scripts:

- `scripts/data_operations.sh` applies JSON feed manifests through `synnergy data feed apply`, exports snapshots and prunes stale keys so downstream analytics stay consistent.
- `scripts/data_resource_manage.sh` imports binary resources from declarative manifests, triggers CLI audits for the referenced keys and can prune any catalogue entry not listed.
- `scripts/data_retention_policy_check.sh` validates resource freshness by comparing `synnergy data resource info` timestamps against manifest-defined retention windows and can fail CI/CD pipelines when violations occur.

All commands persist state beneath `~/.synnergy/data` (override with `SYN_DATA_DIR`), enabling repeatable governance workflows across CLI invocations, automation scripts and the web tooling described in `docs/guides/cli_quickstart.md` and the whitepaper appendices.

### Operational Automation

New Stage 100 scripts extend the operational toolkit so container builds,
runtime orchestration and smoke tests share the same CLI flows as production
nodes:

- `scripts/devnet_start.sh` spins up a wallet-backed swarm, adjusts consensus
  weights and announces readiness on the network pub/sub bus.
- `scripts/docker_build.sh` wraps `docker build` with retries, environment
  overrides and optional pushes to private registries.
- `scripts/docker_compose_up.sh` launches the reference stack with structured
  logging, profile selection and post-start health checks.
- `scripts/e2e_network_tests.sh` exercises the VM, consensus engine and swarm
  membership end-to-end, validating that thresholds and weight adjustments stay
  within safe bounds.

Each helper honours the shared `--dry-run`, `--timeout` and `--log-file`
switches provided by `scripts/lib/common.sh` so workflows are deterministic
across CI, staging and developer laptops.
=======
## Stage 78 Enterprise Diagnostics
Stage 78 upgrades the runtime with a hardened enterprise orchestrator that validates end-to-end readiness across the virtual machine, consensus mesh, wallets, node registries and gas documentation. The orchestrator powers the `synnergy orchestrator` CLI and exports JSON suitable for the function web dashboards so operators can embed live diagnostics into existing tooling.

```bash
./synnergy orchestrator status          # human-readable snapshot
./synnergy orchestrator status --json   # machine-readable diagnostics
./synnergy orchestrator sync            # refresh gas schedule & authority counts
```

Diagnostics confirm VM mode and concurrency, consensus network registration, wallet provenance, authority node totals, and whether Stage 78 opcodes are documented with enterprise-grade gas costs. Results surface through the updated Next.js API (`web/pages/api/orchestrator.js`) and dashboard widgets on the control panel home page, ensuring parity between CLI automation and browser operations. Stress, situational and real-world tests under `core/enterprise_orchestrator_test.go` and `cli/orchestrator_test.go` assert fault tolerance, security controls and regulatory alignment, keeping performance predictable even under high-throughput workloads.


## Production Deployment

### Docker Compose
```bash
scripts/docker_build.sh --tag synnergy/devnet:latest
scripts/docker_compose_up.sh --project synnergy-dev --profile gui
```
`docker_build.sh` keeps retry logs under `scripts/logs/` and can tag/push to
private registries. `docker_compose_up.sh` streams Compose output and records
service health summaries so CI pipelines can assert readiness without parsing
raw logs.

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
scripts/e2e_network_tests.sh  # stage 100 smoke test for VM/consensus/network
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

