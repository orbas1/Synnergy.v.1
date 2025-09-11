# Compliance and Regulatory Architecture

## Overview
Compliance components ensure network activity adheres to legal and policy requirements. They provide audit trails, sanction screening and tools for regulators to review transactions without halting decentralized operation.

## Key Modules
- `compliance.go` – core checks for sanctioned addresses and suspicious transfers.
- `compliance_management.go` – aggregates reports and exposes review APIs.
- `regulatory_management.go` – coordinates rulesets from multiple jurisdictions.
- `regulatory_node.go` – special node type that can request logs and attest to audits.

## Workflow
1. **Policy ingestion** – `regulatory_management` loads jurisdictional rules and feeds them to compliance services.
2. **Transaction evaluation** – `compliance` hooks into validation to screen participants before blocks finalize.
3. **Reporting** – flagged events are persisted by `compliance_management` for later inspection.
4. **Regulatory access** – authorized `regulatory_node` instances can query logs and sign off on reviews.

## Security Considerations
- Access to compliance reports is gated through `access_control` and logged for traceability.
- Regulatory nodes operate read-only paths and cannot alter ledger state.
- Rulesets are versioned so historical decisions can be reconstructed during audits.

## CLI Integration
- `synnergy compliance` – run ad-hoc checks and view current rulesets.
- `synnergy compliance-mgmt` – manage reports and export audit logs.
