#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BIN="$SCRIPT_DIR/../cmd/scripts/synnergy"

if [[ -z "${1:-}" ]]; then
  echo "Usage: $0 <file> [provider] [price] [capacity]" >&2
  exit 1
fi

if [[ -n "${2:-}" ]]; then
  provider="$2"
  price="${3:-0}"
  capacity="${4:-0}"
  "$BIN" marketplace pin "$1" "$provider" "$price" "$capacity"
else
  "$BIN" storage pin "$1"
fi
