#!/usr/bin/env bash
set -e
if [ -z "$1" ]; then
  echo "Usage: $0 <file>"
  exit 1
fi
./synnergy storage pin "$1"
