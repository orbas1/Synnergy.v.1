#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BIN="$SCRIPT_DIR/../cmd/scripts/synnergy"

if [[ $# -lt 5 ]]; then
  echo "Usage: $0 <wasm_file> <model_hash> <manifest> <gas_limit> <owner>" >&2
  exit 1
fi

"$BIN" ai_contract deploy "$1" "$2" "$3" "$4" "$5"
