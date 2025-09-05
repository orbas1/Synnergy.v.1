# Blackridge Group Ltd – AI Marketplace Module

The AI Marketplace is a modular interface for discovering, evaluating, and acquiring AI services within the Synnergy Network. Developed and maintained by **Blackridge Group Ltd**, this package delivers a production‑ready scaffold for building marketplace capabilities or integrating external tooling.

## Overview

The module exposes a lightweight HTTP service and optional command‑line entry point. It ships with opinionated defaults for linting, formatting, and testing so teams can move from prototype to enterprise deployment with minimal friction.

## Key Features

- **TypeScript runtime** targeting Node.js 18 for modern language features and reliability.
- **HTTP service** defined in `src/main.ts` that responds with a liveness message and can be extended with REST or GraphQL routes.
- **In‑memory state store** in `src/state/store.ts` for quick prototyping of marketplace logic.
- **Testing harness** using Jest and `ts-jest` with both unit tests (`tests/unit`) and end‑to‑end tests (`tests/e2e`).
- **Code quality tooling** via ESLint and Prettier to enforce consistent style across contributors.
- **Containerization and orchestration** assets including Dockerfile, Docker Compose, and Kubernetes deployment manifests under `k8s/`.
- **GitHub Actions pipeline** (`ci/pipeline.yml`) that installs dependencies, runs tests, and builds distributable artifacts on every push or pull request.

## Project Structure

```text
GUI/ai-marketplace/
├── config/           # Environment‑specific settings
├── docs/             # Architectural notes and operational guides
├── src/              # Application source code
│   ├── main.ts       # HTTP entry point
│   └── state/store.ts# In‑memory key‑value store
├── tests/            # Unit and e2e test suites
├── ci/               # Continuous integration pipeline
├── k8s/              # Kubernetes manifests for example deployment
├── Dockerfile        # Production container definition
├── docker-compose.yml# Local development stack
├── Makefile          # Common developer tasks
└── package.json      # Scripts, dependencies, and metadata
```

## Getting Started

### Prerequisites
- Node.js 18+
- npm

### Installation and Local Run

```bash
npm install
npm run build
npm start
```

`npm start` invokes the compiled entry point and boots an HTTP service that defaults to port `3000`. Replace or extend this service with real marketplace functionality as needed.

### Makefile Shortcuts

```bash
make install   # npm install
make build     # npm run build
make test      # npm test
```

## Available npm Scripts

| Command | Description |
| ------- | ----------- |
| `npm run build` | Compile TypeScript sources into `dist/` |
| `npm start` | Launch the compiled HTTP service |
| `npm test` | Execute Jest unit and e2e tests |
| `npm run lint` | Lint the project using ESLint |
| `npm run format` | Format source using Prettier |

## Configuration

Runtime configuration lives in `config/production.ts`. Two environment variables are supported:

| Variable | Purpose | Default |
| -------- | ------- | ------- |
| `API_URL` | Target Synnergy API endpoint | `https://api.synnergy.example` |
| `PORT` | HTTP port for the service | `3000` |

These values can be overridden via shell environment variables or container orchestration tools.

## Testing & Quality Assurance

Jest configuration (`jest.config.js`) uses `ts-jest` to compile TypeScript on the fly and collects coverage reports. The test suites demonstrate basic store operations and HTTP liveness checks to provide a template for further scenarios.

Run tests locally:

```bash
npm test
npm run lint
```

## Continuous Integration

The GitHub Actions workflow in `ci/pipeline.yml` provisions Node.js 18, installs dependencies, runs the test suite, and builds the TypeScript output on every commit. This ensures consistent quality across development branches and pull requests.

## Deployment

- **Docker:** `Dockerfile` builds a minimal Node.js 18 image with compiled assets. Use `docker build` and `docker run` for custom deployments.
- **Docker Compose:** `docker-compose.yml` exposes the service on port `3000` for local experimentation.
- **Kubernetes:** `k8s/deployment.yaml` offers an example Deployment manifest for cluster orchestration.

## Further Documentation

Additional architectural context and operational guides are located in [`docs/README.md`](docs/README.md).

## License

This module is released under the terms of the repository's primary license. © Blackridge Group Ltd.

