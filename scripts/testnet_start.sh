#!/usr/bin/env bash
set -e
BIN="$(dirname "$0")/../cmd/scripts/synnergy"
if [ -z "$1" ]; then
  echo "Usage: $0 path/to/testnet.yaml"
  exit 1
fi
"$BIN" network start --config "$1"
