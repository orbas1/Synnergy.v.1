#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck disable=SC1091
source "$SCRIPT_DIR/lib/common.sh"
set_log_file "$(basename "$0")"

usage() {
  cat <<'USAGE'
Deploy an arbitrary WASM smart contract via the Synnergy CLI.

Usage: deploy_contract.sh --wasm FILE [options]

Options:
  --wasm FILE        Path to the compiled contract artifact (.wasm)
  --ric FILE         Optional Ricardian contract/metadata file
  --owner ADDRESS    Owner address that will administer the deployment
  --gas UNITS        Gas limit to allocate (default: 100000)
  --label NAME       Friendly label stored in deployment ledger logs
  --dry-run          Validate inputs and print actions without executing
  --timeout SEC      Override command timeout (default: 120)
  --log-file FILE    Custom log file destination
  -h, --help         Show this help message

Environment variables:
  SYN_CLI_BIN        Explicit path to synnergy binary (defaults to go run)
  SYN_ENV_FILE       Optional .env file to source before running
USAGE
}

if [[ "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
  usage
  exit 0
fi

load_env_file "${SYN_ENV_FILE:-}"

parse_common_flags "$@"
set -- "${POSITIONAL_ARGS[@]}"

WASM=""
RIC=""
OWNER=""
GAS="100000"
LABEL=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --wasm)
      WASM="${2:-}"
      shift 2
      ;;
    --ric)
      RIC="${2:-}"
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
    --label)
      LABEL="${2:-}"
      shift 2
      ;;
    -*)
      log_error "Unknown argument: $1"
      usage
      exit 1
      ;;
    *)
      if [[ -z "$WASM" ]]; then
        WASM="$1"
        shift
        continue
      fi
      log_error "Unknown argument: $1"
      usage
      exit 1
      ;;
  esac
done

if [[ -z "$WASM" ]]; then
  log_error "--wasm is required"
  usage
  exit 1
fi

if [[ ! -f "$WASM" ]]; then
  log_error "contract file not found: $WASM"
  exit 1
fi

if [[ -n "$RIC" && ! -f "$RIC" ]]; then
  log_error "Ricardian manifest not found: $RIC"
  exit 1
fi

if ! [[ "$GAS" =~ ^[0-9]+$ ]]; then
  log_error "Gas value must be numeric"
  exit 1
fi

log_info "Preparing to deploy $WASM with gas=$GAS owner=${OWNER:-<none>}"

args=(contracts deploy --wasm "$WASM" --gas "$GAS")
[[ -n "$RIC" ]] && args+=(--ric "$RIC")
[[ -n "$OWNER" ]] && args+=(--owner "$OWNER")

if [[ "$DRY_RUN" == true ]]; then
  log_info "[dry-run] synnergy ${args[*]}"
  exit 0
fi

if address=$(synnergy_cli "${args[@]}"); then
  log_info "Contract deployed at $address"
  if [[ -n "$LABEL" ]]; then
    log_info "Recording label $LABEL"
    printf '%s,%s,%s\n' "$LABEL" "$address" "$(log_timestamp)" >>"$DEFAULT_LOG_DIR/contracts.csv"
  fi
else
  log_error "Contract deployment failed"
  exit 1
fi
