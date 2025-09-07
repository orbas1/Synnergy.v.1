#!/usr/bin/env bash
set -euo pipefail

CLI=${SYN_CLI:-./synnergy}
if [ ! -x "$CLI" ]; then
  echo "synnergy CLI not found at $CLI" >&2
  exit 1
fi

if [ $# -lt 1 ]; then
  echo "Usage: $0 <address> [role]" >&2
  exit 1
fi
role=${2:-validator}
"$CLI" authority_apply submit "$1" "$role" "auto-generated"
