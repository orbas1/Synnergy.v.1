# Deployment and Container Architecture

## Overview
Docker images provide a reproducible environment for running Synnergy nodes, wallets and supporting services. Multi‑stage builds compile Go binaries and produce small runtime images that can be orchestrated with Compose or Kubernetes.

## Components
- **Dockerfiles** – each major service has a dedicated Dockerfile with build and runtime stages.
- **docker-compose.yml** – orchestrates local clusters for development and testing.
- **Entrypoint scripts** – scripts in `cmd/scripts/` configure nodes before startup.

## Workflow
1. **Build** – the Dockerfile compiles binaries and bundles configuration defaults.
2. **Package** – artifacts are copied into a minimal base image to reduce attack surface.
3. **Run** – containers expose required ports and volumes for keys and data directories.
4. **Compose orchestration** – `docker-compose.yml` links nodes, wallets and dashboards for quick demos.

## Security Considerations
- Images drop privileges to non‑root users where possible.
- Only necessary binaries and configuration files are included in the runtime stage.
- Environment variables and secrets can be injected at run time rather than baked into the image.

## CLI Integration
Containers typically invoke the `synnergy` binary directly; additional helper scripts can bootstrap networks via Compose.
