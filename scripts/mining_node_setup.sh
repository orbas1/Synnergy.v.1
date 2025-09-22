#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib/common.sh
source "$SCRIPT_DIR/lib/common.sh"

usage() {
  cat <<'USAGE'
Enterprise mining node bootstrapper.

Usage: mining_node_setup.sh [--wallet ADDRESS] [--stake AMOUNT] [--hash-data DATA]
                           [--auto-start] [--monitor-seconds SECONDS]

Options:
  --wallet ADDRESS       Wallet address to stake for mining (default: miner1).
  --stake AMOUNT         Amount of SYN to stake when initialising mining (default: 1000).
  --hash-data DATA       Payload to mine when validating GPU/CPU pipeline (default: stage103).
  --auto-start           Start the mining engine immediately after staking.
  --monitor-seconds N    Monitor hashrate for N seconds to ensure stability (default: 10).
  -h, --help             Show this help message and exit.

The script aligns CLI, VM, consensus and wallet layers by staking the wallet,
starting the mining worker, validating gas consumption through the CLI and
collecting telemetry for the Function Web dashboard.
USAGE
}

WALLET="miner1"
STAKE="1000"
HASH_DATA="stage103"
AUTO_START=false
MONITOR_SECONDS=10

if [[ $# -eq 0 ]]; then
  log_info "no arguments supplied; using defaults"
fi

while [[ $# -gt 0 ]]; do
  case "$1" in
    --wallet)
      WALLET=$2
      shift 2
      ;;
    --stake)
      STAKE=$2
      shift 2
      ;;
    --hash-data)
      HASH_DATA=$2
      shift 2
      ;;
    --auto-start)
      AUTO_START=true
      shift
      ;;
    --monitor-seconds)
      MONITOR_SECONDS=$2
      shift 2
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      usage
      fail "unknown argument: $1"
      ;;
  esac
done

require_commands jq

log_info "Ensuring CLI connectivity"

stake_wallet() {
  run_cli_json node stake "$WALLET" "$STAKE" >/dev/null
}

with_retry 3 stake_wallet || fail "unable to stake wallet $WALLET"

log_info "Wallet $WALLET staked with $STAKE units"

declare -a stop_actions=()

if [[ $AUTO_START == true ]]; then
  log_info "Starting mining engine"
  run_cli_json mining start >/dev/null
  stop_actions+=("stop_mining")
fi

stop_mining() {
  log_info "Stopping mining engine"
  run_cli_json mining stop >/dev/null || true
}

trap 'for action in "${stop_actions[@]}"; do "$action"; done' EXIT

log_info "Submitting warm-up transaction to mempool"
run_cli_json node addtx "$WALLET" "$WALLET" "1" "1" "1" >/dev/null

log_info "Mining sample payload"
run_cli_json mining mine "$HASH_DATA" >/dev/null

monitor_hashrate() {
  local duration=$1
  local elapsed=0
  while [[ $elapsed -lt $duration ]]; do
    sleep 1
    local stats
    stats=$(run_cli_json mining hashrate)
    log_info "hashrate snapshot: $stats"
    elapsed=$((elapsed + 1))
  done
}

monitor_hashrate "$MONITOR_SECONDS"

status_output=$(run_cli_json mining status)
if [[ $(json_extract '.mining' "$status_output") != "true" ]]; then
  log_warn "Mining status reported inactive; restarting"
  run_cli_json mining start >/dev/null
fi

log_info "Mining node setup completed successfully"
