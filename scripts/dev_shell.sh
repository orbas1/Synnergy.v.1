#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck disable=SC1091
source "$SCRIPT_DIR/lib/common.sh"
set_log_file "$(basename "$0")"

usage() {
  cat <<'USAGE'
Open an interactive Synnergy development shell with all required environment
variables and helper aliases configured.

Usage: dev_shell.sh [options]

Options:
  --env-file FILE     Additional environment file to source
  --cmd COMMAND       Execute a single command within the prepared shell
  --check             Run dependency diagnostics and exit
  --dry-run           Show the resolved environment without launching
  --timeout SEC       Reserved for interface compatibility (ignored)
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

ENV_FILE=""
EXEC_CMD=""
CHECK_ONLY=false

while [[ $# -gt 0 ]]; do
  case "$1" in
    --env-file)
      ENV_FILE="${2:-}"
      shift 2
      ;;
    --cmd)
      EXEC_CMD="${2:-}"
      shift 2
      ;;
    --check)
      CHECK_ONLY=true
      shift
      ;;
    *)
      log_error "Unknown argument: $1"
      usage
      exit 1
      ;;
  esac
done

load_env_file "${SYN_ENV_FILE:-}"
[[ -n "$ENV_FILE" ]] && load_env_file "$ENV_FILE"

require_command go git >/dev/null

CLI_BIN="$(synnergy_cli_path)"
log_info "Resolved synnergy CLI path: $CLI_BIN"

if [[ "$CHECK_ONLY" == true ]]; then
  synnergy_cli --help >/dev/null
  log_info "CLI help command executed successfully"
  exit 0
fi

if [[ "$DRY_RUN" == true ]]; then
  log_info "[dry-run] would launch shell with SYN_CLI_BIN=$CLI_BIN"
  exit 0
fi

RC_FILE=$(mktemp)
cat >"$RC_FILE" <<'RC'
PS1='[synnergy-dev] \u@\h:\w$ '
alias syn="__CLI_BIN__"
export SYN_PROJECT_ROOT="__PROJECT_ROOT__"
export PATH="__PROJECT_ROOT__/bin:$PATH"
RC
CLI_BIN_VALUE="$CLI_BIN" PROJECT_ROOT_VALUE="$PROJECT_ROOT" python3 - "$RC_FILE" <<'PY'
import os
import sys

path = sys.argv[1]
content = open(path, "r", encoding="utf-8").read()
content = content.replace("__CLI_BIN__", os.environ["CLI_BIN_VALUE"])
content = content.replace("__PROJECT_ROOT__", os.environ["PROJECT_ROOT_VALUE"])
open(path, "w", encoding="utf-8").write(content)
PY

log_info "Launching interactive shell (rc: $RC_FILE)"
export SYN_CLI_BIN="$CLI_BIN"
export PROJECT_ROOT

if [[ -n "$EXEC_CMD" ]]; then
  bash --noprofile --rcfile "$RC_FILE" -lc "$EXEC_CMD"
else
  bash --noprofile --rcfile "$RC_FILE"
fi

rm -f "$RC_FILE"
