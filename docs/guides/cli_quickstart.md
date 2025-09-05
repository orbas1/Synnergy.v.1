# CLI Quickstart

The Synnergy CLI allows interaction with the blockchain for common tasks.

## Installation
Ensure the project is built:
```bash
make build
```

## Common Commands
- `synnergy wallet create` – generate a new wallet
- `synnergy tx send` – send a transaction
- `synnergy node status` – display node synchronization status
- `synnergy charity_pool --json registration <addr>` – view charity registration info as JSON
- `synnergy charity_mgmt donate <from> <amount>` – donate tokens to the charity pool

## Help
Run `synnergy --help` or `synnergy <command> --help` for more details.
