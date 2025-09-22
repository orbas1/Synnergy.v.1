#!/usr/bin/env bash
# Common helpers shared by Synnergy operational scripts.
# shellcheck disable=SC2034,SC2155

if [[ -z "${BASH_VERSION:-}" ]]; then
  echo "common.sh requires bash" >&2
  exit 1
fi

LIB_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SCRIPTS_ROOT="$(cd "$LIB_DIR/.." && pwd)"
PROJECT_ROOT="$(cd "$SCRIPTS_ROOT/.." && pwd)"
DEFAULT_LOG_DIR="${LOG_DIR:-$SCRIPTS_ROOT/logs}"
DEFAULT_TIMEOUT="${SYN_TIMEOUT:-120}"
DRY_RUN=false
LOG_FILE=""

mkdir -p "$DEFAULT_LOG_DIR"

log_timestamp() {
  date -u +"%Y-%m-%dT%H:%M:%SZ"
}

log() {
  local level="$1"
  shift
  local ts="$(log_timestamp)"
  local msg="$*"
  printf '%s [%s] %s\n' "$ts" "$level" "$msg"
  if [[ -n "$LOG_FILE" ]]; then
    printf '%s [%s] %s\n' "$ts" "$level" "$msg" >>"$LOG_FILE"
  fi
}

log_info() { log INFO "$@"; }
log_warn() { log WARN "$@"; }
log_error() { log ERROR "$@"; }

require_command() {
  local missing=()
  for cmd in "$@"; do
    if ! command -v "$cmd" &>/dev/null; then
      missing+=("$cmd")
    fi
  done
  if ((${#missing[@]})); then
    log_error "Missing required command(s): ${missing[*]}"
    return 1
  fi
  return 0
}

set_log_file() {
  local name="$1"
  LOG_FILE="$DEFAULT_LOG_DIR/${name%.sh}.log"
  : >"$LOG_FILE"
  log_info "Logging to $LOG_FILE"
}

parse_common_flags() {
  POSITIONAL_ARGS=()
  while [[ $# -gt 0 ]]; do
    case "$1" in
      --dry-run)
        DRY_RUN=true
        shift
        ;;
      --timeout)
        if [[ -z "${2:-}" ]]; then
          log_error "--timeout requires a value"
          return 1
        fi
        DEFAULT_TIMEOUT="$2"
        shift 2
        ;;
      --log-file)
        if [[ -z "${2:-}" ]]; then
          log_error "--log-file requires a value"
          return 1
        fi
        LOG_FILE="$2"
        mkdir -p "$(dirname "$LOG_FILE")"
        : >"$LOG_FILE"
        shift 2
        ;;
      --)
        shift
        while [[ $# -gt 0 ]]; do
          POSITIONAL_ARGS+=("$1")
          shift
        done
        break
        ;;
      *)
        POSITIONAL_ARGS+=("$1")
        shift
        ;;
    esac
  done
  return 0
}

with_timeout() {
  local desc="$1"
  shift
  local -a cmd=("$@")
  if [[ "$DRY_RUN" == true ]]; then
    log_info "[dry-run] $desc -> ${cmd[*]}"
    return 0
  fi
  local status=0
  if command -v timeout &>/dev/null; then
    timeout "$DEFAULT_TIMEOUT" "${cmd[@]}" 2>&1 | tee -a "$LOG_FILE"
    status=${PIPESTATUS[0]}
  else
    "${cmd[@]}" 2>&1 | tee -a "$LOG_FILE"
    status=${PIPESTATUS[0]}
  fi
  if ((status != 0)); then
    log_error "$desc failed with status $status"
  else
    log_info "$desc completed"
  fi
  return $status
}

retry() {
  local attempts="$1"
  shift
  local delay="$1"
  shift
  local desc="$1"
  shift
  local -a cmd=("$@")
  local try
  for ((try = 1; try <= attempts; try++)); do
    log_info "Attempt $try/$attempts: $desc"
    if with_timeout "$desc" "${cmd[@]}"; then
      return 0
    fi
    if ((try < attempts)); then
      log_warn "Retrying in ${delay}s"
      sleep "$delay"
    fi
  done
  return 1
}

synnergy_cli_path() {
  local candidate=""
  if [[ -n "${SYN_CLI_BIN:-}" ]]; then
    candidate="$SYN_CLI_BIN"
  elif [[ -n "${BIN_PATH:-}" ]]; then
    candidate="$BIN_PATH"
  fi
  if [[ -n "$candidate" ]]; then
    if [[ -x "$candidate" ]]; then
      printf '%s\n' "$candidate"
      return 0
    fi
    log_error "binary not found: $candidate"
    return 1
  fi
  local compiled="$PROJECT_ROOT/bin/synnergy"
  if [[ -x "$compiled" ]]; then
    printf '%s\n' "$compiled"
    return 0
  fi
  printf 'go run ./cmd/synnergy'
  return 0
}

synnergy_cli() {
  local bin
  local bin_output
  if ! bin_output="$(synnergy_cli_path)"; then
    printf '%s\n' "$bin_output" >&2
    return 1
  fi
  bin="$bin_output"
  local -a args=()
  if [[ "$bin" == go* ]]; then
    args=(go run ./cmd/synnergy "$@")
  else
    args=("$bin" "$@")
  fi
  retry 3 2 "synnergy ${*}" "${args[@]}"
}

load_env_file() {
  local file="$1"
  if [[ -f "$file" ]]; then
    # shellcheck disable=SC1090
    source "$file"
    log_info "Loaded environment from $file"
  fi
}

ensure_directory() {
  local dir="$1"
  if [[ ! -d "$dir" ]]; then
    mkdir -p "$dir"
    log_info "Created directory $dir"
  fi
}

