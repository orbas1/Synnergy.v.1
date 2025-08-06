#!/usr/bin/env bash
set -e
BIN="$(dirname "$0")/../cmd/scripts/synnergy"
PORT=${1:-3030}
"$BIN" network start --port "$PORT"
