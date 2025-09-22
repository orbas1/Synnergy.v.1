# Kubernetes Architecture

## Overview
Kubernetes manifests enable scalable deployments of Synnergy components in cluster environments. They describe desired state for nodes, dashboards and auxiliary services, allowing automated rollout and recovery.

## Components
- **Deployment manifests** – under `deploy/k8s/` for authority nodes, explorers and dashboards.
- **Services** – expose gRPC and HTTP endpoints within the cluster.
- **ConfigMaps and Secrets** – inject node configuration and keys at runtime.
- **Horizontal Pod Autoscalers** – optionally adjust validator or dashboard replicas based on metrics.

## Workflow
1. Apply manifests with `kubectl` or through CI pipelines.
2. Kubernetes schedules pods and mounts ConfigMaps/Secrets containing network settings.
3. Services route traffic to healthy pods; failed instances are restarted automatically.
4. Autoscalers monitor CPU or custom metrics to add or remove replicas.

## Security Considerations
- Secrets are mounted as temporary volumes and not baked into images.
- NetworkPolicies can restrict pod communication to required ports.
- Role Based Access Control governs who may apply or modify manifests.

## CLI Integration
While Kubernetes handles orchestration, containers still invoke the `synnergy` CLI for node operations.

## Enterprise Diagnostics
- Cluster health checks can execute `synnergy integration status --format json` as a readiness probe, ensuring consensus, wallet and authority subsystems are online before routing traffic.
- The integration status output mirrors the dashboard widget so SRE runbooks use consistent telemetry regardless of whether they operate via kubectl, the CLI or the browser.
