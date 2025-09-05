# Authority Node Index GUI

Enterprise-grade interface for managing and discovering authority nodes within the Synnergy network. Developed and maintained by **Blackridge Group Ltd.**, this module provides a TypeScript foundation with comprehensive tooling for linting, formatting, testing, and containerised deployment. The codebase targets Node.js 18 and ships as a lightweight container image for consistent execution across environments.

## Table of Contents
- [Overview](#overview)
- [Features](#features)
- [Architecture](#architecture)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Development Workflow](#development-workflow)
- [Scripts](#scripts)
- [Testing](#testing)
- [Deployment](#deployment)
  - [Docker](#docker)
  - [Docker Compose](#docker-compose)
  - [Kubernetes](#kubernetes)
- [Operations](#operations)
  - [Logging](#logging)
  - [Monitoring and Health](#monitoring-and-health)
  - [Security](#security)
  - [Troubleshooting](#troubleshooting)
- [Continuous Integration](#continuous-integration)
- [Contributing](#contributing)
- [License](#license)
- [Support](#support)

## Overview
The Authority Node Index GUI delivers a searchable registry of authority nodes. The current implementation focuses on shared tooling and configuration with the companion CLI, providing a solid foundation for future React components that will surface node metadata and operational controls.

## Features
- **Centralised Configuration** – Runtime parameters originate from `config/production.ts` and can be overridden via environment variables or an `.env` file for flexible deployments.
- **Extensible Service Layer** – `src/main.ts` exposes a minimal boot function utilised by both CLI and GUI entry points and includes defensive error handling for reliable startup.
- **UI-Ready Scaffolding** – Placeholder folders for components, hooks, pages, state, and styles provide a clean starting point for a React-driven interface.
- **Comprehensive Testing** – Jest-based unit and end-to-end suites validate configuration defaults, environment overrides, and enforce coverage thresholds.
- **Strict TypeScript Foundation** – `tsconfig.json` enables `strict` mode, ES2020 targets, and source maps for a maintainable codebase.
- **Code Quality Tooling** – ESLint and Prettier are wired through npm scripts and the `Makefile` to guarantee consistent style and formatting.
- **Containerised Distribution** – Multi-stage Docker build produces a compact runtime image suitable for Kubernetes or Docker Compose deployments.
- **Automated CI Pipeline** – `ci/pipeline.yml` mirrors local tasks with `npm ci`, linting, tests, and compilation to keep branches buildable.
- **Live-Reload Friendly** – `docker-compose.yml` mounts the project directory for instant rebuilds during development.

## Architecture
### Configuration
`config/production.ts` reads core settings such as API endpoint, port, and log level. Defaults favour local development and map directly to the environment variables shown below, allowing overrides via shell exports or a dedicated `.env` file.

### Service Layer
`src/main.ts` centralises bootstrapping. The exported `start()` function returns a human‑readable status string built from the configuration, while the CLI entry point logs the message and wraps initialisation in a `try/catch` block to surface failures without crashing the process. The function currently reports the configured port and can be extended to initialise database connections or HTTP servers as the GUI evolves.

### TypeScript Configuration
`tsconfig.json` targets ES2020 and enables `strict` mode, `noImplicitAny`, and source maps for improved debugging. Output is emitted to `dist/`, which is the entry point for runtime execution and Docker images.

### UI Composition
The `src/` tree is organised for a React application:
- `components/` – reusable widgets
- `pages/` – route-level containers
- `hooks/` – shared logic
- `state/` – application stores
- `styles/` – style modules
- `services/` – API adapters and ancillary utilities

Each directory contains a `.gitkeep` placeholder so that the structure is preserved until concrete implementations are committed.

### Testing
Unit tests live beside source code, while additional unit and e2e suites under `tests/` verify both default configuration and environment overrides. Coverage gates in `jest.config.js` require a minimum of 60% line and 50% branch coverage.

## Project Structure
```
config/           # environment-specific configuration
src/              # application source code (components, hooks, services, state, styles)
tests/            # unit and e2e tests
docs/             # supplementary documentation
k8s/              # reference Kubernetes manifests
ci/               # continuous integration pipeline
Dockerfile        # multi-stage build definition
docker-compose.yml# local orchestration with live reloading
Makefile          # common development tasks
package.json      # npm metadata and script definitions
tsconfig.json     # TypeScript compiler settings
jest.config.js    # test runner configuration
```

## Prerequisites
- Node.js 18+
- npm
- Make

## Installation
```bash
make install
```

## Configuration
Copy `.env.example` to `.env` and adjust as needed.

| Variable   | Default | Description |
|------------|---------|-------------|
| `API_URL`  | _empty_ | Base URL for the authority node API; must be supplied for real API calls |
| `PORT`     | `3000`  | HTTP port exposed by the service |
| `LOG_LEVEL`| `info`  | Verbosity for runtime logs |

Defaults are declared in `config/production.ts` and may be overridden at runtime.

Example `.env`:

```
API_URL=https://api.internal
PORT=3000
LOG_LEVEL=debug
```

`LOG_LEVEL` accepts any [pino](https://github.com/pinojs/pino/blob/master/docs/api.md#level-string)‑compatible value such as `error`, `warn`, `info`, or `debug`.

Environment variables override settings in `config/production.ts`; values in a `.env` file are loaded by Docker Compose and local tooling to streamline development. Never commit `.env` to source control—use `.env.example` as the distributable template.

## Development Workflow
| Task                  | Command         |
|-----------------------|-----------------|
| Compile TypeScript    | `make build`    |
| Run tests             | `make test`     |
| Lint sources          | `make lint`     |
| Format with Prettier  | `make format`   |
| Launch compiled service | `make start`  |
| Clean build artifacts | `make clean`    |

## Scripts
| Script | Purpose |
|--------|---------|
| `npm run build` | Compile TypeScript sources to `dist/` using the settings in `tsconfig.json`. |
| `npm start` | Execute the compiled service via Node.js. |
| `npm test` | Run Jest with coverage reporting; same as `make test`. |
| `npm run lint` | Enforce code quality with ESLint. |
| `npm run format` | Apply Prettier formatting across the project. |

## Testing
Run the entire test suite with:

```bash
npm test
```

Coverage reports are emitted to the terminal; the same command is available via `make test`.
`jest.config.js` configures the Node test environment through `ts-jest` and enforces global coverage thresholds of 60% lines and 50% branches. Tests in `tests/unit` validate configuration defaults, while `tests/e2e` confirm that environment overrides are respected.

To watch for file changes during development, run `npm test -- --watch`.

## Deployment
### Docker
```bash
docker build -t authority-node-index .
docker run -p 3000:3000 authority-node-index
```
Images are built in multiple stages to minimise runtime footprint.
Supply environment variables at run time with `-e`, for example `docker run -p 3000:3000 -e API_URL=https://api.internal authority-node-index`.

### Docker Compose
A `docker-compose.yml` file stands up the service together with supporting dependencies and handles environment overrides for local development.

```bash
docker compose up --build
```
The project directory is mounted into the container for live reloads, and environment variables defined in `.env` are injected automatically.

### Kubernetes
`k8s/deployment.yaml` supplies a minimal Deployment manifest with resource requests/limits and HTTP health probes. Adapt it to your cluster policies before production use and wire secrets or ConfigMaps for environment variables.
Scale replicas by adjusting the `spec.replicas` field and apply your own Service or Ingress definitions to expose the application.

## Operations
### Logging
Log verbosity is governed by `LOG_LEVEL`. Messages are emitted to STDOUT/STDERR and can be scraped by your logging stack.

### Monitoring and Health
The Kubernetes manifest expects a `/health` endpoint for liveness and readiness probes. Implement this route within the service when integrating into a production environment.

### Security
Treat `.env` files and runtime secrets with care. Inject sensitive values through Kubernetes Secrets or Docker Compose overrides rather than committing them to the repository. Review logs to ensure no credentials are inadvertently printed.

### Troubleshooting
- **Build failures** – ensure Node.js 18 and npm are installed and run `make clean && make install`.
- **Configuration issues** – verify `.env` values or environment variables; missing `API_URL` prevents API calls.
- **Port conflicts** – change the `PORT` variable if 3000 is already in use.

## Continuous Integration
The GitHub Actions workflow in `ci/pipeline.yml` installs dependencies with `npm ci`, executes the Jest suite, and compiles the project on each push using Node.js 18. Linting and TypeScript compilation run as discrete steps so that failures are surfaced early. This mirrors the tasks in the `Makefile` to keep local and CI runs aligned and ensures that branches remain releasable.

## Contributing
Contributions follow the guidelines in the repository root `CONTRIBUTING.md`. Please run linting and tests before submitting pull requests.

## License
Distributed under the MIT License. Refer to the root `LICENSE` file for full terms.

## Support
For operational assistance or commercial support, contact Blackridge Group Ltd. through your usual support channel.

---
© 2024 Blackridge Group Ltd. All rights reserved.
