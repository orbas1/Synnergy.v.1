#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BIN="$SCRIPT_DIR/../cmd/scripts/synnergy"

if [[ $# -lt 2 ]]; then
  echo "Usage: $0 <address> <role>" >&2
  exit 1
fi

"$BIN" authority register "$1" "$2"
