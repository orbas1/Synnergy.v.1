#!/usr/bin/env bash
set -euo pipefail

CLI=${SYN_CLI:-./synnergy}
if [ ! -x "$CLI" ]; then
  echo "synnergy CLI not found at $CLI" >&2
  exit 1
fi

if [ $# -lt 2 ]; then
  echo "Usage: $0 <creator> <recipient> [type] [amount] [desc]" >&2
  exit 1
fi
type=${3:-personal}
amount=${4:-0}
desc=${5:-""}
"$CLI" loanpool submit "$1" "$2" "$type" "$amount" "$desc"
