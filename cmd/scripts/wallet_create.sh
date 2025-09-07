#!/usr/bin/env bash
set -euo pipefail

CLI=${SYN_CLI:-./synnergy}
if [ ! -x "$CLI" ]; then
  echo "synnergy CLI not found at $CLI" >&2
  exit 1
fi

out=${1:-wallet.json}
pass=${2:-password}

umask 077
"$CLI" wallet new > "$out"
echo "password: $pass" >> "$out"
