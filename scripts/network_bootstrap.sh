#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$SCRIPT_DIR/.."

# Build the synnergy CLI
GO111MODULE=on go build -trimpath -o "$ROOT_DIR/synnergy" "$ROOT_DIR/cmd/synnergy"

# Use network configuration and genesis file
export SYN_CONFIG="$ROOT_DIR/configs/network.yaml"

# Show configured genesis wallets for operators
"$ROOT_DIR/synnergy" genesis show

# Stake initial validator, authority and host nodes
"$ROOT_DIR/synnergy" node stake d22b7fffcb6a72c94713f5c2e2f2142565b8402299842eb4c1039cea3c293ff0 1000000
"$ROOT_DIR/synnergy" node stake 5baf968c0972eab52046bbca74763d61c923d6e81e549f0f7f0db66f8cfddad9 500000
"$ROOT_DIR/synnergy" node stake 20d2698d356868e08411ce2ea8ac824d07011aebe788a6458d33e244172b8c34 500000

# Launch network and consensus services
"$ROOT_DIR/synnergy" network start &
NET_PID=$!
"$ROOT_DIR/synnergy" consensus-service start 3000 &
CONS_PID=$!
trap 'kill "$NET_PID" "$CONS_PID"' INT TERM
wait
