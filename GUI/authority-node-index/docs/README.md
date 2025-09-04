# Authority Node Index Documentation

This module exposes a searchable index of authority nodes for both the
command-line interface and the React-based GUI. It demonstrates how
Synnergy services share configuration, logging and test utilities across
delivery channels.

## Architecture

- **Configuration** – values come from `config/production.ts` and can be
  overridden via environment variables.
- **Service** – `src/main.ts` emits a startup message that higher layers use
  to verify connectivity.
- **Tests** – unit tests validate configuration defaults while e2e tests
  confirm environment overrides.

## Development

Run the test suite with:

```bash
npm test
```

Use `k8s/deployment.yaml` as a reference for deploying the service inside a
Kubernetes cluster.
