#!/usr/bin/env bash
set -e
BIN="$(dirname "$0")/../cmd/scripts/synnergy"
if [ $# -lt 2 ]; then
  echo "Usage: $0 <title> <proposal_file>"
  exit 1
fi
"$BIN" governance propose "$1" "$2"
