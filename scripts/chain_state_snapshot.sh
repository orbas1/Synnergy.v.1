#!/usr/bin/env bash

set -euo pipefail
IFS=$'\n\t'

SCRIPT_NAME=$(basename "$0")

usage() {
  cat <<'USAGE'
Usage: chain_state_snapshot.sh --network <name> --output-dir <path> [options]

Capture a tamper-evident snapshot of the Synnergy Network chain state for
regulatory, audit, and disaster-recovery workflows. The snapshot integrates
with the Synnergy CLI or recorded telemetry, emits metrics, and enforces
retention/locking policies so that infrastructure automation can safely invoke
it in parallel with consensus and VM operations.

Required arguments:
  --network NAME            Logical network name (e.g., mainnet, stage98).
  --output-dir PATH         Directory where snapshots are persisted.

Optional arguments:
  --status-file PATH        Pre-fetched status JSON (used for dry-runs/tests).
  --cli PATH                Explicit Synnergy CLI binary. Defaults to $SYN_CLI,
                             then synnergy, then synnergy-cli.
  --endpoint URL            Override RPC endpoint for CLI queries.
  --profile NAME            CLI profile to use (default: production).
  --format TYPE             Snapshot format: json or archive (default: json).
  --snapshot-prefix PREFIX  File prefix for generated snapshots (default:
                             snapshot-).
  --retain COUNT            Retain the newest COUNT snapshots (default: 5).
  --tag VALUE               Include an additional descriptor in the snapshot
                             metadata (e.g., nightly, compliance).
  --timestamp VALUE         Override the snapshot timestamp. Primarily used for
                             testing; must follow YYYYMMDDTHHMMSSZ format.
  --lock-file PATH          Path to a flock-based mutex to prevent concurrent
                             executions.
  --metrics-file PATH       Emit JSON metrics describing the capture.
  --dry-run                 Skip CLI execution and filesystem mutations. Must
                             be combined with --status-file.
  --no-retention            Disable retention pruning for this invocation.
  --help, -h                Display this help message.

Examples:
  # Capture a snapshot from live CLI telemetry and keep the last 10 artifacts
  chain_state_snapshot.sh --network mainnet --output-dir /var/lib/synnergy/snaps \
    --retain 10 --metrics-file /var/lib/synnergy/metrics/snapshot.json

  # Idempotent CI validation using recorded status
  chain_state_snapshot.sh --network stage98 --output-dir ./tmp/snaps \
    --status-file ./fixtures/status.json --dry-run
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

resolve_cli() {
  local candidate=$1
  if [[ -n "$candidate" ]]; then
    if [[ -x "$candidate" ]]; then
      echo "$candidate"
      return 0
    fi
    if command -v "$candidate" >/dev/null 2>&1; then
      command -v "$candidate"
      return 0
    fi
    fatal "CLI binary '$candidate' not found"
  fi

  if [[ -n "${SYN_CLI:-}" ]]; then
    candidate=$SYN_CLI
    if [[ -x "$candidate" ]]; then
      echo "$candidate"
      return 0
    fi
    if command -v "$candidate" >/dev/null 2>&1; then
      command -v "$candidate"
      return 0
    fi
  fi

  if command -v synnergy >/dev/null 2>&1; then
    command -v synnergy
    return 0
  fi
  if command -v synnergy-cli >/dev/null 2>&1; then
    command -v synnergy-cli
    return 0
  fi

  fatal "CLI binary not found. Provide --cli or set SYN_CLI"
}

acquire_lock() {
  local path=$1
  if [[ -z "$path" ]]; then
    return
  fi
  mkdir -p "$(dirname "$path")"
  LOCK_FD=213
  eval "exec ${LOCK_FD}>\"$path\""
  if ! flock -n "$LOCK_FD"; then
    fatal "another snapshot capture already holds the lock at $path"
  fi
}

validate_timestamp() {
  local value=$1
  if [[ ! $value =~ ^[0-9]{8}T[0-9]{6}Z$ ]]; then
    fatal "invalid --timestamp '$value'; expected YYYYMMDDTHHMMSSZ format"
  fi
}

parse_status() {
  local file=$1
  local prefix=$2
  local stdout
  local stderr
  stderr=$(mktemp)
  if ! stdout=$(python3 - "$file" "$prefix" <<'PY' 2>"$stderr"
import json
import sys
import shlex

file_path = sys.argv[1]
prefix = sys.argv[2]
required = ["network", "height", "hash", "finalized_height", "finalized_hash"]
try:
    with open(file_path, "r", encoding="utf-8") as fh:
        data = json.load(fh)
except Exception as exc:  # pragma: no cover
    print(f"unable to read {file_path}: {exc}", file=sys.stderr)
    sys.exit(2)
missing = [key for key in required if key not in data]
if missing:
    print(f"{file_path} missing required keys: {', '.join(missing)}", file=sys.stderr)
    sys.exit(3)
if data["height"] < 0 or data["finalized_height"] < 0:
    print("heights must be non-negative", file=sys.stderr)
    sys.exit(4)
if data["finalized_height"] > data["height"]:
    print("finalized height cannot exceed head height", file=sys.stderr)
    sys.exit(5)
for key in required:
    value = data[key]
    print(f"{prefix}{key.upper()}={shlex.quote(str(value))}")
optional = ["timestamp", "finality_threshold", "finalized_by", "validators"]
for key in optional:
    if key in data:
        value = data[key]
        if isinstance(value, (dict, list)):
            value = json.dumps(value, separators=(",", ":"))
        print(f"{prefix}{key.upper()}={shlex.quote(str(value))}")
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

fetch_status() {
  local destination=$1
  local cli=$2
  local network=$3
  local endpoint=$4
  local profile=$5
  local tmp
  tmp=$(mktemp)
  local cmd=("$cli" chain status --network "$network" --output json)
  if [[ -n "$endpoint" ]]; then
    cmd+=(--endpoint "$endpoint")
  fi
  if [[ -n "$profile" ]]; then
    cmd+=(--profile "$profile")
  fi
  if ! "${cmd[@]}" >"$tmp" 2>"${tmp}.err"; then
    local message
    message=$(<"${tmp}.err")
    rm -f "$tmp" "${tmp}.err"
    fatal "CLI status query failed: ${message:-unknown error}"
  fi
  mv "$tmp" "$destination"
  rm -f "${tmp}.err"
}

write_snapshot() {
  local file=$1
  local network=$2
  local status_file=$3
  local source=$4
  local tag=$5
  local timestamp=$6
  python3 - "$file" "$network" "$status_file" "$source" "$tag" "$timestamp" <<'PY'
import json
import sys
from pathlib import Path

file_path, network, status_file, source, tag, timestamp = sys.argv[1:7]
with open(status_file, "r", encoding="utf-8") as fh:
    status = json.load(fh)
status.setdefault("timestamp", timestamp)
snapshot = {
    "network": network,
    "captured_at": timestamp,
    "source": source,
    "height": status["height"],
    "hash": status["hash"],
    "finalized_height": status["finalized_height"],
    "finalized_hash": status["finalized_hash"],
    "metadata": {
        "input_timestamp": status.get("timestamp", timestamp),
        "finality_threshold": status.get("finality_threshold"),
        "finalized_by": status.get("finalized_by"),
        "validators": status.get("validators"),
        "tag": tag or None,
    },
}
Path(file_path).write_text(json.dumps(snapshot, indent=2, sort_keys=True) + "\n", encoding="utf-8")
PY
}

prune_retention() {
  local directory=$1
  local prefix=$2
  local keep=$3
  local dry_run=$4
  local deleted
  deleted=$(python3 - "$directory" "$prefix" "$keep" "$dry_run" <<'PY'
import json
import sys
from pathlib import Path

directory, prefix, keep_str, dry_run_flag = sys.argv[1:5]
keep = int(keep_str)
if keep <= 0:
    print(json.dumps({"deleted": []}))
    sys.exit(0)
root = Path(directory)
patterns = [f"{prefix}*.json", f"{prefix}*.tar.gz"]
files = []
for pattern in patterns:
    files.extend(sorted(root.glob(pattern), key=lambda p: p.stat().st_mtime, reverse=True))
seen = set()
ordered = []
for path in files:
    if path in seen:
        continue
    seen.add(path)
    ordered.append(path)
if len(ordered) <= keep:
    print(json.dumps({"deleted": []}))
    sys.exit(0)
to_delete = ordered[keep:]
if dry_run_flag == "1":
    deleted_paths = [str(p) for p in to_delete]
else:
    deleted_paths = []
    for path in to_delete:
        try:
            path.unlink()
            deleted_paths.append(str(path))
        except FileNotFoundError:
            pass
print(json.dumps({"deleted": deleted_paths}))
PY
  )
  echo "$deleted"
}

write_metrics() {
  local file=$1
  local network=$2
  local path=$3
  local height=$4
  local finalized_height=$5
  local start_ns=$6
  local deleted_json=$7
  local timestamp
  timestamp=$(date -u +'%Y-%m-%dT%H:%M:%SZ')
  python3 - "$file" "$network" "$path" "$height" "$finalized_height" "$start_ns" "$timestamp" "$deleted_json" <<'PY'
import json
import sys
from pathlib import Path
import time

(
    file_path,
    network,
    snapshot_path,
    height,
    finalized_height,
    start_ns,
    emitted_at,
    deleted_json,
) = sys.argv[1:9]
height = int(height)
finalized_height = int(finalized_height)
start_ns = int(start_ns)
summary = json.loads(deleted_json)
deleted = summary.get("deleted", [])
metrics = {
    "network": network,
    "snapshot_path": snapshot_path,
    "height": height,
    "finalized_height": finalized_height,
    "emitted_at": emitted_at,
    "duration_ms": max(0, (int(time.time_ns()) - start_ns) // 1_000_000),
    "deleted": deleted,
}
Path(file_path).write_text(json.dumps(metrics, indent=2) + "\n", encoding="utf-8")
PY
}

NETWORK=""
OUTPUT_DIR=""
STATUS_FILE=""
CLI=""
ENDPOINT=""
PROFILE="production"
FORMAT="json"
RETAIN=5
SNAPSHOT_PREFIX="snapshot-"
TAG=""
TIMESTAMP_OVERRIDE=""
LOCK_FILE=""
METRICS_FILE=""
DRY_RUN=0
DISABLE_RETENTION=0

while [[ $# -gt 0 ]]; do
  case "$1" in
    --network)
      NETWORK=${2:-}
      shift 2
      ;;
    --output-dir)
      OUTPUT_DIR=${2:-}
      shift 2
      ;;
    --status-file)
      STATUS_FILE=${2:-}
      shift 2
      ;;
    --cli)
      CLI=${2:-}
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
    --format)
      FORMAT=${2:-}
      shift 2
      ;;
    --snapshot-prefix)
      SNAPSHOT_PREFIX=${2:-}
      shift 2
      ;;
    --retain)
      RETAIN=${2:-}
      shift 2
      ;;
    --tag)
      TAG=${2:-}
      shift 2
      ;;
    --timestamp)
      TIMESTAMP_OVERRIDE=${2:-}
      shift 2
      ;;
    --lock-file)
      LOCK_FILE=${2:-}
      shift 2
      ;;
    --metrics-file)
      METRICS_FILE=${2:-}
      shift 2
      ;;
    --dry-run)
      DRY_RUN=1
      shift
      ;;
    --no-retention)
      DISABLE_RETENTION=1
      shift
      ;;
    --help|-h)
      usage
      exit 0
      ;;
    *)
      log "ERROR" "unknown argument: $1"
      usage
      exit 1
      ;;
  esac
done

if [[ -z "$NETWORK" || -z "$OUTPUT_DIR" ]]; then
  usage
  fatal "--network and --output-dir are required"
fi

if [[ $FORMAT != "json" && $FORMAT != "archive" ]]; then
  fatal "--format must be json or archive"
fi

if [[ -n "$TIMESTAMP_OVERRIDE" ]]; then
  validate_timestamp "$TIMESTAMP_OVERRIDE"
fi

if ! [[ $RETAIN =~ ^[0-9]+$ ]]; then
  fatal "--retain must be a non-negative integer"
fi

if [[ $DRY_RUN -eq 1 && -z "$STATUS_FILE" ]]; then
  fatal "--dry-run requires --status-file"
fi

acquire_lock "$LOCK_FILE"

START_NS=$(python3 -c 'import time; print(time.time_ns())')

if [[ -n "$STATUS_FILE" ]]; then
  if [[ ! -f "$STATUS_FILE" ]]; then
    fatal "status file $STATUS_FILE does not exist"
  fi
  parse_status "$STATUS_FILE" "STATUS_"
  SOURCE="file:${STATUS_FILE}"
else
  CLI=$(resolve_cli "$CLI")
  TMP_STATUS_FILE=$(mktemp)
  if [[ $DRY_RUN -eq 1 ]]; then
    fatal "dry-run requested but no status-file provided"
  fi
  fetch_status "$TMP_STATUS_FILE" "$CLI" "$NETWORK" "$ENDPOINT" "$PROFILE"
  parse_status "$TMP_STATUS_FILE" "STATUS_"
  SOURCE="cli:${CLI}"
fi

if [[ "${STATUS_NETWORK}" != "$NETWORK" ]]; then
  fatal "status network '${STATUS_NETWORK}' does not match requested network '$NETWORK'"
fi

TIMESTAMP=${TIMESTAMP_OVERRIDE:-$(date -u +'%Y%m%dT%H%M%SZ')}

if [[ $DRY_RUN -eq 0 ]]; then
  mkdir -p "$OUTPUT_DIR"
fi

SNAPSHOT_BASENAME="${SNAPSHOT_PREFIX}${NETWORK}-${TIMESTAMP}"
SNAPSHOT_PATH_JSON="${OUTPUT_DIR%/}/${SNAPSHOT_BASENAME}.json"

if [[ $DRY_RUN -eq 0 ]]; then
  write_snapshot "$SNAPSHOT_PATH_JSON" "$NETWORK" "${STATUS_FILE:-$TMP_STATUS_FILE}" "$SOURCE" "$TAG" "$TIMESTAMP"
  CREATED_PATH="$SNAPSHOT_PATH_JSON"
  if [[ $FORMAT == "archive" ]]; then
    local_archive="${SNAPSHOT_PATH_JSON%.json}.tar.gz"
    tar -czf "$local_archive" -C "$(dirname "$SNAPSHOT_PATH_JSON")" "$(basename "$SNAPSHOT_PATH_JSON")"
    rm -f "$SNAPSHOT_PATH_JSON"
    CREATED_PATH="$local_archive"
  fi
else
  CREATED_PATH="$SNAPSHOT_PATH_JSON"
  log "INFO" "dry-run: would write snapshot to $CREATED_PATH"
fi

DELETED_INFO='{"deleted": []}'
if [[ $DISABLE_RETENTION -eq 0 && $RETAIN -gt 0 ]]; then
  if [[ $DRY_RUN -eq 0 ]]; then
    DELETED_INFO=$(prune_retention "$OUTPUT_DIR" "${SNAPSHOT_PREFIX}${NETWORK}-" "$RETAIN" 0)
  else
    DELETED_INFO=$(prune_retention "$OUTPUT_DIR" "${SNAPSHOT_PREFIX}${NETWORK}-" "$RETAIN" 1)
  fi
fi

if [[ -n "$METRICS_FILE" ]]; then
  write_metrics "$METRICS_FILE" "$NETWORK" "$CREATED_PATH" "${STATUS_HEIGHT}" "${STATUS_FINALIZED_HEIGHT}" "$START_NS" "$DELETED_INFO"
fi

log "INFO" "snapshot captured for $NETWORK at height ${STATUS_HEIGHT} -> $CREATED_PATH"

exit 0
