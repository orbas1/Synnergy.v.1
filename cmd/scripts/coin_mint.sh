#!/usr/bin/env bash
set -euo pipefail

CLI=${SYN_CLI:-./synnergy}
if [ ! -x "$CLI" ]; then
  echo "synnergy CLI not found at $CLI" >&2
  exit 1
fi

if [ $# -lt 1 ]; then
  echo "Usage: $0 <address> [amount]" >&2
  exit 1
fi
addr=$1
amt=${2:-100}
"$CLI" ledger mint "$addr" "$amt"
