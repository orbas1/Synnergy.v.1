# Docker Deployment

This directory provides container images and a compose configuration for running a minimal Synnergy network.

## Images
- `synnergy` – node binary exposing the blockchain and consensus engine on port `8080`.
- `walletserver` – backend service used by GUIs and CLI wallets on port `8090`.

Both images are built with multi-stage Docker builds producing minimal `alpine` images.

## Usage
```
docker compose -f docker/docker-compose.yml up --build
```
The compose file mounts configuration files from `configs/` and starts the node before launching the wallet server.

## Customisation
Set `SYN_CONFIG` to point at an alternative configuration file or extend the compose file with additional services such as the web GUI.

## Fault tolerance
For production deployments use orchestration platforms like Kubernetes and enable health checks and restart policies to ensure nodes automatically recover from failures.
