#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib/common.sh
source "$SCRIPT_DIR/lib/common.sh"

usage() {
  cat <<'USAGE'
Bootstrap a local Synnergy cluster using CLI orchestration.

Usage: multi_node_cluster_setup.sh [--nodes N] [--start-network] [--topology FILE]
                                   [--output FILE]

Options:
  --nodes N           Number of logical nodes to prepare (default: 3).
  --start-network     Start the network service and publish discovery information.
  --topology FILE     Optional JSON file describing peer links; generated if absent.
  --output FILE       Persist resulting peer list for dashboards.
  -h, --help          Show this help message.

The script initialises wallets, stakes them, provisions mining/mobile nodes and
ensures consensus and VM modules are live for Function Web integration.
USAGE
}

NODE_COUNT=3
START_NETWORK=false
TOPOLOGY_FILE=""
OUTPUT=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --nodes)
      NODE_COUNT=$2
      shift 2
      ;;
    --start-network)
      START_NETWORK=true
      shift
      ;;
    --topology)
      TOPOLOGY_FILE=$2
      shift 2
      ;;
    --output)
      OUTPUT=$2
      shift 2
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      usage
      fail "unknown flag $1"
      ;;
  esac
done

require_commands jq openssl

log_info "Preparing $NODE_COUNT nodes"

addresses=()
for i in $(seq 1 "$NODE_COUNT"); do
  pw=$(openssl rand -hex 8)
  wallet_path="${PWD}/wallets/node-$i.json"
  mkdir -p "$(dirname "$wallet_path")"
  wallet_json=$(run_cli_json wallet new --out "$wallet_path" --password "$pw")
  addr=$(json_extract '.address' "$wallet_json")
  log_info "Wallet $addr created for node $i"
  addresses+=("$addr")
  run_cli_json node stake "$addr" "1000" >/dev/null
  if (( i == 1 )); then
    run_cli_json mining start >/dev/null
  else
    run_cli_json mobile-mining set-power "$((50 + i))" >/dev/null
  fi
  run_cli_json node addtx "$addr" "$addr" "1" "1" "$i" >/dev/null
  run_cli_json mining mine "bootstrap-$i" >/dev/null
  log_info "Node $i primed"
done

if [[ $START_NETWORK == true ]]; then
  log_info "Starting network stack"
  run_cli_json network start >/dev/null
fi

if [[ -z $TOPOLOGY_FILE ]]; then
  TOPOLOGY_FILE="${PWD}/cluster-topology.json"
fi

log_info "Generating topology file $TOPOLOGY_FILE"
cat >"$TOPOLOGY_FILE" <<JSON
{
  "nodes": [
    $(printf '"%s",' "${addresses[@]}" | sed 's/,$//')
  ],
  "timestamp": "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
}
JSON
chmod 600 "$TOPOLOGY_FILE"

peers_json=$(run_cli_json network peers)
log_info "Peers: $peers_json"

if [[ -n $OUTPUT ]]; then
  write_secure_file "$OUTPUT" <<JSON
$peers_json
JSON
fi

log_info "Cluster setup complete"
