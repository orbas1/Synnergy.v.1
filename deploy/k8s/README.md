# Kubernetes Manifests

This directory contains example Kubernetes manifests for deploying Synnergy components.

- `node.yaml` defines a `Deployment` and `Service` for a replicated Synnergy node.
- `wallet.yaml` provides a `Deployment` and `Service` for the wallet server.

Both manifests include resource requests/limits and HTTP health probes for
fault‑tolerant operation. Apply them with `kubectl`:

```bash
kubectl apply -f deploy/k8s/node.yaml
kubectl apply -f deploy/k8s/wallet.yaml
```

Configuration is expected to be provided via ConfigMaps named
`syn-node-config` and `syn-wallet-config` respectively. Adjust images and
ports as needed for your environment.
