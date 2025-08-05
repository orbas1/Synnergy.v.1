#!/usr/bin/env bash
set -e
if [ -z "$1" ]; then
  echo "Usage: $0 \"leaf1,leaf2,...\""
  exit 1
fi
./synnergy security merkle "$1"
