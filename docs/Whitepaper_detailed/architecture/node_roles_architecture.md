# Node Roles Architecture

## Overview
Synnergy defines several specialized node types to fulfill distinct responsibilities, from content distribution to regulatory oversight. Each role implements tailored logic while sharing common networking and consensus layers.

## Example Roles
- `content_node.go` – serves application data and handles content type registration.
- `indexing_node.go` – builds searchable indexes for the explorer and analytics tools.
- `regulatory_node.go` – interfaces with compliance modules to provide audit capabilities.
- `watchtower_node.go` – observes peer health and reports anomalies.
- `mobile_mining_node.go` – lightweight miner suited for resource‑constrained devices.

## Workflow
Nodes advertise their role during handshake, enabling peers to route requests appropriately. Shared services such as consensus and storage remain consistent, while role-specific modules extend functionality.

## Security Considerations
- Role capabilities are limited by `access_control` to prevent privilege escalation.
- Specialized nodes may run in isolated environments or networks depending on sensitivity.
- Heartbeat checks from watchtower nodes help detect compromised or offline roles.

## CLI Integration
- `synnergy indexing-node` – manage indexing services.
- `synnergy watchtower` – administer watchtower nodes.
- `synnergy regulatory-node` – access regulatory functions.
