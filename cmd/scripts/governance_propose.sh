#!/usr/bin/env bash
set -euo pipefail

CLI=${SYN_CLI:-./synnergy}
if [ ! -x "$CLI" ]; then
  echo "synnergy CLI not found at $CLI" >&2
  exit 1
fi

title=${1:-Proposal}
body=${2:-proposal.md}
"$CLI" governance propose "$title" "$body"
