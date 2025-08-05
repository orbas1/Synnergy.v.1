#!/usr/bin/env bash
set -e
if [ -z "$1" ]; then
  echo "Usage: $0 <address>"
  exit 1
fi
./synnergy faucet request "$1"
