#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck disable=SC1091
source "$SCRIPT_DIR/lib/common.sh"

parse_common_flags "$@"
set -- "${POSITIONAL_ARGS[@]:-}"

usage() {
  cat <<'USAGE'
Bootstrap a local multi-node devnet that exercises the CLI, consensus engine,
virtual machine and wallet tooling.

Usage: devnet_start.sh [--count N] [--vm-mode MODE] [--demand FLOAT] [--stake FLOAT]
                       [--wallet-dir DIR] [--wallet-password PASS] [--topic TOPIC]
                       [common flags]

Options:
  --count N             Number of nodes to enrol in the swarm (default: 3)
  --vm-mode MODE        VM resource profile: heavy, light, superlight (default: heavy)
  --demand FLOAT        Initial network demand ratio for consensus weighting (default: 0.55)
  --stake FLOAT         Initial stake concentration ratio (default: 0.45)
  --wallet-dir DIR      Directory to store generated devnet wallets
                        (default: $PROJECT_ROOT/var/devnet/wallets)
  --wallet-password PW  Password used to encrypt generated wallets
                        (default: synnergy-devnet)
  --topic TOPIC         PubSub topic for the bootstrap announcement (default: devnet/announce)

Common flags:
  --dry-run             Print the operations without executing them
  --timeout SEC         Override command timeout window (default: 120s)
  --log-file PATH       Append logs to PATH instead of scripts/logs/
USAGE
}

COUNT=${SYN_DEVNET_COUNT:-3}
VM_MODE=${SYN_DEVNET_VM_MODE:-heavy}
CONSENSUS_DEMAND=${SYN_DEVNET_DEMAND:-0.55}
CONSENSUS_STAKE=${SYN_DEVNET_STAKE:-0.45}
WALLET_DIR=${SYN_DEVNET_WALLET_DIR:-$PROJECT_ROOT/var/devnet/wallets}
WALLET_PASSWORD=${SYN_DEVNET_WALLET_PASSWORD:-synnergy-devnet}
ANNOUNCE_TOPIC=${SYN_DEVNET_TOPIC:-devnet/announce}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --count)
      COUNT="$2"
      shift 2
      ;;
    --vm-mode)
      VM_MODE="$2"
      shift 2
      ;;
    --demand)
      CONSENSUS_DEMAND="$2"
      shift 2
      ;;
    --stake)
      CONSENSUS_STAKE="$2"
      shift 2
      ;;
    --wallet-dir)
      WALLET_DIR="$2"
      shift 2
      ;;
    --wallet-password)
      WALLET_PASSWORD="$2"
      shift 2
      ;;
    --topic)
      ANNOUNCE_TOPIC="$2"
      shift 2
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      if [[ "$1" =~ ^[0-9]+$ ]]; then
        COUNT="$1"
        shift
      else
        log_error "Unknown argument: $1"
        usage
        exit 1
      fi
      ;;
  esac
done

if [[ -z "$LOG_FILE" ]]; then
  set_log_file "$(basename "$0")"
fi

require_command go python3
ensure_directory "$WALLET_DIR"
log_info "Devnet bootstrap starting" "nodes=$COUNT" "vm_mode=$VM_MODE"

cleanup() {
  if [[ "$DRY_RUN" == true ]]; then
    return
  fi
  log_info "Tearing down devnet services"
  synnergy_cli simplevm stop >/dev/null 2>&1 || true
  synnergy_cli network stop >/dev/null 2>&1 || true
}
trap cleanup EXIT

validate_float() {
  local label="$1"
  local value="$2"
  if ! python3 - "$value" <<'PY'
import sys
try:
    float(sys.argv[1])
except ValueError:
    sys.exit(1)
PY
  then
    log_error "Invalid $label value: $value"
    exit 1
  fi
}

if ! [[ "$COUNT" =~ ^[0-9]+$ ]] || (( COUNT < 1 )); then
  log_error "--count expects a positive integer"
  exit 1
fi
validate_float demand "$CONSENSUS_DEMAND"
validate_float stake "$CONSENSUS_STAKE"

if [[ "$DRY_RUN" == true ]]; then
  log_info "[dry-run] would create VM ($VM_MODE), start network, enrol $COUNT nodes and generate wallets in $WALLET_DIR"
  exit 0
fi

synnergy_cli simplevm create "$VM_MODE"
synnergy_cli simplevm start
synnergy_cli network start

node_ids=()
for ((i = 1; i <= COUNT; i++)); do
  node_id="devnet-node-$i"
  synnergy_cli swarm join "$node_id"
  node_ids+=("$node_id")
  wallet_path="$WALLET_DIR/$node_id.wallet"
  synnergy_cli wallet new --out "$wallet_path" --password "$WALLET_PASSWORD"
  log_info "Provisioned wallet" "node=$node_id" "path=$wallet_path"
done

synnergy_cli consensus adjust "$CONSENSUS_DEMAND" "$CONSENSUS_STAKE"
weights_output="$(synnergy_cli --json consensus weights || true)"
python3 - "$weights_output" "$LOG_FILE" <<'PY'
import json
import sys
raw = sys.argv[1]
log_file = sys.argv[2]
payload = None
for idx, ch in enumerate(raw):
    if ch in '{[':
        try:
            payload = json.loads(raw[idx:])
        except json.JSONDecodeError:
            pass
        break
if payload:
    with open(log_file, 'a', encoding='utf-8') as fh:
        fh.write(f"consensus_weights={json.dumps(payload)}\n")
PY

synnergy_cli network broadcast "$ANNOUNCE_TOPIC" "bootstrap-complete"

log_info "Devnet ready" "peers=${node_ids[*]}" "wallet_dir=$WALLET_DIR"
printf 'Devnet started with %d nodes. Wallets stored in %s.\n' "$COUNT" "$WALLET_DIR"
