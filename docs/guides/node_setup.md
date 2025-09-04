# Node Setup Guide

This guide walks through installing and running a Synnergy node.

## Prerequisites
- Go 1.21+
- Git
- Access to the Synnergy repository

## Installation
```bash
git clone https://github.com/example/synnergy.git
cd synnergy
make build
```

## Running a Node
```bash
./synnergy --config configs/example.yaml
```

## Troubleshooting
- Ensure ports are open and not blocked by firewalls.
- Check logs under `logs/` for detailed error messages.
