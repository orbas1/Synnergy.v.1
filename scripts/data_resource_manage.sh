#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=../scripts/lib/common.sh
source "$SCRIPT_DIR/lib/common.sh"

usage() {
  cat <<'USAGE'
Usage: data_resource_manage.sh [global options] <command> [command options]

Commands:
  sync        Import resources from a manifest and optionally prune everything else
  audit       Print resource metadata for keys defined in a manifest
  usage       Report overall byte consumption for the resource catalog

Global options:
  --dry-run            Log actions without performing them
  --timeout <seconds>  Override command timeout (default 120)
  --log-file <path>    Write logs to a custom location
USAGE
}

ensure_requirements() {
  require_command jq || return 1
}

run_sync() {
  local manifest="" prune=false
  while [[ $# -gt 0 ]]; do
    case "$1" in
      --manifest)
        manifest="$2"
        shift 2
        ;;
      --prune)
        prune=true
        shift
        ;;
      *)
        log_error "Unknown sync flag: $1"
        return 1
        ;;
    esac
  done
  if [[ -z "$manifest" ]]; then
    log_error "sync requires --manifest"
    return 1
  fi
  if [[ ! -f "$manifest" ]]; then
    log_error "Manifest $manifest does not exist"
    return 1
  fi
  log_info "Synchronising resource manifest" "manifest" "$manifest" "prune" "$prune"
  local args=(data resource import --manifest "$manifest")
  if [[ "$prune" == true ]]; then
    args+=(--prune)
  fi
  synnergy_cli "${args[@]}"
}

run_audit() {
  local manifest=""
  while [[ $# -gt 0 ]]; do
    case "$1" in
      --manifest)
        manifest="$2"
        shift 2
        ;;
      *)
        log_error "Unknown audit flag: $1"
        return 1
        ;;
    esac
  done
  if [[ -z "$manifest" ]]; then
    log_error "audit requires --manifest"
    return 1
  fi
  if [[ ! -f "$manifest" ]]; then
    log_error "Manifest $manifest does not exist"
    return 1
  fi
  log_info "Auditing resources" "manifest" "$manifest"
  local keys
  mapfile -t keys < <(jq -r '.[].key | select(. != null)' "$manifest")
  if [[ ${#keys[@]} -eq 0 ]]; then
    log_warn "No keys found in manifest"
    return 0
  fi
  for key in "${keys[@]}"; do
    log_info "Fetching metadata" "key" "$key"
    synnergy_cli data --json resource info "$key"
  done
}

run_usage() {
  synnergy_cli data resource usage
}

main() {
  if [[ ${1:-} == "--help" || ${1:-} == "-h" ]]; then
    usage
    exit 0
  fi

  parse_common_flags "$@" || exit 1
  ensure_requirements || exit 1
  set_log_file "$(basename "$0")"

  if [[ ${#POSITIONAL_ARGS[@]} -eq 0 ]]; then
    usage
    exit 1
  fi

  case "${POSITIONAL_ARGS[0]}" in
    sync)
      run_sync "${POSITIONAL_ARGS[@]:1}"
      ;;
    audit)
      run_audit "${POSITIONAL_ARGS[@]:1}"
      ;;
    usage)
      run_usage
      ;;
    *)
      log_error "Unknown command: ${POSITIONAL_ARGS[0]}"
      usage
      return 1
      ;;
  esac
}

main "$@"
