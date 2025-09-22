#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=../scripts/lib/common.sh
source "$SCRIPT_DIR/lib/common.sh"

usage() {
  cat <<'USAGE'
Usage: data_retention_policy_check.sh [global options] --manifest <file> [--fail-on-violation]

The manifest must be a JSON array of objects with the fields:
  key      - Resource key to inspect (required)
  max_age  - Allowed age (supports suffix s, m, h, d) (required)
  severity - Optional string flag for documentation

Example manifest entry:
  {"key":"asset","max_age":"720h","severity":"critical"}
USAGE
}

parse_duration() {
  local input="$1"
  if [[ $input =~ ^([0-9]+)([smhd])$ ]]; then
    local value="${BASH_REMATCH[1]}"
    local unit="${BASH_REMATCH[2]}"
    case "$unit" in
      s) echo "$value" ;;
      m) echo $((value * 60)) ;;
      h) echo $((value * 3600)) ;;
      d) echo $((value * 86400)) ;;
    esac
    return 0
  fi
  log_error "Unsupported duration format: $input"
  return 1
}

check_policy() {
  local manifest="" fail_on_violation=false
  while [[ $# -gt 0 ]]; do
    case "$1" in
      --manifest)
        manifest="$2"
        shift 2
        ;;
      --fail-on-violation)
        fail_on_violation=true
        shift
        ;;
      *)
        log_error "Unknown option: $1"
        return 1
        ;;
    esac
  done
  if [[ -z "$manifest" ]]; then
    log_error "--manifest is required"
    return 1
  fi
  if [[ ! -f "$manifest" ]]; then
    log_error "Manifest $manifest does not exist"
    return 1
  fi

  local exit_code=0
  local now_epoch
  now_epoch=$(date -u +%s)
  local entries
  mapfile -t entries < <(jq -c '.[]' "$manifest")
  if [[ ${#entries[@]} -eq 0 ]]; then
    log_warn "No retention entries found"
    return 0
  fi

  for entry in "${entries[@]}"; do
    local key max_age severity
    key=$(jq -r '.key // empty' <<<"$entry")
    max_age=$(jq -r '.max_age // empty' <<<"$entry")
    severity=$(jq -r '.severity // "info"' <<<"$entry")
    if [[ -z "$key" || -z "$max_age" ]]; then
      log_warn "Skipping entry missing key or max_age: $entry"
      continue
    fi
    local allowed
    if ! allowed=$(parse_duration "$max_age"); then
      exit_code=1
      continue
    fi
    log_info "Checking retention" "key" "$key" "max_age" "$max_age" "severity" "$severity"
    local info_json
    if ! info_json=$(synnergy_cli data --json resource info "$key" 2>/dev/null); then
      log_warn "Resource $key not found"
      if [[ "$fail_on_violation" == true ]]; then
        exit_code=1
      fi
      continue
    fi
    local updated
    updated=$(jq -r '.updated_at // empty' <<<"$info_json")
    if [[ -z "$updated" ]]; then
      log_warn "Resource $key missing updated_at field"
      continue
    fi
    local updated_epoch
    updated_epoch=$(date -u -d "$updated" +%s)
    local age=$((now_epoch - updated_epoch))
    if (( age > allowed )); then
      log_warn "Retention violation" "key" "$key" "age_seconds" "$age" "allowed" "$allowed"
      if [[ "$fail_on_violation" == true ]]; then
        exit_code=1
      fi
    else
      log_info "Resource within policy" "key" "$key" "age_seconds" "$age"
    fi
  done

  return $exit_code
}

main() {
  if [[ ${1:-} == "--help" || ${1:-} == "-h" ]]; then
    usage
    exit 0
  fi

  parse_common_flags "$@" || exit 1
  set_log_file "$(basename "$0")"
  require_command jq date

  if check_policy "${POSITIONAL_ARGS[@]}"; then
    exit 0
  fi
  exit 1
}

main "$@"
