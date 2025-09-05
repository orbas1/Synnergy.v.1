# Graphical User Interfaces (GUIs)

Blackridge Group Ltd develops a comprehensive portfolio of graphical front‑ends that sit atop the Synnergy command‑line interface. Every module ships with opinionated TypeScript and Node.js 18 tooling, strict linting and formatting, and reproducible Docker and Kubernetes assets. These interfaces present complex blockchain interactions through hardened HTTP or gRPC gateways while preserving the deterministic execution model of the underlying CLI.

## Platform Architecture
- **CLI Parity** – GUIs invoke the same binaries as automation scripts, ensuring identical transaction semantics across delivery channels.
- **Typed and Tested** – All source targets strict TypeScript settings with Jest and `ts-jest` harnesses; REST endpoints are exercised with Supertest where applicable.
- **Environment Driven** – `.env.example` files and `config/production.ts` templates expose `API_URL`, `PORT`, `LOG_LEVEL`, `DB_URL` and related knobs so behaviour can be tuned without code changes.
- **Secure Containers** – Multi‑stage Dockerfiles build non‑root runtime images and reference Kubernetes manifests enable health probes and resource controls.
- **Continuous Delivery** – Makefiles and GitHub Actions mirror install, lint, test and build steps to keep every branch releasable.

## Module Catalogue
The following enterprise modules illustrate how diverse Synnergy workflows are rendered through a thin, secure GUI layer.

### Governance & Registry
#### Authority Node Index
Searchable registry of authority nodes with centralised configuration, a shared service layer and comprehensive unit and e2e tests. Defaults for `API_URL`, `PORT` and `LOG_LEVEL` live in `config/production.ts`, with Docker, Compose and Kubernetes manifests providing deployment parity.

#### Validator Governance Portal
Manages validator policies and voting actions. npm scripts handle build, test, lint and format tasks, while Docker recipes expose the service on port 3000 for managed environments.

#### DAO Explorer
TypeScript application for exploring DAO state. Environment variables `API_URL`, `LOG_LEVEL` and `DATABASE_URL` drive runtime settings; unit tests and container images ensure deterministic roll‑outs.

### Security & Compliance
#### Compliance Dashboard
Enterprise dashboard for monitoring regulatory adherence. Requires Node 18+, offers `npm test` and `npm run lint` workflows, and includes Docker and Kubernetes deployment manifests.

#### Security Operations Center
Interface for responding to network security events. Make targets and npm scripts compile TypeScript, run placeholder tests and build multi‑stage Docker images for secure operational deployment.

### Identity & Wallet
#### Identity Management Console
Supports user registration and key management through both CLI and GUI entry points. Ships with build, lint and format scripts plus Docker and Kubernetes manifests for containerised rollout.

#### Wallet
Generates encrypted wallets and queries balances. The scaffold enforces strict TypeScript compilation, Jest tests and Docker Compose orchestration for local development.

#### Wallet Admin Interface
Secure HTTP API for wallet administration and signature verification. Designed for automated CLI interaction and packaged with unit tests, Docker recipes and Kubernetes manifests.

### Marketplaces & Tokenization
#### AI Marketplace
HTTP service and optional CLI for publishing AI‑enhanced contracts. Features an in‑memory state store, comprehensive Jest suites and containerisation assets for Docker, Compose and Kubernetes.

#### Smart Contract Marketplace
Deploys and trades WebAssembly contracts. Exposes a `GET /contracts` endpoint tested with Supertest, and relies on ESLint, Prettier and Jest to enforce quality.

#### Storage Marketplace
Web interface for listing and leasing storage offers. npm scripts cover build, test and start workflows, enabling operators to monetise spare capacity through a consistent CLI wrapper.

#### NFT Marketplace
Front‑end for minting and trading NFTs via opcode‑priced CLI commands. Docker support and scripted tests keep asset workflows deterministic across environments.

#### Token Creation Tool
CLI‑driven interface for issuing new token contracts with predictable gas costs. Provides Prettier, ESLint and Jest tooling plus Docker images for repeatable deployments.

### Network Operations & Analytics
#### Cross‑Chain Management
CLI‑focused module supervising cross‑chain bridges. Environment variables configure API endpoints, log levels and database connections; Jest and ESLint enforce reliability.

#### Cross‑Chain Bridge Monitor
Lightweight monitor for bridge activity across networks. Reads the `API_URL` variable, compiles to a minimal runtime image and includes Jest coverage and formatting checks.

#### Explorer
Reference front‑end for inspecting blockchain state such as chain height and block details. Makefile targets deliver build, test, lint, format and container image creation.

#### Data Distribution Monitor
GUI and CLI for observing data flow across the network. Uses Jest for unit tests and exposes status commands that integrate with higher‑level orchestration.

#### DEX Screener
Monitors liquidity pools through the Synnergy CLI. Includes TypeScript tooling, linting, formatting and comprehensive tests for real‑time market visibility.

#### System Analytics Dashboard
Scaffold for system‑wide analytics. Developers can install dependencies, compile TypeScript, run tests, lint and format using provided npm scripts and environment templates.

#### Node Operations Dashboard
Dashboard for monitoring node health and performance. Builds with TypeScript, runs Jest tests and offers Docker images for deployment.

#### Mining Staking Manager
Scaffolded GUI for managing mining and staking activities with TypeScript build and test scripts.

### Unified CLI Web Panel
The Next.js‑based Synnergy Web Control Panel enumerates every CLI command and exposes them through browser forms. API routes (`/api/commands`, `/api/help`, `/api/run`) invoke `go run ../cmd/synnergy/main.go`, allowing full‑fidelity CLI execution from Vercel or self‑hosted deployments.

## Development Timeline
| Stage | Milestone |
|-------|-----------|
| 13    | Data Distribution Monitor and DEX Screener scaffolds establish foundational dashboards. |
| 23    | Node Operations Dashboard and Security Operations Center gain hardened configuration and service layers. |
| 27‑30 | Storage Marketplace and System Analytics Dashboard mature with full CI pipelines. |
| 31‑36 | Token Creation Tool, Validator Governance Portal, Wallet Admin Interface and Wallet GUI reach production readiness. |
| 37‑40 | Administrative dashboards for authority nodes and cross‑chain management round out the operational toolkit. |

---
© 2024 Blackridge Group Ltd. All rights reserved.

