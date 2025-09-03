#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BIN="${BIN_PATH:-$SCRIPT_DIR/../cmd/synnergy}"

usage() {
  echo "Usage: $(basename "$0") path/to/testnet.yaml" >&2
  exit 1
}

CONFIG=${1:-}

if [[ -z "$CONFIG" ]]; then
  usage
fi

if [[ ! -f "$CONFIG" ]]; then
  echo "config file not found: $CONFIG" >&2
  exit 1
fi

if ! command -v "$BIN" &>/dev/null; then
  echo "binary not found: $BIN" >&2
  exit 1
fi

"$BIN" network start --config "$CONFIG"
