# Build Error Tracking

### Stage 1
core/address.go:84:2: syntax error: non-declaration statement outside function body

### Stage 2
core/cross_chain_bridge.go:162:1: syntax error: unexpected ==, expected }

### Stage 3
core/cross_consensus_scaling_networks.go:14:1: syntax error: unexpected keyword type, expected field name or embedded type

### Stage 4
Building package `synnergy` fails due to errors in `synnergy/core`.

### Stage 5
Building package `synnergy/cli` fails due to errors in `synnergy/core`.

### Stage 6
Building package `synnergy/cmd/synnergy` fails due to errors in `synnergy/core`.

### Stage 7
Package `synnergy/Tokens` builds successfully.

### Stage 8
Packages under `synnergy/node_ext` build successfully.

### Stage 9
Package `synnergy/nodes` builds successfully.

### Stage 10
Packages under `synnergy/nodesextra` build successfully.
