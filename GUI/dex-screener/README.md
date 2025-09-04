# DEX Screener GUI

This minimal interface lists liquidity pools by spawning the `synnergy` CLI.
It is intended for operational dashboards that need real‑time pool metrics
without embedding chain logic in the browser.

## Usage

```bash
# List all pools
node -e "import('./src/main.ts').then(m => m.listPools()).then(console.log)"

# Inspect a specific pool
node -e "import('./src/main.ts').then(m => m.poolInfo('A-B')).then(console.log)"
```

The script expects the `synnergy` CLI on the `PATH` and a running node with
liquidity pool support.
