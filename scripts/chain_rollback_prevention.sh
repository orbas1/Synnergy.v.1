#!/usr/bin/env bash

set -euo pipefail
IFS=$'\n\t'

SCRIPT_NAME=$(basename "$0")

usage() {
  cat <<'USAGE'
Usage: chain_rollback_prevention.sh --network <name> --checkpoint-file <path> [options]

Enterprise-grade rollback guard for Synnergy Network chains. The guard compares the
current head/finality status with a persisted checkpoint and initiates mitigation
when a rollback is detected.

Required arguments:
  --network NAME            Logical network name (e.g., mainnet, testnet).
  --checkpoint-file PATH    JSON checkpoint written by previous guard runs.

Optional arguments:
  --status-file PATH        Pre-recorded chain status JSON (used for tests and dry-runs).
  --cli PATH                Explicit Synnergy CLI binary used to fetch status/mitigate.
                             Defaults to $SYN_CLI, then synnergy, then synnergy-cli.
  --endpoint URL            Override API endpoint for CLI calls.
  --profile NAME            CLI profile (default: production).
  --mitigation-action ACT   Mitigation strategy when rollback occurs: none, quarantine,
                             or halt (default: quarantine).
  --allow-finality-drift N  Allow checkpoint to exceed current finality by up to N blocks
                             before classifying as rollback (default: 0).
  --max-finality-lag N      Maximum tolerated difference between head and finalized height
                             before flagging degraded finality (default: 500).
  --lock-file PATH          File used for flock-based mutual exclusion.
  --metrics-file PATH       Emit JSON metrics describing the guard evaluation.
  --hook PATH               Executable invoked after evaluation. Receives the evaluation
                             result (ok|rollback|degraded) as its first argument.
  --dry-run                 Skip CLI invocations while still performing validations.
  --no-update-checkpoint    Do not update the checkpoint even when evaluation succeeds.
  --help, -h                Display this help message.

The guard expects chain status JSON to provide the following fields:
  network, height, hash, finalized_height, finalized_hash

Examples:
  # Validate using recorded status and update checkpoint.
  chain_rollback_prevention.sh --network mainnet \
    --checkpoint-file /var/lib/synnergy/rollback/checkpoint.json \
    --status-file /tmp/status.json --mitigation-action none

  # Production invocation that fetches live status and mitigates automatically.
  chain_rollback_prevention.sh --network mainnet --checkpoint-file /var/lib/synnergy/rollback/checkpoint.json \
    --lock-file /var/run/synnergy/rollback.lock --metrics-file /var/log/synnergy/rollback.json
USAGE
}

log() {
  local level=$1
  local message=$2
  printf '%s [%s] %s\n' "$(date -u +'%Y-%m-%dT%H:%M:%SZ')" "$level" "$message" >&2
}

fatal() {
  local message=$1
  local code=${2:-1}
  log "ERROR" "$message"
  exit "$code"
}

cleanup() {
  local status=$?
  if [[ -n "${TMP_STATUS_FILE:-}" && -f "${TMP_STATUS_FILE}" ]]; then
    rm -f "$TMP_STATUS_FILE"
  fi
  if [[ -n "${LOCK_FD:-}" ]]; then
    eval "exec ${LOCK_FD}>&-"
  fi
  if [[ $status -ne 0 ]]; then
    log "ERROR" "${SCRIPT_NAME} terminated with status $status"
  fi
}
trap cleanup EXIT
trap 'fatal "${SCRIPT_NAME} interrupted" 2' INT TERM

require_command() {
  local cmd=$1
  local description=${2:-$1}
  if ! command -v "$cmd" >/dev/null 2>&1; then
    fatal "$description is required but not installed"
  fi
}

resolve_cli() {
  local candidate=$1
  if [[ -z "$candidate" ]]; then
    if [[ -n "${SYN_CLI:-}" ]]; then
      candidate=$SYN_CLI
    elif command -v synnergy >/dev/null 2>&1; then
      candidate=$(command -v synnergy)
    elif command -v synnergy-cli >/dev/null 2>&1; then
      candidate=$(command -v synnergy-cli)
    else
      fatal "CLI binary not found. Provide --cli or set SYN_CLI"
    fi
  fi

  if [[ -x "$candidate" ]]; then
    echo "$candidate"
    return 0
  fi

  if command -v "$candidate" >/dev/null 2>&1; then
    command -v "$candidate"
    return 0
  fi

  fatal "CLI binary '$candidate' not found"
}

acquire_lock() {
  local path=$1
  if [[ -z "$path" ]]; then
    return
  fi

  mkdir -p "$(dirname "$path")"
  LOCK_FD=212
  eval "exec ${LOCK_FD}>\"$path\""
  if ! flock -n "$LOCK_FD"; then
    fatal "another rollback guard instance already holds the lock at $path"
  fi
}

parse_json_file() {
  local file=$1
  local prefix=$2
  local stdout
  local stderr
  stderr=$(mktemp)
  if ! stdout=$(python3 - "$file" "$prefix" <<'PY' 2>"$stderr"
import json
import shlex
import sys

file_path = sys.argv[1]
prefix = sys.argv[2]
required = ["network", "height", "hash", "finalized_height", "finalized_hash"]
try:
    with open(file_path, "r", encoding="utf-8") as fh:
        data = json.load(fh)
except Exception as exc:  # pragma: no cover - handled in shell
    print(f"unable to read {file_path}: {exc}", file=sys.stderr)
    sys.exit(2)
missing = [key for key in required if key not in data]
if missing:
    print(f"{file_path} missing required keys: {', '.join(missing)}", file=sys.stderr)
    sys.exit(3)
for key in required:
    value = data[key]
    if isinstance(value, bool):
        value = "1" if value else "0"
    else:
        value = str(value)
    print(f"{prefix}{key.upper()}={shlex.quote(value)}")
optional_keys = ["timestamp", "finality_threshold", "finalized_by"]
for key in optional_keys:
    if key in data:
        value = data[key]
        if isinstance(value, bool):
            value = "1" if value else "0"
        else:
            value = str(value)
        print(f"{prefix}{key.upper()}={shlex.quote(value)}")
PY
  ); then
    local message
    message=$(<"$stderr")
    rm -f "$stderr"
    fatal "failed to parse $file: ${message:-unknown error}"
  fi
  rm -f "$stderr"
  eval "$stdout"
}

write_checkpoint() {
  local file=$1
  local network=$2
  local head_height=$3
  local head_hash=$4
  local finalized_height=$5
  local finalized_hash=$6
  local timestamp
  timestamp=$(date -u +'%Y-%m-%dT%H:%M:%SZ')
  mkdir -p "$(dirname "$file")"
  cat >"$file" <<JSON
{
  "network": "${network}",
  "height": ${head_height},
  "hash": "${head_hash}",
  "finalized_height": ${finalized_height},
  "finalized_hash": "${finalized_hash}",
  "timestamp": "${timestamp}"
}
JSON
}

emit_metrics() {
  local file=$1
  local network=$2
  local head_height=$3
  local head_hash=$4
  local finalized_height=$5
  local finalized_hash=$6
  local checkpoint_finalized_height=$7
  local checkpoint_finalized_hash=$8
  local rollback_detected=$9
  local mitigation=${10}
  local dry_run=${11}
  local updated_checkpoint=${12}
  local degraded=${13}
  local initialized=${14}

  if [[ -z "$file" ]]; then
    return
  fi

  mkdir -p "$(dirname "$file")"
  cat >"$file" <<JSON
{
  "network": "${network}",
  "head_height": ${head_height},
  "head_hash": "${head_hash}",
  "finalized_height": ${finalized_height},
  "finalized_hash": "${finalized_hash}",
  "checkpoint_finalized_height": ${checkpoint_finalized_height},
  "checkpoint_finalized_hash": "${checkpoint_finalized_hash}",
  "rollback_detected": ${rollback_detected},
  "mitigation": "${mitigation}",
  "dry_run": ${dry_run},
  "updated_checkpoint": ${updated_checkpoint},
  "degraded_finality": ${degraded},
  "initialized_checkpoint": ${initialized},
  "evaluated_at": "$(date -u +'%Y-%m-%dT%H:%M:%SZ')"
}
JSON
}

invoke_cli() {
  local dry_run_flag=$1
  shift
  if [[ "$dry_run_flag" -eq 1 ]]; then
    log "INFO" "dry-run: would invoke CLI $*"
    return 0
  fi
  "$@"
}

run_hook() {
  local result=$1
  if [[ -z "${HOOK_PATH:-}" ]]; then
    return
  fi
  if [[ ! -x "$HOOK_PATH" ]]; then
    fatal "hook $HOOK_PATH is not executable"
  fi
  ROLLBACK_RESULT=$result \
  ROLLBACK_DETECTED=$ROLLBACK_DETECTED \
  ROLLBACK_REASON="$ROLLBACK_REASON" \
  ROLLBACK_NETWORK="$NETWORK" \
  ROLLBACK_FINALIZED_HEIGHT="$STATUS_FINALIZED_HEIGHT" \
  ROLLBACK_FINALIZED_HASH="$STATUS_FINALIZED_HASH" \
  ROLLBACK_HEAD_HEIGHT="$STATUS_HEAD_HEIGHT" \
  ROLLBACK_HEAD_HASH="$STATUS_HEAD_HASH" \
  ROLLBACK_CHECKPOINT="$CHECKPOINT_FILE" \
  "$HOOK_PATH" "$result"
}

NETWORK=""
CHECKPOINT_FILE=""
STATUS_FILE=""
CLI_BINARY=""
ENDPOINT=""
PROFILE="production"
MITIGATION_ACTION="quarantine"
ALLOW_FINALITY_DRIFT=0
MAX_FINALITY_LAG=500
LOCK_PATH=""
METRICS_FILE=""
HOOK_PATH=""
DRY_RUN=0
UPDATE_CHECKPOINT=1

if [[ $# -eq 0 ]]; then
  usage
  exit 1
fi

while [[ $# -gt 0 ]]; do
  case "$1" in
    --network)
      NETWORK=${2:-}
      shift 2
      ;;
    --checkpoint-file)
      CHECKPOINT_FILE=${2:-}
      shift 2
      ;;
    --status-file)
      STATUS_FILE=${2:-}
      shift 2
      ;;
    --cli)
      CLI_BINARY=${2:-}
      shift 2
      ;;
    --endpoint)
      ENDPOINT=${2:-}
      shift 2
      ;;
    --profile)
      PROFILE=${2:-}
      shift 2
      ;;
    --mitigation-action)
      MITIGATION_ACTION=${2:-}
      shift 2
      ;;
    --allow-finality-drift)
      ALLOW_FINALITY_DRIFT=${2:-0}
      shift 2
      ;;
    --max-finality-lag)
      MAX_FINALITY_LAG=${2:-500}
      shift 2
      ;;
    --lock-file)
      LOCK_PATH=${2:-}
      shift 2
      ;;
    --metrics-file)
      METRICS_FILE=${2:-}
      shift 2
      ;;
    --hook)
      HOOK_PATH=${2:-}
      shift 2
      ;;
    --dry-run)
      DRY_RUN=1
      shift
      ;;
    --no-update-checkpoint)
      UPDATE_CHECKPOINT=0
      shift
      ;;
    --help|-h)
      usage
      exit 0
      ;;
    *)
      fatal "unknown argument: $1"
      ;;
  esac
done

if [[ -z "$NETWORK" ]]; then
  fatal "--network is required"
fi
if [[ -z "$CHECKPOINT_FILE" ]]; then
  fatal "--checkpoint-file is required"
fi
if ! [[ "$ALLOW_FINALITY_DRIFT" =~ ^[0-9]+$ ]]; then
  fatal "--allow-finality-drift must be a non-negative integer"
fi
if ! [[ "$MAX_FINALITY_LAG" =~ ^[0-9]+$ ]]; then
  fatal "--max-finality-lag must be a non-negative integer"
fi

require_command python3 "python3"
acquire_lock "$LOCK_PATH"

TMP_STATUS_FILE=$(mktemp)

if [[ -n "$STATUS_FILE" ]]; then
  if [[ ! -f "$STATUS_FILE" ]]; then
    fatal "status file '$STATUS_FILE' not found"
  fi
  cp "$STATUS_FILE" "$TMP_STATUS_FILE"
else
  CLI_BINARY=$(resolve_cli "$CLI_BINARY")
  STATUS_CMD=("$CLI_BINARY" "chain" "status" "--network" "$NETWORK" "--output" "json" "--profile" "$PROFILE")
  if [[ -n "$ENDPOINT" ]]; then
    STATUS_CMD+=("--endpoint" "$ENDPOINT")
  fi
  log "INFO" "collecting chain status via CLI"
  if ! invoke_cli "$DRY_RUN" "${STATUS_CMD[@]}" >"$TMP_STATUS_FILE"; then
    fatal "failed to collect chain status"
  fi
fi

parse_json_file "$TMP_STATUS_FILE" "STATUS_"

STATUS_HEAD_HEIGHT=${STATUS_HEIGHT}
STATUS_HEAD_HASH=${STATUS_HASH}
STATUS_FINALIZED_HEIGHT=${STATUS_FINALIZED_HEIGHT}
STATUS_FINALIZED_HASH=${STATUS_FINALIZED_HASH}

if [[ "$STATUS_NETWORK" != "$NETWORK" ]]; then
  fatal "status network '$STATUS_NETWORK' does not match requested network '$NETWORK'"
fi

FINALITY_LAG=$(( STATUS_HEAD_HEIGHT - STATUS_FINALIZED_HEIGHT ))
if (( FINALITY_LAG < 0 )); then
  FINALITY_LAG=$(( -1 * FINALITY_LAG ))
fi

ROLLBACK_DETECTED=0
ROLLBACK_REASON=""
DEGRADED_FINALITY=0
INITIALIZED_CHECKPOINT=0
UPDATED_CHECKPOINT=0

if [[ ! -f "$CHECKPOINT_FILE" ]]; then
  log "INFO" "no checkpoint found; initializing at finalized height ${STATUS_FINALIZED_HEIGHT}"
  write_checkpoint "$CHECKPOINT_FILE" "$NETWORK" "$STATUS_HEAD_HEIGHT" "$STATUS_HEAD_HASH" "$STATUS_FINALIZED_HEIGHT" "$STATUS_FINALIZED_HASH"
  UPDATED_CHECKPOINT=1
  INITIALIZED_CHECKPOINT=1
  emit_metrics "$METRICS_FILE" "$NETWORK" "$STATUS_HEAD_HEIGHT" "$STATUS_HEAD_HASH" "$STATUS_FINALIZED_HEIGHT" "$STATUS_FINALIZED_HASH" "$STATUS_FINALIZED_HEIGHT" "$STATUS_FINALIZED_HASH" "$ROLLBACK_DETECTED" "$MITIGATION_ACTION" "$DRY_RUN" "$UPDATED_CHECKPOINT" "$DEGRADED_FINALITY" "$INITIALIZED_CHECKPOINT"
  log "INFO" "checkpoint initialized at height ${STATUS_FINALIZED_HEIGHT}"
  run_hook "ok"
  printf 'CHECKPOINT_PATH=%s\n' "$CHECKPOINT_FILE"
  exit 0
fi

parse_json_file "$CHECKPOINT_FILE" "CHECKPOINT_"

if [[ "$CHECKPOINT_NETWORK" != "$NETWORK" ]]; then
  fatal "checkpoint network '$CHECKPOINT_NETWORK' does not match requested network '$NETWORK'"
fi

CHECKPOINT_FINALIZED_HEIGHT=${CHECKPOINT_FINALIZED_HEIGHT}
CHECKPOINT_FINALIZED_HASH=${CHECKPOINT_FINALIZED_HASH}

if (( STATUS_FINALIZED_HEIGHT + ALLOW_FINALITY_DRIFT < CHECKPOINT_FINALIZED_HEIGHT )); then
  ROLLBACK_DETECTED=1
  ROLLBACK_REASON="finalized height ${STATUS_FINALIZED_HEIGHT} behind checkpoint ${CHECKPOINT_FINALIZED_HEIGHT}"
elif (( STATUS_FINALIZED_HEIGHT == CHECKPOINT_FINALIZED_HEIGHT )) && [[ "$STATUS_FINALIZED_HASH" != "$CHECKPOINT_FINALIZED_HASH" ]]; then
  ROLLBACK_DETECTED=1
  ROLLBACK_REASON="finalized hash mismatch at height ${STATUS_FINALIZED_HEIGHT}"
fi

if (( FINALITY_LAG > MAX_FINALITY_LAG )); then
  DEGRADED_FINALITY=1
fi

if (( ROLLBACK_DETECTED == 1 )); then
  log "ERROR" "rollback detected: ${ROLLBACK_REASON}"
  emit_metrics "$METRICS_FILE" "$NETWORK" "$STATUS_HEAD_HEIGHT" "$STATUS_HEAD_HASH" "$STATUS_FINALIZED_HEIGHT" "$STATUS_FINALIZED_HASH" "$CHECKPOINT_FINALIZED_HEIGHT" "$CHECKPOINT_FINALIZED_HASH" "$ROLLBACK_DETECTED" "$MITIGATION_ACTION" "$DRY_RUN" "$UPDATED_CHECKPOINT" "$DEGRADED_FINALITY" "$INITIALIZED_CHECKPOINT"
  run_hook "rollback"
  case "$MITIGATION_ACTION" in
    none)
      log "WARN" "mitigation disabled; manual intervention required"
      ;;
    quarantine)
      CLI_BINARY=$(resolve_cli "$CLI_BINARY")
      invoke_cli "$DRY_RUN" "$CLI_BINARY" "consensus" "quarantine" "--network" "$NETWORK" "--reason" "rollback-detected"
      ;;
    halt)
      CLI_BINARY=$(resolve_cli "$CLI_BINARY")
      invoke_cli "$DRY_RUN" "$CLI_BINARY" "consensus" "halt" "--network" "$NETWORK" "--reason" "rollback-detected"
      ;;
    *)
      fatal "unknown mitigation action: $MITIGATION_ACTION"
      ;;
  esac
  exit 3
fi

if (( DEGRADED_FINALITY == 1 )); then
  log "WARN" "finality lag ${FINALITY_LAG} exceeds threshold ${MAX_FINALITY_LAG}"
  run_hook "degraded"
fi

if (( UPDATE_CHECKPOINT == 1 )); then
  write_checkpoint "$CHECKPOINT_FILE" "$NETWORK" "$STATUS_HEAD_HEIGHT" "$STATUS_HEAD_HASH" "$STATUS_FINALIZED_HEIGHT" "$STATUS_FINALIZED_HASH"
  UPDATED_CHECKPOINT=1
  log "INFO" "checkpoint updated to finalized height ${STATUS_FINALIZED_HEIGHT}"
fi

emit_metrics "$METRICS_FILE" "$NETWORK" "$STATUS_HEAD_HEIGHT" "$STATUS_HEAD_HASH" "$STATUS_FINALIZED_HEIGHT" "$STATUS_FINALIZED_HASH" "$CHECKPOINT_FINALIZED_HEIGHT" "$CHECKPOINT_FINALIZED_HASH" "$ROLLBACK_DETECTED" "$MITIGATION_ACTION" "$DRY_RUN" "$UPDATED_CHECKPOINT" "$DEGRADED_FINALITY" "$INITIALIZED_CHECKPOINT"
run_hook "ok"

printf 'STATUS_FINALIZED_HEIGHT=%s\n' "$STATUS_FINALIZED_HEIGHT"
printf 'STATUS_FINALIZED_HASH=%s\n' "$STATUS_FINALIZED_HASH"
printf 'CHECKPOINT_PATH=%s\n' "$CHECKPOINT_FILE"

exit 0
