#!/usr/bin/env bash
set -euo pipefail

CLI=${SYN_CLI:-./synnergy}
if [ ! -x "$CLI" ]; then
  echo "synnergy CLI not found at $CLI" >&2
  exit 1
fi

if [ $# -lt 1 ]; then
  echo "Usage: $0 <file> [provider] [price] [capacity]" >&2
  exit 1
fi
provider=${2:-provider}
price=${3:-0}
capacity=${4:-0}
"$CLI" marketplace pin "$1" "$provider" "$price" "$capacity"
