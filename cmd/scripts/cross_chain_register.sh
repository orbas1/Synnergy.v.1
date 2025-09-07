#!/usr/bin/env bash
set -euo pipefail

CLI=${SYN_CLI:-./synnergy}
if [ ! -x "$CLI" ]; then
  echo "synnergy CLI not found at $CLI" >&2
  exit 1
fi

if [ $# -ne 3 ]; then
  echo "Usage: $0 <srcChain> <dstChain> <relayerAddress>" >&2
  exit 1
fi
"$CLI" cross_chain register "$1" "$2" "$3"
