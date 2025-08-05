#!/usr/bin/env bash
set -e
if [ -z "$1" ]; then
  echo "Usage: $0 <batch.json>"
  exit 1
fi
txs=$(cat "$1")
./synnergy rollups submit $txs
