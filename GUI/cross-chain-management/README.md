# Cross-Chain Management GUI

This dashboard drives the `cross-chain-management` CLI module to list
and establish cross-chain bridges. It executes the `synnergy` binary and
presents the returned JSON data.

## Usage

```bash
# Show registered bridges
node src/main.ts bridges

# Connect to a new chain
node src/main.ts connect chainA https://endpoint
```

Ensure the `synnergy` CLI is configured with credentials and gas policy
suitable for cross-chain operations.
