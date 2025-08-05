#!/usr/bin/env bash
set -e
if [ -z "$1" ]; then
  echo "Usage: $0 <file.wasm>"
  exit 1
fi
./synnergy contracts deploy --wasm "$1"
