#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib/common.sh
source "$SCRIPT_DIR/lib/common.sh"

usage() {
  cat <<'USAGE'
Generate a regulatory compliance report from CLI telemetry.

Usage: regulatory_report.sh [--address ADDRESS] [--output FILE]
                            [--jurisdiction CODE]

Options:
  --address ADDRESS    Address to audit (default: random sample).
  --output FILE        File to write report JSON.
  --jurisdiction CODE  Filter regulations by jurisdiction (default: EU-103).
  -h, --help           Show help message.
USAGE
}

ADDRESS=""
OUTPUT=""
JURISDICTION="EU-103"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --address)
      ADDRESS=$2
      shift 2
      ;;
    --output)
      OUTPUT=$2
      shift 2
      ;;
    --jurisdiction)
      JURISDICTION=$2
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

if [[ -z $ADDRESS ]]; then
  tmp_wallet=$(run_cli_json wallet new)
  ADDRESS=$(json_extract '.address' "$tmp_wallet")
fi

regs=$(run_cli_json regulator list)
report=$(cat <<JSON
{
  "address": "$ADDRESS",
  "jurisdiction": "$JURISDICTION",
  "regulations": $regs,
  "audit": $(run_cli_json regnode audit "$ADDRESS"),
  "timestamp": "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
}
JSON
)

if [[ -n $OUTPUT ]]; then
  write_secure_file "$OUTPUT" <<JSON
$report
JSON
else
  echo "$report"
fi

log_info "Regulatory report generated"
