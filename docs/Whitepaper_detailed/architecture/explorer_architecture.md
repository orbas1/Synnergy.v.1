# Explorer Architecture

## Overview
The Explorer provides read‑only insight into ledger state and network activity. A small TypeScript frontend shells out to the `synnergy` CLI to fetch blocks, transactions and account balances.

## Key Components
- **Frontend** – SPA built with TypeScript that renders tables and charts.
- **CLI Adapter** – Node layer that executes commands like `synnergy ledger head` and `synnergy tx show`.
- **Indexing node** – optional backend using `indexing_node.go` for faster queries.

## Workflow
1. The user requests a block or account view in the browser.
2. The adapter runs the appropriate CLI command and captures structured JSON output.
3. Results are parsed and displayed with minimal processing.

## Security Considerations
- The explorer only performs read operations; signing keys never leave the node.
- Rate limits on the adapter protect nodes from heavy scraping.
- Content security policy headers restrict external scripts in the frontend.

## CLI Integration
- `synnergy ledger` – retrieve block headers and history.
- `synnergy tx` – inspect transaction details.

## Enterprise Diagnostics
- The Explorer’s landing page now consumes `/api/integration`, which proxies `synnergy integration status`, to display validator, consensus and wallet readiness before any block queries are executed.
- Operations teams can compare the Explorer widget with CLI output to confirm UI, CLI and backend data paths remain in sync after upgrades.
