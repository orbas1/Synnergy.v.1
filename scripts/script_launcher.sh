#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck disable=SC1091
source "$SCRIPT_DIR/lib/common.sh"
# shellcheck disable=SC1091
source "$SCRIPT_DIR/lib/workflows.sh"

print_global_usage() {
  cat <<'USAGE'
Usage: script_launcher.sh --target SCRIPT [options]
       SCRIPT.sh [options]

Dispatch Synnergy workflow scripts through a unified implementation that
supports dry-run planning, parameter binding and manifest generation.

Options:
  --target SCRIPT        Run the named workflow (when invoking script_launcher.sh)
  --list                 List workflows managed by the launcher
  -h, --help             Show help for the selected workflow
  --plan                 Generate a plan without performing side effects
  --output PATH          Write workflow manifest to PATH
  --env-file FILE        Source FILE before executing the workflow
  --note TEXT            Attach TEXT to the workflow manifest (repeatable)
  --set key=value        Provide a workflow parameter (repeatable)
  --param KEY VALUE      Provide a workflow parameter (repeatable)

Common flags from lib/common.sh:
  --dry-run              Alias for --plan; no commands are executed
  --timeout SEC          Override the per-command timeout (default: 120)
  --log-file PATH        Append logs to PATH instead of scripts/logs/

Example:
  ./consensus_start.sh --plan --set profile=balanced
  ./script_launcher.sh --target wallet_transfer --set from=wallet-1 --set to=wallet-2 --set amount=25
USAGE
}

list_workflows() {
  local file
  while IFS= read -r file; do
    local name
    name="$(basename "$file")"
    if [[ "$name" == "script_launcher.sh" ]]; then
      continue
    fi
    printf '%s\n' "${name%.sh}"
  done < <(grep -l 'script_launcher.sh' "$SCRIPT_DIR"/*.sh | sort)
}

print_workflow_usage() {
  local script_label="$1"
  local base="$2"
  determine_workflow_context "$base"
  describe_workflow "$base"
  cat <<USAGE
Usage: $script_label [options] [common flags]

$WF_DESCRIPTION

Options:
  --plan                 Generate a plan without side effects
  --output PATH          Write workflow manifest to PATH
  --env-file FILE        Load environment variables from FILE before running
  --note TEXT            Append TEXT to the manifest (repeatable)
  --set key=value        Provide a workflow parameter (repeatable)
  --param KEY VALUE      Provide a workflow parameter (repeatable)
  -h, --help             Show this help text

Common flags from lib/common.sh:
  --dry-run              Alias for --plan
  --timeout SEC          Override the per-command timeout (default: 120)
  --log-file PATH        Append logs to PATH instead of scripts/logs/
USAGE
  if [[ -n "$WF_PARAM_HELP" ]]; then
    printf '\nParameter hints:\n%s\n' "$WF_PARAM_HELP"
  fi
}

main() {
  local invoked="$(basename "$0")"
  local target=""
  local args=()

  if [[ $# -eq 0 && "$invoked" == "script_launcher.sh" ]]; then
    print_global_usage
    exit 1
  fi

  if [[ "${1:-}" == "--list" ]]; then
    list_workflows
    exit 0
  fi

  if [[ "${1:-}" == "--target" ]]; then
    if [[ $# -lt 2 ]]; then
      log_error "--target requires a value"
      exit 1
    fi
    target="$2"
    shift 2
  elif [[ "$invoked" == "script_launcher.sh" ]]; then
    target="${1:-}"
    shift $(( $# > 0 ? 1 : 0 ))
  else
    target="$invoked"
  fi

  if [[ -z "$target" ]]; then
    print_global_usage
    exit 1
  fi

  target="${target##*/}"
  local base="${target%.sh}"
  local script_label="${base}.sh"

  parse_common_flags "$@"
  set -- "${POSITIONAL_ARGS[@]:-}"

  WORKFLOW_PLAN_MODE=false
  WORKFLOW_OUTPUT="${WORKFLOW_OUTPUT:-}"
  WORKFLOW_NOTES=()
  declare -gA WORKFLOW_PARAMS=()

  local env_file=""
  while [[ $# -gt 0 ]]; do
    case "$1" in
      --help|-h)
        print_workflow_usage "$script_label" "$base"
        exit 0
        ;;
      --plan)
        WORKFLOW_PLAN_MODE=true
        shift
        ;;
      --output)
        if [[ -z "${2:-}" ]]; then
          log_error "--output requires a path"
          exit 1
        fi
        WORKFLOW_OUTPUT="$2"
        shift 2
        ;;
      --env-file)
        if [[ -z "${2:-}" ]]; then
          log_error "--env-file requires a path"
          exit 1
        fi
        env_file="$2"
        shift 2
        ;;
      --note)
        if [[ -z "${2:-}" ]]; then
          log_error "--note requires text"
          exit 1
        fi
        WORKFLOW_NOTES+=("$2")
        shift 2
        ;;
      --set)
        if [[ -z "${2:-}" || "${2}" != *=* ]]; then
          log_error "--set expects key=value"
          exit 1
        fi
        WORKFLOW_PARAMS["${2%%=*}"]="${2#*=}"
        shift 2
        ;;
      --param)
        if [[ $# -lt 3 ]]; then
          log_error "--param expects KEY VALUE"
          exit 1
        fi
        WORKFLOW_PARAMS["$2"]="$3"
        shift 3
        ;;
      --)
        shift
        while [[ $# -gt 0 ]]; do
          WORKFLOW_NOTES+=("$1")
          shift
        done
        ;;
      *=*)
        WORKFLOW_PARAMS["${1%%=*}"]="${1#*=}"
        shift
        ;;
      *)
        WORKFLOW_NOTES+=("$1")
        shift
        ;;
    esac
  done

  if [[ -n "$env_file" ]]; then
    load_env_file "$env_file"
  fi

  if [[ -z "${LOG_FILE:-}" ]]; then
    set_log_file "$script_label"
  fi

  if [[ "$WORKFLOW_PLAN_MODE" == true ]]; then
    DRY_RUN=true
  fi

  local status=0
  if ! run_workflow_engine "$base"; then
    status=$?
  fi

  if (( status == 0 )); then
    log_info "Workflow $base completed"
  else
    log_error "Workflow $base finished with status $status"
  fi
  return $status
}

main "$@"
