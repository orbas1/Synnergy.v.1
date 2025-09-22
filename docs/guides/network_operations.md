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

## Automated Devnet Bootstrap

Stage 100 ships `scripts/devnet_start.sh`, a wrapper that primes the CLI for
local swarm testing:

1. Creates a heavy `simplevm` profile and starts it.
2. Boots the network stack, joins the requested number of nodes via `synnergy
   swarm join` and generates encrypted wallets for each participant.
3. Adjusts adaptive consensus weights and publishes a readiness event on the
   `devnet/announce` topic.

The script honours `--dry-run`, `--timeout` and `--log-file` flags and stores
wallets in `var/devnet/wallets`.

## Container Orchestration

- `scripts/docker_build.sh` adds retries, environment overrides and optional
  registry pushes to `docker build`.
- `scripts/docker_compose_up.sh` launches the stack with profile selection,
  scaling hints and post-start health summaries sourced from
  `docker compose ps --format json`.

Together they provide an auditable path from source to containers without
relying on ad-hoc shell invocations.

## Smoke Testing

Run `scripts/e2e_network_tests.sh` to exercise the simple VM, consensus weight
adjustments and swarm consensus loop. The helper validates that:

- Swarm membership matches the requested node count.
- Consensus weights remain balanced (sum close to 1.0).
- Transition thresholds stay positive after adjustments.
- Network broadcasts succeed through the CLI.

## Upgrades

- Pull new releases and run `go build` to update the binary.
- Restart the node with minimal downtime using systemd or Docker restarts.
- Follow release notes for any required migrations.

## Key Management

- Back up validator keys stored under `keys/`.
- Rotate keys with `synnergy validator rotate` and update peers accordingly.

For advanced procedures, consult the whitepaper and module guides.

