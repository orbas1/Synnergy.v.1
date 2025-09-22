# Blocks and Sub-Blocks – Stage 77 Architecture

Synnergy organises transactions into sub-blocks that are validated by
specialised authority or validator nodes before being sealed into blocks by the
hybrid consensus engine. Stage 77 enhances this pipeline with runtime
integration, automated gas validation and deployment blueprints that keep CLI,
web and infrastructure layers aligned.

## Sub-Block Assembly

- Nodes collect transactions into sub-blocks using deterministic ordering and
  validator selection to maintain fairness and prevent censorship【F:core/node.go†L1-L80】【F:core/consensus.go†L144-L206】.
- Runtime integration ensures the node, consensus engine and wallet share the
  same gas schedule and regulatory approvals before sub-block creation, allowing
  CLI automation and web dashboards to observe identical state transitions【F:internal/runtime/integration.go†L15-L192】【F:cmd/synnergy/main.go†L9-L125】.
- Watchtower nodes and pod-level health checks monitor sub-block propagation and
  consensus status, restarting containers or shifting workloads when issues are
  detected without interrupting CLI workflows【F:watchtower_node.go†L13-L86】【F:deploy/k8s/node.yaml†L1-L143】.

## Block Finalisation

- Once validated, sub-blocks are aggregated into blocks with fee distribution
  and stake adjustments calculated via the ledger and validator manager helpers,
  guaranteeing reward transparency and economic stability【F:core/node.go†L81-L160】【F:core/consensus_validator_management.go†L5-L84】.
- The virtual machine executes smart-contract opcodes during block finalisation,
  now consulting the shared gas table per opcode to avoid budget overruns and to
  keep documentation, CLI and contract execution synchronised【F:virtual_machine.go†L24-L142】【F:gas_table.go†L108-L126】.
- Integration tests exercise this path end-to-end by signing transactions,
  feeding them through the runtime integration and verifying ledger height and
  health telemetry, providing regression coverage for enterprise deployments【F:integration_test.go†L5-L66】.

## Deployment and Observability

- Kubernetes manifests define PodDisruptionBudgets, ConfigMaps and telemetry
  sidecars to maintain block production continuity during upgrades or failures
  while exporting metrics for dashboards and alerting systems【F:deploy/k8s/node.yaml†L1-L143】.
- Terraform provisions auto-scaling groups, load balancers and encrypted
  configuration stores so block producers can be rolled out across availability
  zones with consistent runtime configuration and access control【F:deploy/terraform/main.tf†L1-L231】.
- Docker Compose bundles the node, wallet and web UI, enabling local operators to
  inspect block and sub-block state using the same RPC interfaces deployed in
  production clusters【F:docker/docker-compose.yml†L1-L64】.

Stage 77 therefore delivers a block production pipeline that is deterministic,
observable and fault tolerant across CLI, automation and browser-based control
planes, ensuring enterprise-grade reliability for Synnergy’s hybrid consensus.
