#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib/common.sh
source "$SCRIPT_DIR/lib/common.sh"

usage() {
  cat <<'USAGE'
General node provisioning pipeline.

Usage: node_setup.sh [--address ADDRESS] [--stake AMOUNT] [--role ROLE]
                     [--start-mining]

Options:
  --address ADDRESS  Optional existing address to use; otherwise generated.
  --stake AMOUNT     Stake amount (default: 1000).
  --role ROLE        Role to grant through access CLI (default: validator).
  --start-mining     Start mining after provisioning.
  -h, --help         Show this help message.
USAGE
}

ADDRESS=""
STAKE=1000
ROLE="validator"
START_MINING=false

while [[ $# -gt 0 ]]; do
  case "$1" in
    --address)
      ADDRESS=$2
      shift 2
      ;;
    --stake)
      STAKE=$2
      shift 2
      ;;
    --role)
      ROLE=$2
      shift 2
      ;;
    --start-mining)
      START_MINING=true
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

require_commands jq openssl

if [[ -z $ADDRESS ]]; then
  log_info "Generating new wallet"
  wallet_json=$(run_cli_json wallet new)
  ADDRESS=$(json_extract '.address' "$wallet_json")
fi

log_info "Staking address $ADDRESS"
run_cli_json node stake "$ADDRESS" "$STAKE" >/dev/null

log_info "Granting role $ROLE"
run_cli_json access grant "$ROLE" "$ADDRESS" >/dev/null

if [[ $START_MINING == true ]]; then
  log_info "Starting mining worker"
  run_cli_json mining start >/dev/null
fi

log_info "Node setup complete for $ADDRESS"
