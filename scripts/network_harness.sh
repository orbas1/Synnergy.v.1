#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib/common.sh
source "$SCRIPT_DIR/lib/common.sh"

usage() {
  cat <<'USAGE'
Run an integration harness exercising network, consensus and VM operations.

Usage: network_harness.sh [--cycles N] [--topic TOPIC]
                          [--payload PAYLOAD] [--stress]

Options:
  --cycles N      Number of broadcast/mine cycles to execute (default: 3).
  --topic TOPIC   Topic to use for broadcasts (default: harness).
  --payload DATA  Payload to mine for each cycle (default: harness-payload).
  --stress        Enable stress mode which issues additional CLI calls.
  -h, --help      Show usage.

The harness validates CLI wiring and ensures operations emit deterministic JSON
for Function Web ingestion.
USAGE
}

CYCLES=3
TOPIC="harness"
PAYLOAD="harness-payload"
STRESS=false

while [[ $# -gt 0 ]]; do
  case "$1" in
    --cycles)
      CYCLES=$2
      shift 2
      ;;
    --topic)
      TOPIC=$2
      shift 2
      ;;
    --payload)
      PAYLOAD=$2
      shift 2
      ;;
    --stress)
      STRESS=true
      shift
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

run_cli_json network start >/dev/null
run_cli_json mining start >/dev/null
trap 'run_cli_json network stop >/dev/null || true; run_cli_json mining stop >/dev/null || true' EXIT

for i in $(seq 1 "$CYCLES"); do
  log_info "Cycle $i broadcasting on $TOPIC"
  run_cli_json network broadcast "$TOPIC" "message-$i" >/dev/null
  run_cli_json node addtx "harness-$i" "receiver" "1" "1" "$i" >/dev/null
  run_cli_json mining mine "$PAYLOAD-$i" >/dev/null
  peers=$(run_cli_json network peers)
  log_info "Peers snapshot: $peers"
  if [[ $STRESS == true ]]; then
    log_info "Stress mode enabled - mining additional block"
    run_cli_json mining mine "stress-$i" >/dev/null
  fi
done

log_info "Network harness completed"
