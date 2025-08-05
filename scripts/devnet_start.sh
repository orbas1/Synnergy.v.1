#!/usr/bin/env bash
set -e
BIN="$(dirname "$0")/../cmd/scripts/synnergy"
N=${1:-1}
for i in $(seq 1 "$N"); do
  "$BIN" network start --port $((3030+i)) &
done
trap 'kill 0' INT TERM
wait
