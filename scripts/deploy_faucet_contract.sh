#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck disable=SC1091
source "$SCRIPT_DIR/lib/common.sh"
set_log_file "$(basename "$0")"

usage() {
  cat <<'USAGE'
Deploy the canonical faucet template shipped with Synnergy.

Usage: deploy_faucet_contract.sh [options]

Options:
  --name TEMPLATE     Template to deploy (default: token_faucet)
  --owner ADDRESS     Owner address recorded for the faucet contract
  --gas UNITS         Gas limit override (default: 100000)
  --list              List available templates and exit
  --verify ADDRESS    Verify an existing faucet deployment by contract address
  --dry-run           Validate inputs without executing any CLI command
  --timeout SEC       Override command timeout (default: 120)
  --log-file FILE     Custom log destination
  -h, --help          Show this help message
USAGE
}

if [[ "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
  usage
  exit 0
fi

parse_common_flags "$@"
set -- "${POSITIONAL_ARGS[@]}"

NAME="token_faucet"
OWNER=""
GAS="100000"
LIST_ONLY=false
VERIFY_ADDR=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --name)
      NAME="${2:-}"
      shift 2
      ;;
    --owner)
      OWNER="${2:-}"
      shift 2
      ;;
    --gas)
      GAS="${2:-}"
      shift 2
      ;;
    --list)
      LIST_ONLY=true
      shift
      ;;
    --verify)
      VERIFY_ADDR="${2:-}"
      shift 2
      ;;
    *)
      log_error "Unknown argument: $1"
      usage
      exit 1
      ;;
  esac
done

if [[ "$LIST_ONLY" == true ]]; then
  synnergy_cli contracts list-templates
  exit 0
fi

if [[ -n "$VERIFY_ADDR" ]]; then
  if out=$(synnergy_cli contracts info "$VERIFY_ADDR"); then
    log_info "Contract $VERIFY_ADDR manifest:\n$out"
    exit 0
  else
    log_error "Unable to fetch contract metadata for $VERIFY_ADDR"
    exit 1
  fi
fi

if ! [[ "$GAS" =~ ^[0-9]+$ ]]; then
  log_error "Gas value must be numeric"
  exit 1
fi

log_info "Deploying faucet template $NAME"
args=(contracts deploy-template --name "$NAME" --gas "$GAS")
[[ -n "$OWNER" ]] && args+=(--owner "$OWNER")

if [[ "$DRY_RUN" == true ]]; then
  log_info "[dry-run] synnergy ${args[*]}"
  exit 0
fi

if address=$(synnergy_cli "${args[@]}"); then
  log_info "Faucet deployed at $address"
  printf '%s,%s,%s\n' "$NAME" "$address" "$(log_timestamp)" >>"$DEFAULT_LOG_DIR/faucet_contracts.csv"
else
  log_error "Faucet deployment failed"
  exit 1
fi
