#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BIN="$SCRIPT_DIR/../cmd/scripts/synnergy"

if [[ -z "${1:-}" ]]; then
  echo "Usage: $0 path/to/testnet.yaml" >&2
  exit 1
fi

"$BIN" network start --config "$1"
