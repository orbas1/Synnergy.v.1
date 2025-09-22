#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib/common.sh
source "$SCRIPT_DIR/lib/common.sh"

usage() {
  cat <<'USAGE'
Configure the regulatory manager and regulatory node workflows.

Usage: regulatory_node_setup.sh [--jurisdiction CODE] [--max AMOUNT]
                                [--wallet FILE] [--password PASS]

Options:
  --jurisdiction CODE  Jurisdiction identifier (default: EU-103).
  --max AMOUNT         Maximum transaction amount before manual review (default: 5000).
  --wallet FILE        Wallet file for signing approvals.
  --password PASS      Password to decrypt wallet file.
  -h, --help           Display help message.
USAGE
}

JURISDICTION="EU-103"
MAX_AMOUNT=5000
WALLET_FILE=""
PASSWORD=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --jurisdiction)
      JURISDICTION=$2
      shift 2
      ;;
    --max)
      MAX_AMOUNT=$2
      shift 2
      ;;
    --wallet)
      WALLET_FILE=$2
      shift 2
      ;;
    --password)
      PASSWORD=$2
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

log_info "Registering regulation"
run_cli_json regulator add "reg-$JURISDICTION" "$JURISDICTION" "Stage103 limits" "$MAX_AMOUNT" >/dev/null

if [[ -z $WALLET_FILE ]]; then
  temp_wallet=$(run_cli_json wallet new --out "${PWD}/wallets/regnode.json" --password "synnergy")
  WALLET_FILE="${PWD}/wallets/regnode.json"
  PASSWORD="synnergy"
  address=$(json_extract '.address' "$temp_wallet")
else
  temp_wallet=$(run_cli_json wallet new)
  address=$(json_extract '.address' "$temp_wallet")
fi

log_info "Approving transaction via regulatory node"
run_cli regnode approve "$address" "$MAX_AMOUNT" --wallet "$WALLET_FILE" --password "$PASSWORD" --json
run_cli_json regnode flag "$address" "manual-review" >/dev/null
logs=$(run_cli_json regnode logs "$address")
log_info "Regulatory logs: $logs"

log_info "Regulatory node setup completed"
