#!/usr/bin/env bash
set -euo pipefail

CLI=${SYN_CLI:-./synnergy}
if [ ! -x "$CLI" ]; then
  echo "synnergy CLI not found at $CLI" >&2
  exit 1
fi

if [ $# -lt 3 ]; then
  echo "Usage: $0 <token> <from> <to> [amount]" >&2
  exit 1
fi
token=$1
from=$2
to=$3
amt=${4:-1}
"$CLI" "$token" transfer "$from" "$to" "$amt"
