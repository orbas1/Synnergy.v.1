#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=../scripts/lib/common.sh
source "$SCRIPT_DIR/lib/common.sh"

usage() {
  cat <<'USAGE'
Usage: data_operations.sh [global options] <command> [command options]

Global options:
  --dry-run            Log actions without performing them
  --timeout <seconds>  Override command timeout (default 120)
  --log-file <path>    Write logs to a custom location

Commands:
  feed-sync     Apply a JSON manifest to a named feed and optionally emit a snapshot
  feed-export   Write the latest feed snapshot to a JSON file
  feed-prune    Remove a key from a feed to enforce retention rules
  resource-sync Import resources from a manifest and optionally prune stale entries
  resource-list Emit the registered resources to STDOUT or JSON

See the repository docs/guides for manifest schemas.
USAGE
}

ensure_requirements() {
  require_command jq || return 1
}

run_feed_sync() {
  local feed="" manifest="" snapshot=""
  while [[ $# -gt 0 ]]; do
    case "$1" in
      --feed)
        feed="$2"
        shift 2
        ;;
      --file)
        manifest="$2"
        shift 2
        ;;
      --snapshot-out)
        snapshot="$2"
        shift 2
        ;;
      *)
        log_error "Unknown feed-sync flag: $1"
        return 1
        ;;
    esac
  done
  if [[ -z "$feed" || -z "$manifest" ]]; then
    log_error "feed-sync requires --feed and --file"
    return 1
  fi
  if [[ ! -f "$manifest" ]]; then
    log_error "Manifest $manifest does not exist"
    return 1
  fi
  log_info "Applying feed manifest" "feed" "$feed" "manifest" "$manifest"
  synnergy_cli data feed apply --feed "$feed" --file "$manifest"
  if [[ -n "$snapshot" ]]; then
    ensure_directory "$(dirname "$snapshot")"
    log_info "Writing feed snapshot" "feed" "$feed" "output" "$snapshot"
    synnergy_cli data --json feed snapshot "$feed" >"$snapshot"
  fi
}

run_feed_export() {
  local feed="" output=""
  while [[ $# -gt 0 ]]; do
    case "$1" in
      --feed)
        feed="$2"
        shift 2
        ;;
      --out)
        output="$2"
        shift 2
        ;;
      *)
        log_error "Unknown feed-export flag: $1"
        return 1
        ;;
    esac
  done
  if [[ -z "$feed" || -z "$output" ]]; then
    log_error "feed-export requires --feed and --out"
    return 1
  fi
  ensure_directory "$(dirname "$output")"
  log_info "Exporting feed snapshot" "feed" "$feed" "output" "$output"
  synnergy_cli data --json feed snapshot "$feed" >"$output"
}

run_feed_prune() {
  local feed="" key=""
  while [[ $# -gt 0 ]]; do
    case "$1" in
      --feed)
        feed="$2"
        shift 2
        ;;
      --key)
        key="$2"
        shift 2
        ;;
      *)
        log_error "Unknown feed-prune flag: $1"
        return 1
        ;;
    esac
  done
  if [[ -z "$feed" || -z "$key" ]]; then
    log_error "feed-prune requires --feed and --key"
    return 1
  fi
  log_info "Deleting feed key" "feed" "$feed" "key" "$key"
  if ! synnergy_cli data feed delete "$feed" "$key"; then
    log_warn "Key $key not found in feed $feed"
  fi
}

run_resource_sync() {
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
        log_error "Unknown resource-sync flag: $1"
        return 1
        ;;
    esac
  done
  if [[ -z "$manifest" ]]; then
    log_error "resource-sync requires --manifest"
    return 1
  fi
  if [[ ! -f "$manifest" ]]; then
    log_error "Manifest $manifest does not exist"
    return 1
  fi
  log_info "Importing resource manifest" "manifest" "$manifest" "prune" "$prune"
  local args=(data resource import --manifest "$manifest")
  if [[ "$prune" == true ]]; then
    args+=(--prune)
  fi
  synnergy_cli "${args[@]}"
}

run_resource_list() {
  synnergy_cli data --json resource list
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

  local command="${POSITIONAL_ARGS[0]}"
  shift 0
  case "$command" in
    feed-sync)
      run_feed_sync "${POSITIONAL_ARGS[@]:1}"
      ;;
    feed-export)
      run_feed_export "${POSITIONAL_ARGS[@]:1}"
      ;;
    feed-prune)
      run_feed_prune "${POSITIONAL_ARGS[@]:1}"
      ;;
    resource-sync)
      run_resource_sync "${POSITIONAL_ARGS[@]:1}"
      ;;
    resource-list)
      run_resource_list
      ;;
    *)
      log_error "Unknown command: $command"
      usage
      return 1
      ;;
  esac
}

main "$@"
