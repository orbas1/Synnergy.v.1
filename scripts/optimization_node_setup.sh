#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib/common.sh
source "$SCRIPT_DIR/lib/common.sh"

usage() {
  cat <<'USAGE'
Optimise transaction ordering using the optimisation node CLI module.

Usage: optimization_node_setup.sh [--tx HASH:FEE:SIZE ...] [--output FILE]

Options:
  --tx VALUE       Transaction spec formatted as hash:fee:size. Multiple allowed.
  --output FILE    Persist optimisation results to file.
  -h, --help       Show help.

If no transactions are supplied, a synthetic batch is generated for benchmarking.
USAGE
}

TX_SPECS=()
OUTPUT=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --tx)
      TX_SPECS+=("$2")
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

if [[ ${#TX_SPECS[@]} -eq 0 ]]; then
  TX_SPECS=("tx1:30:250" "tx2:10:128" "tx3:45:512")
fi

require_commands jq

log_info "Running optimisation for ${#TX_SPECS[@]} transactions"
result=$(run_cli_json optimize fee "${TX_SPECS[@]}")
log_info "Optimisation result: $result"

if [[ -n $OUTPUT ]]; then
  write_secure_file "$OUTPUT" <<JSON
$result
JSON
fi

log_info "Optimisation workflow completed"
