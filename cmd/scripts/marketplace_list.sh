#!/usr/bin/env bash
set -euo pipefail

CLI=${SYN_CLI:-./synnergy}
if [ ! -x "$CLI" ]; then
  echo "synnergy CLI not found at $CLI" >&2
  exit 1
fi

if [ $# -ne 2 ]; then
  echo "Usage: $0 <price> <cid>" >&2
  exit 1
fi
"$CLI" marketplace list "$1" "$2"
