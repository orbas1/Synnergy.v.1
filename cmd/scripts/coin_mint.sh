#!/usr/bin/env bash
set -e
if [ -z "$1" ]; then
  echo "Usage: $0 <address> [amount]"
  exit 1
fi
addr=$1
amt=${2:-100}
./synnergy ledger mint "$addr" "$amt"
