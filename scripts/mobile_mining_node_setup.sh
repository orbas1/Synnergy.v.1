#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib/common.sh
source "$SCRIPT_DIR/lib/common.sh"

usage() {
  cat <<'USAGE'
Provision a mobile mining node with adaptive power envelopes.

Usage: mobile_mining_node_setup.sh [--device-id ID] [--power LEVEL] [--start]
                                   [--duration SECONDS]

Options:
  --device-id ID     Identifier used for Function Web dashboards (default: mobile-1).
  --power LEVEL      Desired power level between 10 and 100 (default: 50).
  --start            Start the mobile miner immediately.
  --duration SEC     Runtime monitoring window (default: 15 seconds).
  -h, --help         Show usage.

The workflow uses CLI commands MobileMiningStart/MobileMiningSetPower to keep
consensus, wallet and VM telemetry aligned.
USAGE
}

DEVICE_ID="mobile-1"
POWER=50
START=false
DURATION=15

while [[ $# -gt 0 ]]; do
  case "$1" in
    --device-id)
      DEVICE_ID=$2
      shift 2
      ;;
    --power)
      POWER=$2
      shift 2
      ;;
    --start)
      START=true
      shift
      ;;
    --duration)
      DURATION=$2
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

if (( POWER < 10 || POWER > 100 )); then
  fail "power level must be between 10 and 100"
fi

require_commands jq

log_info "Configuring mobile miner $DEVICE_ID"

run_cli_json mobile-mining set-power "$POWER" >/dev/null

if [[ $START == true ]]; then
  log_info "Starting mobile mining"
  run_cli_json mobile-mining start >/dev/null
  trap 'run_cli_json mobile-mining stop >/dev/null || true' EXIT
fi

log_info "Streaming telemetry for $DURATION seconds"
elapsed=0
while [[ $elapsed -lt $DURATION ]]; do
  sleep 1
  stats=$(run_cli_json mobile-mining power)
  log_info "mobile miner snapshot: $stats"
  elapsed=$((elapsed + 1))
done

status=$(run_cli_json mobile-mining status)
log_info "Final status: $status"

log_info "Mobile mining node ready"
