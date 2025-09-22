# Docker Deployment

Stage 77 refreshes the container toolchain so local and CI environments mirror
enterprise production settings. Images are compiled with `-trimpath`, tested
during the build stage, and run as an unprivileged `synnergy` user with
healthchecks that surface the same endpoints exposed to Kubernetes probes.

## Images

- `synnergy/runtime:stage77` – consensus, virtual-machine, and CLI entrypoint
  exposing P2P (30303), RPC (8080), and Prometheus metrics (9102).
- `synnergy/wallet:stage77` – wallet API with mTLS, HSM, and telemetry hooks on
  port `8081` and metrics on `9103`.

## Compose Orchestration

The updated `docker-compose.yml` provisions a full stack including
OpenTelemetry and Prometheus services plus an optional Next.js web dashboard.

```bash
docker compose -f docker/docker-compose.yml up --build
```

Volumes persist ledger state (`synnergy-ledger`) and configuration files are
mounted read-only from `configs/` so CLI, wallet, and authority workflows
operate against deterministic manifests.

### Profiles

- `default` – launches node, wallet, telemetry, and Prometheus services.
- `ui` – adds the web dashboard for exercising the JavaScript interface.

## Customisation

- Override `SYN_CONFIG` or telemetry endpoints via environment variables.
- Supply custom Prometheus or OpenTelemetry configuration files to integrate
  with existing observability backends.

## Production Alignment

For real deployments pair these images with the Stage 77 Kubernetes manifests or
Terraform modules. Health checks and metrics endpoints are consistent across
orchestrators, enabling seamless migration from local testing to HA clusters.
