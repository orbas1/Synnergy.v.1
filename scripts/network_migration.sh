#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib/common.sh
source "$SCRIPT_DIR/lib/common.sh"

usage() {
  cat <<'USAGE'
Perform a rolling restart to validate network migration workflows.

Usage: network_migration.sh [--checkpoint FILE] [--downtime SEC]
                            [--post-validate]

Options:
  --checkpoint FILE  File to store pre/post migration peer snapshots.
  --downtime SEC     Simulated downtime before restart (default: 3 seconds).
  --post-validate    Run post-migration harness to ensure consensus continuity.
  -h, --help         Display help message.
USAGE
}

CHECKPOINT="${PWD}/network-migration.json"
DOWNTIME=3
POST_VALIDATE=false

while [[ $# -gt 0 ]]; do
  case "$1" in
    --checkpoint)
      CHECKPOINT=$2
      shift 2
      ;;
    --downtime)
      DOWNTIME=$2
      shift 2
      ;;
    --post-validate)
      POST_VALIDATE=true
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      usage
      fail "unknown flag $1"
      ;;
  esac
done

require_commands jq

log_info "Capturing pre-migration peer state"
run_cli_json network start >/dev/null
pre_state=$(run_cli_json network peers)

log_info "Stopping network for migration"
run_cli_json network stop >/dev/null
sleep "$DOWNTIME"

log_info "Restarting network"
run_cli_json network start >/dev/null
post_state=$(run_cli_json network peers)

{
  echo "{"
  echo "  \"pre\": $pre_state,"
  echo "  \"post\": $post_state,"
  echo "  \"timestamp\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\""
  echo "}"
} | write_secure_file "$CHECKPOINT"

if [[ $POST_VALIDATE == true ]]; then
  log_info "Running post-migration validation cycle"
  run_cli_json node addtx "migration" "migration" "1" "1" "1" >/dev/null
  run_cli_json mining mine "migration-check" >/dev/null
  harness_out=$("$SCRIPT_DIR/network_harness.sh" --cycles 1 --topic migration --payload migration)
  log_info "Harness output: $harness_out"
fi

log_info "Network migration workflow complete"
