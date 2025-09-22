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
Run an end-to-end smoke test that bootstraps the VM, consensus engine and
network stack via the Synnergy CLI. The test validates that swarm membership,
consensus weighting and transaction broadcast APIs remain healthy.

Usage: e2e_network_tests.sh [options]

Options:
  --nodes N               Number of swarm nodes to enrol (default: 3)
  --rounds N              Number of consensus rounds to attempt (default: 2)
  --demand FLOAT          Demand ratio used for weight adjustments (default: 0.6)
  --stake FLOAT           Stake ratio for weight adjustments (default: 0.5)
  --topic TOPIC           Topic used for broadcast verification (default: tests/e2e)
  -h, --help              Show this help text

Common flags:
  --dry-run               Print the planned CLI calls without executing them
  --timeout SEC           Override per-command timeout (default: 120)
  --log-file PATH         Append logs to PATH instead of scripts/logs/
USAGE
}

NODE_COUNT=${SYN_E2E_NODE_COUNT:-3}
ROUNDS=${SYN_E2E_ROUNDS:-2}
DEMAND=${SYN_E2E_DEMAND:-0.6}
STAKE=${SYN_E2E_STAKE:-0.5}
TOPIC=${SYN_E2E_TOPIC:-tests/e2e}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --nodes)
      NODE_COUNT="$2"
      shift 2
      ;;
    --rounds)
      ROUNDS="$2"
      shift 2
      ;;
    --demand)
      DEMAND="$2"
      shift 2
      ;;
    --stake)
      STAKE="$2"
      shift 2
      ;;
    --topic)
      TOPIC="$2"
      shift 2
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      log_error "Unknown argument: $1"
      usage
      exit 1
      ;;
  esac
done

if [[ -z "$LOG_FILE" ]]; then
  set_log_file "$(basename "$0")"
fi

require_command go python3

if ! [[ "$NODE_COUNT" =~ ^[0-9]+$ ]] || (( NODE_COUNT < 1 )); then
  log_error "--nodes expects a positive integer"
  exit 1
fi
if ! [[ "$ROUNDS" =~ ^[0-9]+$ ]] || (( ROUNDS < 1 )); then
  log_error "--rounds expects a positive integer"
  exit 1
fi

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

validate_float demand "$DEMAND"
validate_float stake "$STAKE"

if [[ "$DRY_RUN" == true ]]; then
  log_info "[dry-run] would create VM, start network, join $NODE_COUNT nodes, run $ROUNDS consensus rounds"
  exit 0
fi

cleanup() {
  synnergy_cli network stop >/dev/null 2>&1 || true
  synnergy_cli simplevm stop >/dev/null 2>&1 || true
}
trap cleanup EXIT

extract_json() {
  python3 - <<'PY'
import json
import sys
raw = sys.stdin.read()
for idx, ch in enumerate(raw):
    if ch in '{[':
        payload = raw[idx:]
        try:
            json.loads(payload)
        except json.JSONDecodeError:
            continue
        print(payload.strip())
        break
PY
}

cli_json() {
  local raw
  raw="$(synnergy_cli --json "$@" || true)"
  printf '%s' "$raw" | extract_json
}

synnergy_cli simplevm create heavy
synnergy_cli simplevm start
synnergy_cli network start

for ((i = 1; i <= NODE_COUNT; i++)); do
  synnergy_cli swarm join "e2e-node-$i"
done

peers_json="$(cli_json swarm peers)"
if [[ -z "$peers_json" ]]; then
  log_error "Failed to retrieve swarm peers"
  exit 1
fi
peer_count=$(printf '%s' "$peers_json" | python3 - <<'PY'
import json
import sys
peers = json.load(sys.stdin)
print(len(peers))
PY
)
if (( peer_count != NODE_COUNT )); then
  log_error "Expected $NODE_COUNT peers but found $peer_count"
  exit 1
fi

synnergy_cli consensus adjust "$DEMAND" "$STAKE"
weights_json="$(cli_json consensus weights)"
if [[ -z "$weights_json" ]]; then
  log_error "Failed to obtain consensus weights"
  exit 1
fi
printf '%s' "$weights_json" | python3 - <<'PY'
import json
import sys
weights = json.load(sys.stdin)
if abs(sum(weights.values()) - 1.0) > 0.25:
    sys.exit("consensus weights drifted outside expected range")
PY

for ((round = 1; round <= ROUNDS; round++)); do
  synnergy_cli swarm consensus
  synnergy_cli consensus mine 1 >/dev/null 2>&1 || true
done

threshold_json="$(cli_json consensus threshold "$DEMAND" "$STAKE")"
if [[ -z "$threshold_json" ]]; then
  log_error "Failed to compute transition threshold"
  exit 1
fi
threshold_value=$(printf '%s' "$threshold_json" | python3 - <<'PY'
import json
import sys
payload = json.load(sys.stdin)
print(float(payload.get("threshold", 0)))
PY
)
if ! python3 - "$threshold_value" <<'PY'; then
import sys
if float(sys.argv[1]) <= 0:
    sys.exit(1)
PY
  log_error "Consensus threshold must be positive"
  exit 1
fi

broadcast_json="$(cli_json network broadcast "$TOPIC" "test-message")"
if [[ -z "$broadcast_json" ]]; then
  log_warn "Broadcast command did not produce JSON output"
fi

log_info "E2E network checks completed" "peers=$peer_count" "rounds=$ROUNDS"
printf 'E2E network test passed with %d peers and %d rounds.\n' "$peer_count" "$ROUNDS"
