#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck disable=SC1091
source "$SCRIPT_DIR/lib/common.sh"
set_log_file "$(basename "$0")"

usage() {
  cat <<'USAGE'
Deploy all starter smart contract templates bundled with Synnergy.

Usage: deploy_starter_smart_contracts_to_blockchain.sh [options]

Options:
  --templates LIST    Comma-separated template names to deploy
  --owner ADDRESS     Owner address applied to every deployment
  --gas UNITS         Gas limit per contract (default: 120000)
  --parallel N        Number of concurrent deployments (default: 1)
  --dry-run           Validate inputs without executing commands
  --timeout SEC       Override command timeout (default: 120)
  --log-file FILE     Custom log destination
  -h, --help          Show this message
USAGE
}

if [[ "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
  usage
  exit 0
fi

parse_common_flags "$@"
set -- "${POSITIONAL_ARGS[@]}"

OWNER=""
GAS="120000"
PARALLEL=1
TEMPLATE_LIST=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --templates)
      TEMPLATE_LIST="${2:-}"
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
    --parallel)
      PARALLEL="${2:-}"
      shift 2
      ;;
    *)
      log_error "Unknown argument: $1"
      usage
      exit 1
      ;;
  esac
done

if ! [[ "$GAS" =~ ^[0-9]+$ ]]; then
  log_error "Gas value must be numeric"
  exit 1
fi

if ! [[ "$PARALLEL" =~ ^[0-9]+$ ]] || ((PARALLEL < 1)); then
  log_error "--parallel must be a positive integer"
  exit 1
fi

mapfile -t templates < <(
  if [[ -n "$TEMPLATE_LIST" ]]; then
    tr ',' '\n' <<<"$TEMPLATE_LIST" | sed '/^\s*$/d'
  else
    shopt -s nullglob
    for file in "$PROJECT_ROOT"/smart-contracts/*.wasm; do
      basename "$file" .wasm
    done
  fi
)

if ((${#templates[@]} == 0)); then
  log_error "No templates discovered"
  exit 1
fi

log_info "Deploying ${#templates[@]} templates with gas=$GAS owner=${OWNER:-<none>}"

queue=()
for tpl in "${templates[@]}"; do
  args=(contracts deploy-template --name "$tpl" --gas "$GAS")
  [[ -n "$OWNER" ]] && args+=(--owner "$OWNER")
  if [[ "$DRY_RUN" == true ]]; then
    log_info "[dry-run] synnergy ${args[*]}"
    continue
  fi

  (
    if address=$(synnergy_cli "${args[@]}"); then
      log_info "$tpl deployed at $address"
      printf '%s,%s,%s\n' "$tpl" "$address" "$(log_timestamp)" >>"$DEFAULT_LOG_DIR/starter_contracts.csv"
    else
      log_error "Deployment failed for $tpl"
      exit 1
    fi
  ) &
  queue+=($!)

  if (( ${#queue[@]} >= PARALLEL )); then
    wait "${queue[0]}"
    queue=(${queue[@]:1})
  fi
done

for pid in "${queue[@]}"; do
  wait "$pid"
done
