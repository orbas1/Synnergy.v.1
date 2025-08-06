#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BIN="$SCRIPT_DIR/../cmd/scripts/synnergy"

if [[ -z "${1:-}" ]]; then
  echo "Usage: $0 <kyc.json>" >&2
  exit 1
fi

"$BIN" compliance validate "$1"
