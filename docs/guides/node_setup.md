# Node Setup Guide

This guide walks through installing and running a Synnergy node.

## Prerequisites

- Go 1.21+
- Git
- Optional: Docker for containerised runs

## Installation

```bash
git clone https://github.com/example/synnergy.git
cd synnergy
make build
```

The build produces the `synnergy` binary under `./bin`.

## Configuration

Copy a sample config and adjust to your environment:

```bash
cp cmd/config/default.yaml mynode.yaml
```

Edit network ID, ports and logging paths as required.

## Running a Node

```bash
./synnergy --config mynode.yaml
```

To run under systemd, create a service file pointing to the binary and enable it with `systemctl enable synnergy.service`.

## Security

- Restrict RPC ports with a firewall.
- Backup keys under `keys/` and set appropriate file permissions.

## Troubleshooting

- Ensure ports are open and not blocked by firewalls.
- Check logs under `logs/` for detailed error messages.

