# Security Operations Center

This module provides a minimal dashboard for monitoring security events in a Synnergy network.  It exposes a TypeScript entry point that reads configuration from environment variables and communicates with the core network through the function web.

## Usage

```bash
npm install
npm run build
node dist/main.js
```

### Docker

```bash
docker compose up --build
```

### Kubernetes

Deploy using the manifests in `k8s/`:

```bash
kubectl apply -f k8s/deployment.yaml
```

### Testing

Run unit tests with Jest:

```bash
npm test
```
