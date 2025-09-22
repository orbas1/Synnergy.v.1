# Network Operations Guide

This guide outlines day‑to‑day tasks for running Synnergy nodes and maintaining the network.

## Joining the Network

1. Run `scripts/node_setup.sh --start-mining` to provision a wallet, stake funds,
   grant the validator role and start the miner with Function Web telemetry
   enabled. The script wraps the CLI workflow documented in the Node Setup Guide
   and emits JSON suitable for automation pipelines.
2. Register your node using `synnergy node register` if you need deterministic
   peer IDs for static peering.
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

- Use `synnergy health status` for local checks and
  `scripts/network_diagnostics.sh --iterations 10` to record JSON snapshots of
  peer availability. The diagnostics script is resilient to transient failures
  and retries operations with exponential back‑off.
- Aggregated dashboards are available via the node operations dashboard GUI.
- Logs are written to `logs/`; rotate them periodically. The
  `scripts/network_harness.sh` utility exercises broadcast, consensus and mining
  pipelines to validate alerting configurations after upgrades.

## Upgrades

- Pull new releases and run `scripts/package_release.sh --version <tag>` to
  execute the vetted build pipeline. Use `scripts/release_sign_verify.sh` to sign
  or verify release archives before deployment.
- When migrating between releases, execute `scripts/network_migration.sh` to
  capture pre/post peer snapshots and validate that consensus continues to
  produce blocks through a rolling restart.
- For full failover drills, run `scripts/network_partition_test.sh` to simulate
  a partition and confirm that the network recovers with acceptable gas costs.

## Key Management

- Back up validator keys stored under `keys/`.
- Rotate keys with `synnergy validator rotate` and update peers accordingly.

For advanced procedures, consult the whitepaper and module guides.

