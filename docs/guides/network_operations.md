# Network Operations Guide

This guide outlines day‑to‑day tasks for running Synnergy nodes and maintaining the network.

## Joining the Network

1. Follow the [Node Setup Guide](node_setup.md).
2. Register your node using `synnergy node register`.
3. Open required ports and configure firewalls.

## Staking

- Delegate tokens with `synnergy stake add <amount>`.
- Check status with `synnergy stake status`.
- Withdraw with `synnergy stake remove <amount>` after the unbonding period.

## Governance

- List proposals: `synnergy gov proposals`
- Vote on a proposal: `synnergy gov vote <id> <yes|no>`
- Submit a new proposal: `synnergy gov submit <file>`

## Monitoring

- Use `synnergy health status` for local checks.
- Aggregated dashboards are available via the node operations dashboard GUI.
- Logs are written to `logs/`; rotate them periodically.

## Upgrades

- Pull new releases and run `go build` to update the binary.
- Restart the node with minimal downtime using systemd or Docker restarts.
- Follow release notes for any required migrations.

## Key Management

- Back up validator keys stored under `keys/`.
- Rotate keys with `synnergy validator rotate` and update peers accordingly.

For advanced procedures, consult the whitepaper and module guides.

