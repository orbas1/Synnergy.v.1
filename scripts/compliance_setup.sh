#!/usr/bin/env bash
set -e
BIN="$(dirname "$0")/../cmd/scripts/synnergy"
if [ -z "$1" ]; then
  echo "Usage: $0 <kyc.json>"
  exit 1
fi
"$BIN" compliance validate "$1"
