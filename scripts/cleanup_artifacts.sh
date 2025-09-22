#!/usr/bin/env bash

set -euo pipefail
IFS=$'\n\t'

SCRIPT_NAME=$(basename "$0")

usage() {
  cat <<'USAGE'
Usage: cleanup_artifacts.sh [options]

Enterprise artifact cleanup utility that protects consensus-critical data while
reclaiming build caches, coverage dumps, CLI sandboxes, and UI bundles. The
cleanup integrates with CI/CD, operator workflows, and air-gapped rotations by
supporting retention windows, pattern filtering, metrics emission, and
idempotent dry-runs.

Optional arguments:
  --workspace PATH          Base directory (default: current working directory).
  --target NAME             Logical target to purge: build, cache, logs, or tmp.
                             Can be specified multiple times.
  --path PATH               Additional path to purge (relative to workspace by
                             default).
  --include PATTERN         Glob pattern to include when pruning directories.
  --exclude PATTERN         Glob pattern to exclude from pruning.
  --max-age HOURS           Only remove files older than the given age.
  --metrics-file PATH       Write JSON metrics describing the cleanup.
  --lock-file PATH          Flock-based mutex to avoid concurrent runs.
  --dry-run                 Report actions without deleting files.
  --force                   Allow deleting paths outside the workspace.
  --help, -h                Display this help message.

Examples:
  # Dry-run of build/cache cleanup for nightly job
  cleanup_artifacts.sh --workspace /srv/synnergy --target build --target cache --dry-run

  # Remove log archives older than 72 hours and emit metrics
  cleanup_artifacts.sh --workspace /srv/synnergy --target logs --max-age 72 \
    --metrics-file /var/log/synnergy/cleanup.json
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
  if [[ -n "${LOCK_FD:-}" ]]; then
    eval "exec ${LOCK_FD}>&-"
  fi
  if [[ $status -ne 0 ]]; then
    log "ERROR" "${SCRIPT_NAME} terminated with status $status"
  fi
}
trap cleanup EXIT
trap 'fatal "${SCRIPT_NAME} interrupted" 2' INT TERM

acquire_lock() {
  local path=$1
  if [[ -z "$path" ]]; then
    return
  fi
  mkdir -p "$(dirname "$path")"
  LOCK_FD=215
  eval "exec ${LOCK_FD}>\"$path\""
  if ! flock -n "$LOCK_FD"; then
    fatal "another cleanup_artifacts.sh invocation already holds the lock at $path"
  fi
}

WORKSPACE="$(pwd)"
METRICS_FILE=""
LOCK_FILE=""
DRY_RUN=0
FORCE=0
MAX_AGE_HOURS=-1
TARGETS=()
CUSTOM_PATHS=()
INCLUDE_PATTERNS=()
EXCLUDE_PATTERNS=()

while [[ $# -gt 0 ]]; do
  case "$1" in
    --workspace)
      WORKSPACE=${2:-}
      shift 2
      ;;
    --target)
      TARGETS+=("${2:-}")
      shift 2
      ;;
    --path)
      CUSTOM_PATHS+=("${2:-}")
      shift 2
      ;;
    --include)
      INCLUDE_PATTERNS+=("${2:-}")
      shift 2
      ;;
    --exclude)
      EXCLUDE_PATTERNS+=("${2:-}")
      shift 2
      ;;
    --max-age)
      MAX_AGE_HOURS=${2:-}
      shift 2
      ;;
    --metrics-file)
      METRICS_FILE=${2:-}
      shift 2
      ;;
    --lock-file)
      LOCK_FILE=${2:-}
      shift 2
      ;;
    --dry-run)
      DRY_RUN=1
      shift
      ;;
    --force)
      FORCE=1
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

if [[ ! -d "$WORKSPACE" ]]; then
  fatal "workspace $WORKSPACE does not exist"
fi

WORKSPACE=$(cd "$WORKSPACE" && pwd)

if [[ ${#TARGETS[@]} -eq 0 && ${#CUSTOM_PATHS[@]} -eq 0 ]]; then
  TARGETS+=(build)
fi

if [[ -n "$METRICS_FILE" ]]; then
  mkdir -p "$(dirname "$METRICS_FILE")"
fi

if [[ $MAX_AGE_HOURS != -1 && ! $MAX_AGE_HOURS =~ ^[0-9]+$ ]]; then
  fatal "--max-age must be a non-negative integer"
fi

acquire_lock "$LOCK_FILE"

collect_targets() {
  local target=$1
  case "$target" in
    build)
      printf '%s\n' "build" "dist" "out" "bin" "artifacts" "deploy"
      ;;
    cache)
      printf '%s\n' ".cache" "cache" "tmp/cache" "node_modules/.cache" "web/node_modules/.cache" "web/.next"
      ;;
    logs)
      printf '%s\n' "logs" "var/log" "*.log" "log/*.log" "logs/*.log"
      ;;
    tmp|temp)
      printf '%s\n' "tmp" "temp" "sandbox" "run/tmp"
      ;;
    *)
      fatal "unsupported target '$target'"
      ;;
  esac
}

CANDIDATES=()
for target in "${TARGETS[@]}"; do
  while IFS= read -r entry; do
    [[ -z "$entry" ]] && continue
    CANDIDATES+=("$entry")
  done < <(collect_targets "$target")
done

for path in "${CUSTOM_PATHS[@]}"; do
  CANDIDATES+=("$path")
done

if [[ ${#CANDIDATES[@]} -eq 0 ]]; then
  log "INFO" "no candidate paths resolved; exiting"
  exit 0
fi

python_args=("${METRICS_FILE}" "$WORKSPACE" "$MAX_AGE_HOURS" "$DRY_RUN" "$FORCE")
python_args+=("${CANDIDATES[@]}")
python_args+=("--")
python_args+=("${INCLUDE_PATTERNS[@]}")
python_args+=("--")
python_args+=("${EXCLUDE_PATTERNS[@]}")

python3 - "$SCRIPT_NAME" "${python_args[@]}" <<'PY'
import json
import os
import shutil
import sys
import time
from pathlib import Path
from typing import Iterable

SCRIPT_NAME = sys.argv[1]
metrics_path = sys.argv[2]
workspace = Path(sys.argv[3]).resolve()
max_age_hours = int(sys.argv[4])
dry_run = sys.argv[5] == "1"
force = sys.argv[6] == "1"
args = sys.argv[7:]

if "--" not in args:
    print("invalid argument layout", file=sys.stderr)
    sys.exit(2)
divider_one = args.index("--")
candidate_args = args[:divider_one]
rest = args[divider_one + 1:]
if "--" in rest:
    divider_two = rest.index("--")
    include_args = rest[:divider_two]
    exclude_args = rest[divider_two + 1:]
else:
    include_args = rest
    exclude_args = []

includes = [pattern for pattern in include_args if pattern]
excludes = [pattern for pattern in exclude_args if pattern]
max_age_seconds = -1 if max_age_hours < 0 else max_age_hours * 3600
now = time.time()
removed = []

def within_workspace(path: Path) -> bool:
    try:
        path.relative_to(workspace)
        return True
    except ValueError:
        return False

candidates: list[tuple[Path, str]] = []
for raw in candidate_args:
    if not raw:
        continue
    if any(ch in raw for ch in "*?["):
        for match in workspace.glob(raw):
            candidates.append((match.resolve(), raw))
    else:
        path = Path(raw)
        if not path.is_absolute():
            path = (workspace / raw).resolve()
        candidates.append((path, raw))

def should_keep(path: Path) -> bool:
    relative = str(path.relative_to(workspace)) if within_workspace(path) else str(path)
    if includes and not any(Path(relative).match(pattern) for pattern in includes):
        return True
    if excludes and any(Path(relative).match(pattern) for pattern in excludes):
        return True
    if max_age_seconds >= 0 and path.exists():
        age = now - path.stat().st_mtime
        if age < max_age_seconds:
            return True
    return False

for path, source in candidates:
    if not force and not within_workspace(path):
        print(f"Skipping {path} (outside workspace)")
        continue
    if not path.exists():
        continue
    if path.is_dir():
        if includes:
            for child in sorted(path.rglob("*")):
                if not child.exists():
                    continue
                if should_keep(child):
                    continue
                size = child.stat().st_size if child.is_file() else 0
                action = "Would remove" if dry_run else "Removing"
                print(f"{action} {child}")
                if not dry_run:
                    if child.is_file() or child.is_symlink():
                        child.unlink(missing_ok=True)
                    else:
                        shutil.rmtree(child, ignore_errors=True)
                removed.append({"path": str(child), "bytes": size, "type": "dir" if child.is_dir() else "file"})
        else:
            size = 0
            for child in path.rglob("*"):
                if child.is_file():
                    size += child.stat().st_size
            action = "Would remove" if dry_run else "Removing"
            print(f"{action} {path}")
            if not dry_run:
                shutil.rmtree(path, ignore_errors=True)
            removed.append({"path": str(path), "bytes": size, "type": "dir"})
    else:
        if should_keep(path):
            continue
        size = path.stat().st_size
        action = "Would remove" if dry_run else "Removing"
        print(f"{action} {path}")
        if not dry_run:
            path.unlink(missing_ok=True)
        removed.append({"path": str(path), "bytes": size, "type": "file"})

metrics = {
    "script": SCRIPT_NAME,
    "workspace": str(workspace),
    "dry_run": dry_run,
    "removed": removed,
    "count": len(removed),
    "bytes": sum(item["bytes"] for item in removed),
    "generated_at": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime(now)),
}
if metrics_path:
    with open(metrics_path, "w", encoding="utf-8") as fh:
        json.dump(metrics, fh, indent=2)
        fh.write("\n")
print(f"Summary: {metrics['count']} items, {metrics['bytes']} bytes", flush=True)
PY

status=$?
if [[ $status -ne 0 ]]; then
  exit $status
fi

log "INFO" "cleanup_artifacts.sh completed"
exit 0
