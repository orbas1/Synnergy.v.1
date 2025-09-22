# Architecture Documents

This directory contains architecture overviews for major module groups. Each file outlines the components and CLI commands associated with a functional area.

1. [Consensus](consensus_architecture.md)
2. [Loanpool](loanpool_architecture.md)
3. [Storage and Data](storage_architecture.md)
4. [Cross-Chain](cross_chain_architecture.md)
5. [Identity and Access](identity_access_architecture.md)
6. [AI](ai_architecture.md)
7. [Node Roles](node_roles_architecture.md)
8. [Security](security_architecture.md)
9. [Tokens and Transactions](tokens_transactions_architecture.md)
10. [Governance and DAO](governance_architecture.md)
11. [Compliance and Regulatory](compliance_architecture.md)
12. [Monitoring and Logging](monitoring_logging_architecture.md)
13. [Specialized Features](specialized_architecture.md)
14. GUI Modules
   - [Wallet](wallet_architecture.md)
   - [Explorer](explorer_architecture.md)
   - [AI Marketplace](ai_marketplace_architecture.md)

15. [Deployment and Containerization](docker_architecture.md)

There are **15** distinct module groups.

Stage 80 introduces the Synthron Treasury orchestration layer, which spans multiple architecture domains by wiring ledger, VM, wallet, consensus and authority services into a single control plane. The orchestration service aligns gas schedules, registers treasury opcodes (including operator governance), exposes telemetry through the CLI and function web, and ensures enterprise dashboards share the same mint/burn, circulation, operator and compliance insights with form-driven controls for governed actions【F:treasury/synthron_treasury.go†L41-L523】【F:cli/coin.go†L23-L130】【F:web/pages/index.js†L1-L210】.
