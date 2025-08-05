#!/usr/bin/env bash
set -e
if [ -z "$1" ]; then
  echo "Usage: $0 <tx.json>"
  exit 1
fi
./synnergy tx submit "$1"
