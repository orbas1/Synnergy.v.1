# Storage Marketplace GUI

A reference interface for listing and leasing decentralised storage via the Synnergy CLI.

## Usage

```bash
# Create a listing
node src/main.ts list <hash> <price> <owner>

# List available storage
node src/main.ts listings

# Open a deal
node src/main.ts deal <listingID> <buyer>
```

The implementation executes the local `synnergy` binary, inheriting its network configuration and gas pricing.
