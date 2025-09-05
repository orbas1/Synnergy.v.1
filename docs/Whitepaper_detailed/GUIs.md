# Graphical User Interfaces (GUIs)

Blackridge Group Ltd. maintains a portfolio of enterprise‑grade graphical front‑ends that extend the Synnergy command‑line interface. Each package is built on a strict Node.js 18 and TypeScript foundation, enforced by ESLint, Prettier, Jest and `ts-jest`. Multi‑stage Dockerfiles, Kubernetes manifests and Makefiles provide reproducible builds, while GitHub Actions mirror local workflows so every branch can ship with confidence.

## Platform Architecture and Enterprise Standards
- **CLI Fidelity** – GUIs call the same `synnergy` binaries as automation scripts, guaranteeing identical transaction semantics across browser, desktop and server environments.
- **TypeScript Baseline** – `tsconfig.json` files enable `strict` mode and modern ECMAScript targets, with Jest and Supertest exercising both functional paths and HTTP endpoints.
- **Environment‑Driven Configuration** – `.env.example` templates and `config/production.ts` files expose `API_URL`, `PORT`, `LOG_LEVEL`, `DB_URL` and related knobs, allowing behaviour to be tuned without code changes.
- **Security Posture** – Non‑root containers, CORS/Helmet middleware, and explicit dependency locking protect each surface. Kubernetes resources specify liveness/readiness probes and resource quotas.
- **Quality and Delivery** – Makefiles and npm scripts standardise install, lint, test, build and release steps. GitHub Actions pipelines (`ci/pipeline.yml`) run these tasks on every push.
- **Observability Ready** – Modules emit structured logs via Winston and expose health endpoints, enabling integration with enterprise telemetry stacks.
- **UX Foundations** – The `docs/ux` directory codifies accessibility aids, mobile responsiveness, theming, localisation, error handling and command history guidelines that all GUIs follow.

## Module Catalogue
The following modules illustrate how diverse Synnergy workflows are rendered through a secure GUI layer.

### Governance & Registry
#### Authority Node Index (`GUI/authority-node-index`)
- Searchable registry of authority nodes with centralised configuration and a shared service layer.
- Defaults for `API_URL`, `PORT` and `LOG_LEVEL` reside in `config/production.ts`; unit and e2e tests enforce behaviour.
- Ships with Makefile tasks, Docker/Compose/Kubernetes manifests and a CI pipeline.

#### Validator Governance Portal (`GUI/validator-governance-portal`)
- Manages validator policies and voting actions through TypeScript services.
- npm scripts handle build, test, lint and format; Docker recipes expose the service on port 3000.

#### DAO Explorer (`GUI/dao-explorer`)
- Explores DAO state with environment variables `API_URL`, `LOG_LEVEL` and `DATABASE_URL` driving runtime settings.
- Deterministic roll‑outs are ensured via Jest tests and container images.

### Security & Compliance
#### Compliance Dashboard (`GUI/compliance-dashboard`)
- Monitors regulatory adherence with Node 18+, ESLint, Prettier and Jest in the toolchain.
- Docker and Kubernetes manifests provide repeatable deployment targets.

#### Security Operations Center (`GUI/security-operations-center`)
- Responds to network security events through a hardened TypeScript interface.
- Make targets compile sources, run placeholder tests and build multi‑stage Docker images for secure operations.

### Identity & Wallet
#### Identity Management Console (`GUI/identity-management-console`)
- Enables user registration and key management via both CLI and GUI entry points.
- Build, lint and format scripts accompany Docker and Kubernetes assets for containerised rollout.

#### Wallet (`GUI/wallet`)
- Generates encrypted wallets and queries balances with strict TypeScript compilation.
- Jest tests and Docker Compose orchestration support local development.

#### Wallet Admin Interface (`GUI/wallet-admin-interface`)
- Exposes a secure HTTP API for wallet administration and signature verification.
- Designed for automated CLI interaction and packaged with unit tests, Docker recipes and Kubernetes manifests.

### Marketplaces & Tokenization
#### AI Marketplace (`GUI/ai-marketplace`)
- Publishes AI‑enhanced contracts through an HTTP service and optional CLI.
- Uses an in‑memory state store, comprehensive Jest suites and containerisation assets for Docker, Compose and Kubernetes.

#### Smart Contract Marketplace (`GUI/smart-contract-marketplace`)
- Deploys and trades WebAssembly contracts, exposing a `GET /contracts` endpoint.
- Relies on ESLint, Prettier and Jest to enforce quality across builds.

#### Storage Marketplace (`GUI/storage-marketplace`)
- Lists and leases storage offers with build, test and start workflows scripted in npm.

#### NFT Marketplace (`GUI/nft_marketplace`)
- Mints and trades NFTs via opcode‑priced CLI commands.
- Docker support and scripted tests keep asset workflows deterministic.

#### Token Creation Tool (`GUI/token-creation-tool`)
- CLI‑driven interface for issuing new token contracts with predictable gas costs.
- Provides Prettier, ESLint and Jest tooling plus Docker images for repeatable deployments.

### Network Operations & Analytics
#### Cross‑Chain Management (`GUI/cross-chain-management`)
- Supervises cross‑chain bridges with environment variables for `API_URL`, `LOG_LEVEL` and `DB_URL`.
- Jest, ESLint and formatted scripts enforce reliability and code quality.

#### Cross‑Chain Bridge Monitor (`GUI/cross-chain-bridge-monitor`)
- Lightweight monitor for bridge activity; reads `API_URL` and compiles to a minimal runtime image.
- Includes Jest coverage and formatting checks.

#### Explorer (`GUI/explorer`)
- Reference front‑end for inspecting blockchain state such as chain height and block details.
- Makefile targets deliver build, test, lint, format and container image creation.

#### Data Distribution Monitor (`GUI/data-distribution-monitor`)
- GUI and CLI for observing data flow across the network.
- Uses Jest for unit tests and exposes status commands for higher‑level orchestration.

#### DEX Screener (`GUI/dex-screener`)
- Monitors liquidity pools through the Synnergy CLI with TypeScript tooling, linting and comprehensive tests.

#### System Analytics Dashboard (`GUI/system-analytics-dashboard`)
- Provides system‑wide analytics with npm scripts for build, test, lint and format.
- Environment templates simplify deployment across stages.

#### Node Operations Dashboard (`GUI/node-operations-dashboard`)
- Monitors node health and performance with TypeScript builds and Jest tests.
- Docker images enable predictable deployment.

#### Mining Staking Manager (`GUI/mining-staking-manager`)
- Manages mining and staking activities via a TypeScript scaffold and placeholder tests.

### Unified CLI Web Panel (`web`)
- Next.js application that enumerates every Synnergy CLI command and exposes them through browser forms.
- API routes—`/api/commands`, `/api/help` and `/api/run`—invoke `go run ../cmd/synnergy/main.go`, enabling full‑fidelity CLI execution from Vercel or self‑hosted deployments.
- Includes an example `/authority` page that renders authority node data from the CLI.

### User Experience Standards
Documentation under `docs/ux` outlines cross‑module guidelines for:
- mobile responsiveness
- accessibility aids
- error handling and validation
- loading feedback and status indicators
- theming options and localisation support
- command history, onboarding help and authentication roles

### Operational Practices
- **Logging & Metrics** – Winston‑based structured logs stream to STDOUT/STDERR and can be scraped by external observability platforms.
- **Security Controls** – Containers run as non‑root, `.env` files keep secrets external, and Kubernetes manifests wire liveness and readiness probes for zero‑trust environments.
- **Scalability** – Docker and Kubernetes assets define resource limits and allow horizontal scaling via replica counts and load‑balanced services.
- **Delivery Pipeline** – GitHub Actions pipelines and Makefiles codify installation, linting, testing, building and image publishing for consistent releases.

## Development Timeline
| Stage | Milestone |
|-------|-----------|
| 13 | Data Distribution Monitor and DEX Screener scaffolds establish foundational dashboards. |
| 23 | Node Operations Dashboard and Security Operations Center gain hardened configuration and service layers. |
| 27‑30 | Storage Marketplace and System Analytics Dashboard mature with full CI pipelines. |
| 31‑36 | Token Creation Tool, Validator Governance Portal, Wallet Admin Interface and Wallet GUI reach production readiness. |
| 37‑40 | Administrative dashboards for authority nodes and cross‑chain management round out the operational toolkit. |

---
© 2024 Blackridge Group Ltd. All rights reserved.

