# Data Distribution Monitor Docs

This section details architecture and integration notes for the data distribution monitor.

- **CLI**: invoke `node dist/main.js status` to check service health.
- **Configuration**: set `API_URL` and `PORT` environment variables for production deployments.
- **Kubernetes**: see `k8s/deployment.yaml` for an example deployment manifest with resource limits and health probes.
