# Ai Marketplace

The AI Marketplace module provides a web-based interface and command-line entrypoint for interacting with Synnergy network services.  It is built with TypeScript and ships with opinionated linting, formatting, and testing defaults so the project can be used as an enterprise-grade scaffold.

## Getting Started

```bash
npm install
npm run build
npm start
```

`npm start` launches the command‑line interface which in turn boots a minimal HTTP service used by the tests and example tooling.  The service is intended to be replaced with real marketplace logic.

## Scripts

| Command | Description |
| ------- | ----------- |
| `npm run build` | Compile TypeScript sources into `dist/` |
| `npm test` | Run Jest unit tests |
| `npm run lint` | Lint the project using ESLint |
| `npm run format` | Format source using Prettier |

## Container & Deployment

This module includes a `Dockerfile`, `docker-compose.yml`, and Kubernetes manifests under `k8s/` for local experimentation.  See `docs/` for architecture notes.
