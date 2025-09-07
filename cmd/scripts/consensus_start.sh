#!/usr/bin/env bash
set -euo pipefail

CLI=${SYN_CLI:-./synnergy}
if [ ! -x "$CLI" ]; then
  echo "synnergy CLI not found at $CLI" >&2
  exit 1
fi

interval=${1:-1000}
"$CLI" consensus-service start "$interval"
