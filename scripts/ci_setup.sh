#!/usr/bin/env bash

set -euo pipefail
IFS=$'\n\t'

SCRIPT_NAME=$(basename "$0")

usage() {
  cat <<'USAGE'
Usage: ci_setup.sh [options]

Enterprise continuous integration bootstrapper for the Synnergy Network
repository. The workflow aligns CLI automation, consensus builds, VM
regressions, wallet/node toolchains, and JavaScript UI bundles with a single
command that can run locally or inside CI runners.

Optional arguments:
  --project-root PATH       Repository root (default: auto-detected via git).
  --profile NAME            Execution profile: full, verify, lint, or build-only
                             (default: full).
  --toolchain LIST          Comma-separated toolchains to initialize. Supported
                             entries: go, node, docs (default: go,node).
  --skip-tests              Skip unit/integration tests even in profiles that
                             usually run them.
  --skip-build              Skip binary bundling steps.
  --ci-config PATH          Validate that the CI configuration file exists.
  --cache-dir PATH          Ensure cache directory is initialized.
  --lock-file PATH          Use flock-based locking to prevent parallel runs.
  --report PATH             Write JSON execution report for downstream
                             observability.
  --log-file PATH           Append structured logs to the provided file.
  --dry-run                 Print tasks without executing them.
  --no-fail-fast            Continue executing tasks even if one fails.
  --help, -h                Display this help message.

Examples:
  # Full pipeline bootstrap with reporting
  ci_setup.sh --project-root /srv/synnergy --report /var/log/ci/report.json

  # Developer sanity check without build/test steps
  ci_setup.sh --profile lint --dry-run
USAGE
}

LOG_FILE=""
REPORT_FILE=""
LOCK_FILE=""
DRY_RUN=0
FAIL_FAST=1
PROJECT_ROOT=""
PROFILE="full"
TOOLCHAIN="go,node"
SKIP_TESTS=0
SKIP_BUILD=0
CI_CONFIG=""
CACHE_DIR=""
OVERALL_FAILURE=0

log() {
  local level=$1
  local message=$2
  local timestamp
  timestamp=$(date -u +'%Y-%m-%dT%H:%M:%SZ')
  printf '%s [%s] %s\n' "$timestamp" "$level" "$message" >&2
  if [[ -n "$LOG_FILE" ]]; then
    printf '%s [%s] %s\n' "$timestamp" "$level" "$message" >>"$LOG_FILE"
  fi
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
  LOCK_FD=214
  eval "exec ${LOCK_FD}>\"$path\""
  if ! flock -n "$LOCK_FD"; then
    fatal "another ci_setup.sh invocation already holds the lock at $path"
  fi
}

require_command() {
  local cmd=$1
  local description=${2:-$1}
  if ! command -v "$cmd" >/dev/null 2>&1; then
    fatal "$description is required but not installed"
  fi
}

RESULT_ROWS=()
TASK_IDS=()
TASK_DESCS=()
TASK_CMDS=()

add_task() {
  local id=$1
  local description=$2
  local command=$3
  TASK_IDS+=("$id")
  TASK_DESCS+=("$description")
  TASK_CMDS+=("$command")
}

run_tasks() {
  local count=${#TASK_IDS[@]}
  local idx=0
  while [[ $idx -lt $count ]]; do
    local id=${TASK_IDS[$idx]}
    local description=${TASK_DESCS[$idx]}
    local command=${TASK_CMDS[$idx]}
    local start_ns end_ns duration status exit_code
    start_ns=$(python3 -c 'import time; print(time.time_ns())')
    if [[ $DRY_RUN -eq 1 ]]; then
      status="skipped"
      exit_code=0
      log "INFO" "dry-run: [$id] $command"
    else
      log "INFO" "executing [$id] $command"
      if eval "$command"; then
        status="success"
        exit_code=0
      else
        exit_code=$?
        status="failed"
        log "ERROR" "task $id failed with exit code $exit_code"
        OVERALL_FAILURE=1
        if [[ $FAIL_FAST -eq 1 ]]; then
          end_ns=$(python3 -c 'import time; print(time.time_ns())')
          duration=$(( (end_ns - start_ns) / 1000000 ))
          RESULT_ROWS+=("$id|$description|$status|$exit_code|$duration|$command")
          return
        fi
      fi
    fi
    end_ns=$(python3 -c 'import time; print(time.time_ns())')
    duration=$(( (end_ns - start_ns) / 1000000 ))
    RESULT_ROWS+=("$id|$description|$status|$exit_code|$duration|$command")
    idx=$((idx + 1))
  done
}

write_report() {
  local file=$1
  shift
  if [[ -z "$file" ]]; then
    return
  fi
  mkdir -p "$(dirname "$file")"
  if [[ $# -eq 0 ]]; then
    printf '{"tasks": []}\n' >"$file"
    return
  fi
  python3 - "$file" "$@" <<'PY'
import json
import sys

file_path = sys.argv[1]
rows = []
for raw in sys.argv[2:]:
    if not raw:
        continue
    parts = raw.split('|', 5)
    if len(parts) != 6:
        continue
    task = {
        "id": parts[0],
        "description": parts[1],
        "status": parts[2],
        "exit_code": int(parts[3]),
        "duration_ms": int(parts[4]),
        "command": parts[5],
    }
    rows.append(task)
with open(file_path, 'w', encoding='utf-8') as fh:
    json.dump({"tasks": rows}, fh, indent=2)
    fh.write('\n')
PY
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --project-root)
      PROJECT_ROOT=${2:-}
      shift 2
      ;;
    --profile)
      PROFILE=${2:-}
      shift 2
      ;;
    --toolchain)
      TOOLCHAIN=${2:-}
      shift 2
      ;;
    --skip-tests)
      SKIP_TESTS=1
      shift
      ;;
    --skip-build)
      SKIP_BUILD=1
      shift
      ;;
    --ci-config)
      CI_CONFIG=${2:-}
      shift 2
      ;;
    --cache-dir)
      CACHE_DIR=${2:-}
      shift 2
      ;;
    --lock-file)
      LOCK_FILE=${2:-}
      shift 2
      ;;
    --report)
      REPORT_FILE=${2:-}
      shift 2
      ;;
    --log-file)
      LOG_FILE=${2:-}
      shift 2
      ;;
    --dry-run)
      DRY_RUN=1
      shift
      ;;
    --no-fail-fast)
      FAIL_FAST=0
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

if [[ -z "$PROJECT_ROOT" ]]; then
  if git rev-parse --show-toplevel >/dev/null 2>&1; then
    PROJECT_ROOT=$(git rev-parse --show-toplevel)
  else
    PROJECT_ROOT=$(pwd)
  fi
fi

if [[ ! -d "$PROJECT_ROOT" ]]; then
  fatal "project root $PROJECT_ROOT does not exist"
fi

PROJECT_ROOT=$(cd "$PROJECT_ROOT" && pwd)

if [[ -n "$LOG_FILE" ]]; then
  mkdir -p "$(dirname "$LOG_FILE")"
  : >"$LOG_FILE"
fi

if [[ -n "$CACHE_DIR" ]]; then
  if [[ $DRY_RUN -eq 0 ]]; then
    mkdir -p "$CACHE_DIR"
  else
    log "INFO" "dry-run: would create cache directory at $CACHE_DIR"
  fi
fi

if [[ -n "$CI_CONFIG" && ! -f "$CI_CONFIG" ]]; then
  fatal "CI config $CI_CONFIG not found"
fi

IFS=',' read -r -a TOOLCHAIN_ENTRIES <<<"$TOOLCHAIN"
VALID_TOOLCHAINS=(go node docs)
for entry in "${TOOLCHAIN_ENTRIES[@]}"; do
  case "$entry" in
    go|node|docs)
      ;;
    "")
      ;;
    *)
      fatal "unsupported toolchain '$entry'"
      ;;
  esac
done

case "$PROFILE" in
  full)
    ;;
  verify)
    SKIP_BUILD=1
    ;;
  lint)
    SKIP_TESTS=1
    SKIP_BUILD=1
    ;;
  build-only)
    SKIP_TESTS=1
    ;;
  *)
    fatal "unknown profile '$PROFILE'"
    ;;
esac

acquire_lock "$LOCK_FILE"

if [[ $DRY_RUN -eq 0 ]]; then
  require_command bash "bash"
  if printf '%s\n' "${TOOLCHAIN_ENTRIES[@]}" | grep -q 'go'; then
    require_command go "Go toolchain"
  fi
  if printf '%s\n' "${TOOLCHAIN_ENTRIES[@]}" | grep -q 'node'; then
    require_command npm "Node/npm toolchain"
  fi
fi

add_go_tasks() {
  if printf '%s\n' "${TOOLCHAIN_ENTRIES[@]}" | grep -q '^go$'; then
    add_task "go-mod-tidy" "Go modules tidy" "cd \"$PROJECT_ROOT\" && go mod tidy"
    add_task "go-mod-verify" "Verify Go modules" "cd \"$PROJECT_ROOT\" && go mod verify"
    if [[ $SKIP_TESTS -eq 0 ]]; then
      add_task "go-test" "Run Go tests" "cd \"$PROJECT_ROOT\" && go test ./..."
      if [[ -f \"$PROJECT_ROOT/scripts/run_tests.sh\" ]]; then
        add_task "syn-tests" "Synnergy scripted tests" "cd \"$PROJECT_ROOT\" && bash scripts/run_tests.sh"
      fi
    fi
    if [[ $SKIP_BUILD -eq 0 ]]; then
      if [[ -f \"$PROJECT_ROOT/scripts/build_all.sh\" ]]; then
        add_task "syn-build" "Build Synnergy binaries" "cd \"$PROJECT_ROOT\" && bash scripts/build_all.sh"
      fi
    fi
    if [[ -f \"$PROJECT_ROOT/scripts/lint.sh\" && $PROFILE != build-only ]]; then
      add_task "syn-lint" "Run Synnergy linters" "cd \"$PROJECT_ROOT\" && bash scripts/lint.sh"
    fi
  fi
}

add_node_tasks() {
  if printf '%s\n' "${TOOLCHAIN_ENTRIES[@]}" | grep -q '^node$'; then
    local web_dir="$PROJECT_ROOT/web"
    if [[ -d "$web_dir" ]]; then
      add_task "node-install" "Install Node dependencies" "cd \"$web_dir\" && npm ci"
      if [[ $SKIP_BUILD -eq 0 ]]; then
        add_task "node-build" "Build web assets" "cd \"$web_dir\" && npm run build"
      fi
      if [[ $SKIP_TESTS -eq 0 ]]; then
        add_task "node-test" "Run web unit tests" "cd \"$web_dir\" && npm test -- --watch=false"
      fi
    else
      log "INFO" "web directory not present; skipping node toolchain"
    fi
  fi
}

add_docs_tasks() {
  if printf '%s\n' "${TOOLCHAIN_ENTRIES[@]}" | grep -q '^docs$'; then
    if [[ -f "$PROJECT_ROOT/mkdocs.yml" ]]; then
      add_task "docs-build" "Build documentation site" "cd \"$PROJECT_ROOT\" && mkdocs build"
    else
      log "INFO" "mkdocs.yml missing; skipping docs toolchain"
    fi
  fi
}

add_go_tasks
add_node_tasks
add_docs_tasks

run_tasks

write_report "$REPORT_FILE" "${RESULT_ROWS[@]}"

if [[ $OVERALL_FAILURE -eq 1 ]]; then
  fatal "CI setup encountered failures" 1
fi

log "INFO" "ci_setup.sh completed successfully"
exit 0
