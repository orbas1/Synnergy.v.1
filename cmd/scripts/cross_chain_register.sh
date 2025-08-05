#!/usr/bin/env bash
set -e
if [ $# -ne 3 ]; then
  echo "Usage: $0 <srcChain> <dstChain> <relayerAddress>"
  exit 1
fi
./synnergy cross_chain register "$1" "$2" "$3"
