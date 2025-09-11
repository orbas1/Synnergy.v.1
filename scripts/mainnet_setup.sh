#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$SCRIPT_DIR/.."

usage() {
  cat <<'USAGE'
Usage: mainnet_setup.sh
Build and bootstrap a Synnergy mainnet node with initial resources.
USAGE
}

if [[ ${1:-} == "--help" ]]; then
  usage
  exit 0
fi

# Build the synnergy CLI
GO111MODULE=on go build -trimpath -o "$ROOT_DIR/synnergy" "$ROOT_DIR/cmd/synnergy"

# Use main network configuration
export SYN_CONFIG="$ROOT_DIR/configs/network.yaml"

# Display configured genesis wallets
"$ROOT_DIR/synnergy" genesis show

# To redirect the creator fee share to node hosts once the network is live,
# run:
# "$ROOT_DIR/synnergy" creator disable-distribution --wallet /path/to/creator.wallet --password <password>

# Initialise genesis block once
if ! "$ROOT_DIR/synnergy" genesis init 2>/dev/null; then
  echo "Genesis already initialised"
fi

# Generate wallets for initial token distribution
mkdir -p "$ROOT_DIR/wallets"
for name in distribution reserve treasury; do
  "$ROOT_DIR/synnergy" wallet new > "$ROOT_DIR/wallets/$name.wallet"
done

# Compute genesis wallet allocations
"$ROOT_DIR/synnergy" genesis allocate 1000000000 > "$ROOT_DIR/genesis_allocation.txt"

# Deploy bundled smart contracts
if ls "$ROOT_DIR/smart-contracts"/*.wasm > /dev/null 2>&1; then
  for wasm in "$ROOT_DIR/smart-contracts"/*.wasm; do
    echo "Deploying contract $wasm"
    "$ROOT_DIR/synnergy" contracts deploy --wasm "$wasm"
  done
fi

# Basic firewall rule
"$ROOT_DIR/synnergy" firewall allow 0.0.0.0/0

# Prepare storage using helper scripts if available
if [ -f "$ROOT_DIR/cmd/scripts/build_cli.sh" ] && [ -f "$ROOT_DIR/cmd/scripts/storage_pin.sh" ]; then
  (cd "$ROOT_DIR/cmd/scripts" && ./build_cli.sh && ./storage_pin.sh "$ROOT_DIR/README.md")
fi

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

# Setup initial liquidity pool and add bootstrap liquidity
"$ROOT_DIR/synnergy" liquidity_pools create SYN USD 30
"$ROOT_DIR/synnergy" liquidity_pools add SYN-USD bootstrapper 1000000 1000000

wait
