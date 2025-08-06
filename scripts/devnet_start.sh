#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BIN="$SCRIPT_DIR/../cmd/scripts/synnergy"

readonly N="${1:-1}"
for ((i=1; i<=N; i++)); do
  "$BIN" network start --port $((3030+i)) &
done
trap 'kill 0' INT TERM
wait
