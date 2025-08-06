#!/usr/bin/env bash
set -e
BIN="$(dirname "$0")/../cmd/scripts/synnergy"
if [ $# -lt 2 ]; then
  echo "Usage: $0 <address> <role>"
  exit 1
fi
"$BIN" authority register "$1" "$2"
