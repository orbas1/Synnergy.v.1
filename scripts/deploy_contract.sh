#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BIN="${BIN_PATH:-$SCRIPT_DIR/../cmd/synnergy}"

usage() {
  echo "Usage: $(basename "$0") path/to/contract.wasm" >&2
  exit 1
}

FILE=${1:-}

if [[ -z "$FILE" ]]; then
  usage
fi

if [[ ! -f "$FILE" ]]; then
  echo "contract file not found: $FILE" >&2
  exit 1
fi

if ! command -v "$BIN" &>/dev/null; then
  echo "binary not found: $BIN" >&2
  exit 1
fi

"$BIN" contracts deploy --wasm "$FILE"
