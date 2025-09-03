# CLI Node Operations Guide

Stage 25 introduces unified command-line management for core node roles. Each
command supports a global `--json` flag so scripts and GUIs can parse responses
reliably.

## Examples

### Full node
```
synnergy fullnode create --id node1 --mode archive --json
```
Creates a full node and prints its mode.

### Staking node
```
synnergy staking_node stake addr1 100 --json
```
Stakes tokens on behalf of `addr1` and outputs the new balance.

All node commands follow similar patterns with structured output and validated
arguments.
