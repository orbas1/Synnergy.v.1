#!/usr/bin/env bash
set -e
BIN="$(dirname "$0")/../cmd/scripts/synnergy"
if [ $# -lt 5 ]; then
  echo "Usage: $0 <wasm_file> <model_hash> <manifest> <gas_limit> <owner>"
  exit 1
fi
"$BIN" ai_contract deploy "$1" "$2" "$3" "$4" "$5"
