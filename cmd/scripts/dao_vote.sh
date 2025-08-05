#!/usr/bin/env bash
set -e
if [ -z "$1" ]; then
  echo "Usage: $0 <proposal-id> [approve]"
  exit 1
fi
choice=${2:-yes}
./synnergy dao-proposal vote "$1" default 1 "$choice"
