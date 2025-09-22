# Kubernetes Manifests

Stage 77 elevates the Synnergy deployment blueprints to enterprise readiness. The
manifests now align with the production guidance captured in the CLI, gas, and
consensus documentation, providing high-availability defaults, observability
integration, and hardened security postures.

## Files

- `node.yaml` – Declares a `Deployment`, `PodDisruptionBudget`, and `Service`
  for consensus/virtual-machine nodes. The workload ships with OpenTelemetry and
  Prometheus sidecars, runs as an unprivileged user, and mounts persistent
  storage for ledger data. Topology spread constraints and pod anti-affinity
  guarantee distribution across availability zones while readiness and liveness
  probes align with the CLI health commands.
- `wallet.yaml` – Provides a redundant wallet API tier with HSM configuration
  secrets, mTLS certificates, and automatic failover alerts. Health probes are
  tuned for low-latency UX paths and the pod-level security context enforces a
  read-only root filesystem.

## Deployment

1. Provision supporting ConfigMaps and Secrets:
   ```bash
   kubectl apply -f configs/node-config.yaml
   kubectl apply -f configs/wallet-config.yaml
   kubectl apply -f secrets/wallet-hsm.yaml
   ```
2. Apply the workloads:
   ```bash
   kubectl apply -f deploy/k8s/node.yaml
   kubectl apply -f deploy/k8s/wallet.yaml
   ```
3. Monitor rollout status with the CLI to ensure gas-table and consensus
   metadata are in sync:
   ```bash
   synnergy status cluster --namespace synnergy
   ```

## Operational Notes

- Metrics endpoints are exposed on port `9102` for nodes and `9103` for wallet
  pods and integrate with the Stage 77 gas catalogue by exporting opcode
  categories as Prometheus labels.
- Both deployments propagate tracing headers to the OpenTelemetry collector,
  ensuring virtual-machine executions are traceable across node, wallet, and
  authority workflows.
- Pod disruption budgets, node affinity, and persistent volumes ensure the
  cluster tolerates node maintenance without interrupting CLI access or consensus
  quorum.
