# Mining Staking Manager Documentation

This module provides a minimal interface for managing mining and staking
operations on the Synnergy network.  It is designed as a starting point for
building richer dashboards and can run either locally or inside Docker/Kubernetes
environments.

## Getting Started

```bash
npm install
npm run build
npm start
```

The service reads configuration from environment variables as defined in
`config/production.ts`.

## Testing

Run unit and end‑to‑end tests with:

```bash
npm test
```

## Deployment

Refer to `docker-compose.yml` for local containerised execution and
`k8s/deployment.yaml` for a minimal Kubernetes deployment.
