# Internal Security Platform

The internal security package provides production-grade primitives that power the
Synnergy Network validator core, CLI utilities, the function web interface, and
the regulatory compliance workflows. Each component is designed to be auditable,
self-healing, and horizontally scalable.

## Components

### Key Management
- Deterministic key purposes allow the VM, consensus engine, wallets, and CLI to
  request the exact key material they require without exposing unrelated secrets.
- High-entropy generation with audit logging, deterministic testing hooks, and
  support for importing keys from external HSMs.
- Ed25519 signing helpers surface both the signature and key version so remote
  verifiers can enforce rotation policies.

### Encryption Service
- XChaCha20-Poly1305 authenticated encryption with Ed25519 signatures to deliver
  zero-trust envelopes consumable by Go services, the CLI, and the web UI.
- Seamless key rotation with fingerprinting that feeds governance and regulatory
  dashboards.
- Support for associated data so transactions, governance votes, and authority
  node attestations can be bound into the cryptographic domain.

### DDoS Mitigation
- Adaptive scoring and burst controls tuned for validator RPC workloads.
- Deterministic snapshots support CLI inspection, consensus telemetry, and web
  dashboards.
- Manual quarantine hooks allow authority nodes or automated scripts to respond
  to targeted attacks in real time.

### Patch Governance
- Digitally signed patch metadata ensures that only authorised releases are
  recorded.
- Metadata exports feed compliance reports and serve as the backbone for the
  authority node governance portal.
- Backwards compatible `Applied()` helper keeps legacy CLI integrations working
  while exposing richer metadata for the GUI.

## Integration Points
- **CLI**: new commands read the audit logs, patch metadata, and DDoS snapshots
  to assist operators during incident response.
- **Virtual Machine**: sealed envelopes are used for state checkpointing and for
  hot-swappable WASM modules. Rotated keys are pushed to the VM through the key
  manager subscribers.
- **Consensus**: the consensus orchestrator polls the DDoS mitigator to throttle
  malicious peers and consumes key rotation events to refresh Noise/TLS
  transports.
- **Wallet + Node Infrastructure**: wallets consume the signing helpers for
  transaction approval flows, while nodes use the patch manager to coordinate
  release rollouts.

## Operational Best Practices
1. Configure the key manager with an HSM backed entropy source in production.
2. Mirror the audit log into the governance ledger for immutable record keeping.
3. Rotate Noise and TLS keys using the key rotator at least once per epoch.
4. Monitor the DDoS snapshot metrics and escalate sustained high scores to the
   security operations centre.
5. Require digitally signed patch metadata before applying release artifacts to
   authority nodes.

