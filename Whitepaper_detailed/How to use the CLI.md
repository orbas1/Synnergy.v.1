# How to use the CLI

The `synnergy` binary exposes network and contract functionality.
Stage 39 extends the toolset with liquidity pool commands used by the DEX screener.

```bash
# Create a new liquidity pool with default fee
synnergy liquidity_pools create TOKENA TOKENB

# Add liquidity and then inspect all pools
synnergy liquidity_pools add TOKENA-TOKENB provider 100 100
synnergy liquidity_views list

# Query a specific pool
synnergy liquidity_views info TOKENA-TOKENB
```

Most commands accept the `--json` flag to produce machine readable output for GUIs.
