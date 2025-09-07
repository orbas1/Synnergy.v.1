#!/usr/bin/env bash
set -euo pipefail

CLI=${SYN_CLI:-./synnergy}
if [ ! -x "$CLI" ]; then
  echo "synnergy CLI not found at $CLI" >&2
  exit 1
fi

if [ $# -lt 1 ]; then
  echo "Usage: $0 <proposal-id> [approve]" >&2
  exit 1
fi
choice=${2:-yes}
"$CLI" dao-proposal vote "$1" default 1 "$choice"
