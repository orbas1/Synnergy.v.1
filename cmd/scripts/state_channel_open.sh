#!/usr/bin/env bash
set -euo pipefail

CLI=${SYN_CLI:-./synnergy}
if [ ! -x "$CLI" ]; then
  echo "synnergy CLI not found at $CLI" >&2
  exit 1
fi

if [ $# -lt 2 ]; then
  echo "Usage: $0 <from> <to> [amount]" >&2
  exit 1
fi
amt=${3:-1}
"$CLI" state-channel open "$1" "$2" "$amt"
