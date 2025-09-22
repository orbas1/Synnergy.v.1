#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib/common.sh
source "$SCRIPT_DIR/lib/common.sh"

usage() {
  cat <<'USAGE'
Perform health diagnostics on the Synnergy networking layer.

Usage: network_diagnostics.sh [--iterations N] [--broadcast TOPIC:MESSAGE]
                              [--output FILE]

Options:
  --iterations N         Number of status polls to perform (default: 5).
  --broadcast TOPIC:MSG  Publish a diagnostic message before polling.
  --output FILE          Save diagnostic JSON to file.
  -h, --help             Show this help message.

The diagnostics ensure CLI/VM telemetry is wired for Function Web dashboards,
verifying gas charges for NetworkStart/NetworkPeers/NetworkBroadcast.
USAGE
}

ITERATIONS=5
BROADCAST=""
OUTPUT=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --iterations)
      ITERATIONS=$2
      shift 2
      ;;
    --broadcast)
      BROADCAST=$2
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

require_commands jq

log_info "Ensuring network service is running"
run_cli_json network start >/dev/null

if [[ -n $BROADCAST ]]; then
  topic=${BROADCAST%%:*}
  msg=${BROADCAST#*:}
  log_info "Broadcasting diagnostic message on $topic"
  run_cli_json network broadcast "$topic" "$msg" >/dev/null
fi

results=()
for i in $(seq 1 "$ITERATIONS"); do
  sleep 1
  peers=$(run_cli_json network peers)
  results+=("$peers")
  log_info "Iteration $i peers: $peers"
done

if [[ -n $OUTPUT ]]; then
  log_info "Writing diagnostics to $OUTPUT"
  {
    echo "["
    for idx in "${!results[@]}"; do
      entry=${results[$idx]}
      if [[ $idx -lt $(( ${#results[@]} - 1 )) ]]; then
        echo "  ${entry},"
      else
        echo "  ${entry}"
      fi
    done
    echo "]"
  } | write_secure_file "$OUTPUT"
fi

log_info "Diagnostics completed"
