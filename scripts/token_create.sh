#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BIN="$SCRIPT_DIR/../cmd/scripts/synnergy"

if [[ $# -lt 5 ]]; then
  echo "Usage: $0 <name> <symbol> <owner> <decimals> <supply>" >&2
  exit 1
fi

"$BIN" syn500 create --name "$1" --symbol "$2" --owner "$3" --dec "$4" --supply "$5"
