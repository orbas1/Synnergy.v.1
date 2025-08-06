#!/usr/bin/env bash
set -e
BIN="$(dirname "$0")/../cmd/scripts/synnergy"
if [ $# -lt 5 ]; then
  echo "Usage: $0 <name> <symbol> <owner> <decimals> <supply>"
  exit 1
fi
"$BIN" syn500 create --name "$1" --symbol "$2" --owner "$3" --dec "$4" --supply "$5"
