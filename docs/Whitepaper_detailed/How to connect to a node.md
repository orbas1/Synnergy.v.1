# How to Connect to a Node

## Overview
Blackridge Group Ltd provides the Synnergy network as a modular, enterprise-grade blockchain. Connecting to a node allows developers and operators to participate in consensus, broadcast transactions, and monitor the system. This guide walks through prerequisites, configuration, secure connection procedures, and diagnostics.

## Prerequisites
- **Source build**: compile the `synnergy` client using Go 1.24 or later.
- **Configuration file**: a YAML file describing network identifiers, peer limits, and listening ports.
- **Network access**: ability to reach peer addresses over the configured P2P port.
- **Security credentials**: user identifiers and biometrics when broadcasting transactions.
- **Hardware resources**: at least 4 CPU cores, 8 GB of RAM, and persistent storage for the blockchain database.
- **Firewall rules**: inbound access to the configured P2P port and outbound egress to known peers.

## Configure Network Parameters
Create or edit a configuration file to define the target network. Key fields include the network identifier, chain ID, maximum peers, and transport settings【F:configs/network.yaml†L1-L10】. Example:

```yaml
network:
  id: synnergy-mainnet
  chain_id: 1215
  max_peers: 50
  p2p_port: 30303
  bootstrap_peers: []
```

Set the `SYN_CONFIG` environment variable to point to this file before launching the client.

## Advanced Configuration
Beyond basic network identifiers, the configuration file controls consensus, virtual machine limits, storage paths, and logging options【F:configs/network.yaml†L12-L40】. Adjust these fields to align with enterprise policies:

- **Consensus** – select algorithm and weighting and tune block times for throughput.
- **RPC exposure** – toggle the JSON‑RPC interface for remote management.
- **Listen address** – bind to specific interfaces for multi‑homed servers.
- **Storage and logging** – choose data directories and log verbosity for audit retention.

## Starting Local Services
Start the networking stack and inspect connectivity using the CLI【F:README.md†L110-L117】:

```bash
export SYN_CONFIG=configs/dev.yaml
./synnergy network start   # launch networking stack
./synnergy network peers   # list currently known peers
```

Stop services with `./synnergy network stop` when finished.

## High Availability and Failover
Production deployments should avoid single points of failure. The `FailoverManager` monitors node heartbeats and promotes the most recent backup if the primary becomes unresponsive【F:high_availability.go†L8-L69】. Deploy at least one standby instance and send regular heartbeats to maintain readiness.

## Dialing a Seed Node
To join an existing mesh, run the base node module and dial a seed peer. The module exposes lifecycle and peer management commands【F:cli/base_node.go†L14-L56】:

```bash
./synnergy basenode start
./synnergy basenode dial <seed-address>
./synnergy basenode peers
```

After the seed handshake completes, the node populates its peer table and begins syncing.

## Peer Discovery and Management
Peer discovery is handled by the `peer` command set built on the `PeerManager` structure【F:cli/peer_management.go†L1-L44】【F:core/peer_management.go†L21-L60】. Each subcommand reports its gas cost and honors the global `--json` flag for machine‑readable responses. Typical workflow:

```bash
./synnergy peer connect <addr> --json      # record a peer by address
./synnergy peer advertise <topic> --json   # announce presence on a topic
./synnergy peer discover <topic> --json    # list peers advertising the topic
```

Connections are pooled and reused through `ConnectionPool`, which creates, dials, and releases link objects while enforcing capacity limits【F:core/connection_pool.go†L22-L75】.

## Firewall and Access Controls
Synnergy includes a built‑in `Firewall` to block malicious wallets, tokens, or IP addresses【F:firewall.go†L5-L44】. Maintain allow and deny lists that mirror corporate security policies, and export rules for compliance reviews.

## Publishing and Subscribing
Use the `network` module to publish messages or subscribe to topics across the peer mesh【F:cli/network.go†L41-L80】:

```bash
./synnergy network broadcast blocks "payload"  # publish data
./synnergy network subscribe blocks             # stream messages until Ctrl-C
```

## Secure Data Channels
For confidential workloads, open encrypted channels through the `ZeroTrustEngine`, which generates per‑channel keys and signs every payload【F:zero_trust_data_channels.go†L24-L66】. Messages can later be verified and decrypted by authorised recipients.

## Security and Authentication
Before a transaction is propagated, the network verifies biometric signatures and attaches them to the payload for accountability. This requirement is enforced by the `Network` service when broadcasting transactions【F:core/network.go†L104-L113】. Operators may also wrap privileged routines with `BiometricSecurityNode` to enforce biometric checks around sensitive actions【F:biometric_security_node.go†L8-L47】.

## Cross-Chain Connectivity
When bridging to external ledgers, manage links with the `ConnectionManager`, which tracks active and historic cross‑chain sessions【F:cross_chain_connection.go†L10-L58】. Connections can be opened, listed, and safely closed as business requirements evolve.

## Monitoring and Diagnostics
- **Peer list**: `./synnergy network peers`
- **Connection pool stats**: expose metrics by instrumenting `ConnectionPool.Stats()` for active vs. capacity numbers【F:core/connection_pool.go†L85-L97】.
- **System health**: `./synnergy system_health snapshot` captures runtime metrics via the `SystemHealthLogger` and can be automated through a `WatchtowerNode` for continuous monitoring【F:system_health_logging.go†L11-L40】【F:watchtower_node.go†L13-L68】.

## Troubleshooting
- **No peers discovered**: verify `bootstrap_peers` and network reachability. Use `peer connect` to add known addresses manually.
- **Pool exhausted**: increase the `max_peers` setting or call `ConnectionPool.Release`/`Close` to free resources.
- **Biometric errors**: confirm the correct user ID and signature payloads are supplied with each broadcast.
- **Failover not triggered**: ensure backups send regular heartbeats and that the timeout window reflects operational SLAs.
- **Firewall blocks legitimate traffic**: review and adjust address and IP allow lists.

## Further Resources
Refer to the CLI help (`./synnergy --help`) for command options and to additional whitepaper sections for governance, security, and cross-chain operations. For enterprise support, contact Blackridge Group Ltd.
