# Kubernetes Architecture

The Synnergy network can be deployed on Kubernetes to provide automated
scaling, self‑healing and declarative configuration.  The manifests in
`deploy/k8s/` define two core components:

- **Node Deployment** – Runs replicated validator or full nodes. Each pod
  exposes P2P and RPC ports and uses liveness/readiness probes so failed
  instances are replaced automatically.
- **Wallet Server Deployment** – Hosts the HTTP wallet API behind a
  `ClusterIP` service. Probes and resource quotas ensure responsive
  behaviour under load.

Configuration is supplied via ConfigMaps and mounted into the containers.
Operators can extend the manifests with ingress controllers, persistent
storage or PodDisruptionBudgets to satisfy their reliability targets.
