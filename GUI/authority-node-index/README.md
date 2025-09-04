# Authority Node Index GUI

This lightweight dashboard interfaces with the Synnergy CLI to manage
and inspect authority nodes. It spawns the `synnergy` binary with the
`authority-node-index` module and displays the JSON output.

## Usage

```bash
# List registered authority nodes
node src/main.ts list

# Register a new authority node
node src/main.ts register 0xabc123 0xpubkey
```

The script assumes the `synnergy` CLI is on the PATH and configured with
appropriate network access and gas pricing.
