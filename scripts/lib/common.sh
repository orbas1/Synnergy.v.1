#!/usr/bin/env bash
# Common helpers for Synnergy operational scripts.
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
LOG_LEVEL=${LOG_LEVEL:-info}
LOG_FILE=${LOG_FILE:-}
JSON_OUTPUT=${JSON_OUTPUT:-true}
MAX_RETRIES=${MAX_RETRIES:-5}
RETRY_BACKOFF=${RETRY_BACKOFF:-2}

_log_ts() {
  date -u +"%Y-%m-%dT%H:%M:%SZ"
}

log_info() {
  local msg=$1
  local ts=$(_log_ts)
  echo "[$ts] [INFO] $msg"
  if [[ -n $LOG_FILE ]]; then
    echo "[$ts] [INFO] $msg" >>"$LOG_FILE"
  fi
}

log_warn() {
  local msg=$1
  local ts=$(_log_ts)
  echo "[$ts] [WARN] $msg" >&2
  if [[ -n $LOG_FILE ]]; then
    echo "[$ts] [WARN] $msg" >>"$LOG_FILE"
  fi
}

log_error() {
  local msg=$1
  local ts=$(_log_ts)
  echo "[$ts] [ERROR] $msg" >&2
  if [[ -n $LOG_FILE ]]; then
    echo "[$ts] [ERROR] $msg" >>"$LOG_FILE"
  fi
}

fail() {
  log_error "$1"
  exit 1
}

require_commands() {
  local missing=0
  for cmd in "$@"; do
    if ! command -v "$cmd" >/dev/null 2>&1; then
      log_error "required command '$cmd' not found"
      missing=1
    fi
  done
  if [[ $missing -eq 1 ]]; then
    fail "missing required tooling"
  fi
}

_detect_cli() {
  if [[ -n ${SYN_CLI_BIN:-} && -x ${SYN_CLI_BIN} ]]; then
    echo "$SYN_CLI_BIN"
    return
  fi
  if command -v synnergy >/dev/null 2>&1; then
    echo "synnergy"
    return
  fi
  echo "go"
}

_cli_args() {
  if [[ -n ${SYN_CONFIG:-} ]]; then
    echo "--config" "$SYN_CONFIG"
  fi
  if [[ ${LOG_LEVEL,,} == "debug" ]]; then
    echo "--log-level" "debug"
  fi
  if [[ ${JSON_OUTPUT} == true ]]; then
    echo "--json"
  fi
}

run_cli() {
  local bin
  bin=$(_detect_cli)
  local args
  mapfile -t args < <(_cli_args)
  local cmd=("$bin")
  if [[ $bin == "go" ]]; then
    cmd=("go" "run" "./cmd/synnergy")
  fi

  if [[ ${#args[@]} -gt 0 ]]; then
    cmd+=("${args[@]}")
  fi
  cmd+=("$@")

  if [[ $bin == "go" ]]; then
    (
      cd "$ROOT_DIR"
      log_info "running CLI: ${cmd[*]}"
      GO111MODULE=on "${cmd[@]}"
    )
  else
    log_info "running CLI: ${cmd[*]}"
    "${cmd[@]}"
  fi
}

run_cli_json() {
  local output
  if ! output=$(run_cli "$@" 2>&1); then
    log_error "cli command failed: $output"
    return 1
  fi
  echo "$output"
}

with_retry() {
  local attempt=1
  local max=${1:-$MAX_RETRIES}
  shift
  local delay=${RETRY_BACKOFF}
  until "$@"; do
    if [[ $attempt -ge $max ]]; then
      log_error "command failed after $attempt attempts"
      return 1
    fi
    log_warn "retrying command (attempt $((attempt + 1))/$max)"
    sleep $((delay * attempt))
    attempt=$((attempt + 1))
  done
}

write_secure_file() {
  local path=$1
  umask 0077
  cat >"$path"
  chmod 600 "$path"
}

json_extract() {
  require_commands jq
  local key=$1
  local input=$2
  echo "$input" | jq -r "$key"
}

