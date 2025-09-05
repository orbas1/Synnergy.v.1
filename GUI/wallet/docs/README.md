# Wallet Module Documentation

This directory contains reference material for the Synnergy Wallet GUI. The
wallet provides a lightweight interface for managing Synnergy accounts and
tokens. Developers can use this guide to understand local development,
testing and deployment workflows.

## Features

- Account overview with real‑time balance display
- Modular architecture with pluggable hooks and services
- Jest based unit testing with TypeScript support
- Kubernetes deployment manifest for container orchestration

## Development

```bash
npm install
npm run build   # compile TypeScript sources
npm start       # launch the CLI entry point
npm test        # execute unit tests
```

## Deployment

The `k8s/deployment.yaml` manifest deploys the wallet service behind a single
container. It includes basic health probes and resource limits suitable for a
small‑scale environment. Adjust replicas and resource requests as needed for
production workloads.

## Further Reading

- `../README.md` – project overview and CLI examples
- `../k8s/deployment.yaml` – Kubernetes deployment reference

