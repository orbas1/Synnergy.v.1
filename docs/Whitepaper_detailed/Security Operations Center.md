# Security Operations Center

The security operations center (SOC) provides real‑time visibility into threats against the Synnergy network. It consumes logs and alerts from nodes, aggregates them and presents actionable insights through a web interface. The SOC communicates with the rest of the platform via the function web, allowing administrators to query status and trigger responses through the CLI.

## Features
- Environment‑driven configuration compatible with Docker and Kubernetes.
- Fault‑tolerant startup with explicit error handling.
- Extensible service layers for custom alerting and dashboards.

## Integration
The SOC uses standard REST calls to interact with Synnergy nodes and authority services.  Authentication is handled through API tokens issued by the wallet service. Future work includes deeper integration with consensus metrics and token‑based authorization.
