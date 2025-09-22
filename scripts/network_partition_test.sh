#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib/common.sh
source "$SCRIPT_DIR/lib/common.sh"

usage() {
  cat <<'USAGE'
Simulate a network partition and ensure recovery workflows succeed.

Usage: network_partition_test.sh [--duration SEC] [--report FILE]
                                 [--topic TOPIC]

Options:
  --duration SEC   Duration of simulated partition before reconnection (default: 4).
  --report FILE    File to write detailed JSON report.
  --topic TOPIC    Topic used for pre/post partition broadcasts (default: partition).
  -h, --help       Show help message.
USAGE
}

DURATION=4
REPORT=""
TOPIC="partition"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --duration)
      DURATION=$2
      shift 2
      ;;
    --report)
      REPORT=$2
      shift 2
      ;;
    --topic)
      TOPIC=$2
      shift 2
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

run_cli_json network start >/dev/null
pre_peers=$(run_cli_json network peers)
run_cli_json network broadcast "$TOPIC" "pre-partition" >/dev/null

log_info "Simulating partition"
run_cli_json network stop >/dev/null
sleep "$DURATION"

log_info "Restoring network"
run_cli_json network start >/dev/null
run_cli_json network broadcast "$TOPIC" "post-partition" >/dev/null
post_peers=$(run_cli_json network peers)

consensus_block=$(run_cli_json mining mine "partition-check")

if [[ -n $REPORT ]]; then
  {
    echo "{"
    echo "  \"prePeers\": $pre_peers,"
    echo "  \"postPeers\": $post_peers,"
    echo "  \"consensus\": $consensus_block,"
    echo "  \"duration\": $DURATION,"
    echo "  \"timestamp\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\""
    echo "}"
  } | write_secure_file "$REPORT"
fi

log_info "Network partition simulation complete"
