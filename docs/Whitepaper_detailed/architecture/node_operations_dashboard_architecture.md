# Node Operations Dashboard Architecture

## Overview
The Node Operations Dashboard aggregates metrics from authority and validator nodes to give operators a single view of network health. It combines a TypeScript frontend with REST endpoints exposed by nodes and watchtower services.

## Key Components
- `GUI/node-operations-dashboard` – frontend that renders charts and status tables.
- Node status endpoints provided by authority nodes.
- `watchtower_node.go` – feeds health information to the dashboard.

## Workflow
1. **Data collection** – nodes expose JSON endpoints with uptime, block height and peer counts.
2. **Aggregation** – the dashboard polls these endpoints and normalizes results.
3. **Visualization** – operators view charts, alerts and logs through the GUI.

## Security Considerations
- Dashboard requests are read-only and can be gated behind authentication.
- TLS is recommended for all endpoints to prevent tampering.
- Rate limits stop the dashboard from overwhelming nodes during polling.

## CLI Integration
- Operators can script status checks using `synnergy watchtower` and feed output into the dashboard.
