#!/usr/bin/env bash
set -e
if [ $# -ne 2 ]; then
  echo "Usage: $0 <price> <cid>"
  exit 1
fi
./synnergy marketplace list "$1" "$2"
